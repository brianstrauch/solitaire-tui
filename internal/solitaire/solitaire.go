package solitaire

import (
	"github.com/brianstrauch/solitaire-tui/pkg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type deckType int

const (
	stock      deckType = 0
	waste      deckType = 1
	empty      deckType = 2
	foundation deckType = 3
	tableau    deckType = 7
)

var deckTypes = []deckType{
	stock,
	waste,
	empty,
	foundation,
	foundation,
	foundation,
	foundation,
	tableau,
	tableau,
	tableau,
	tableau,
	tableau,
	tableau,
	tableau,
}

type Solitaire struct {
	decks    []*pkg.Deck
	selected *index

	mouse        tea.MouseMsg
	windowHeight int
	maxHeight    int
}

type index struct {
	deck int
	card int
}

func New() *Solitaire {
	decks := make([]*pkg.Deck, len(deckTypes))
	for i := range decks {
		switch deckTypes[i] {
		case stock:
			decks[i] = pkg.NewFullDeck()
		case empty:
			decks[i] = nil
		default:
			decks[i] = pkg.NewEmptyDeck()
		}
	}

	for i := 0; i < len(decks)-int(tableau); i++ {
		deck := decks[int(tableau)+i]
		for j := 0; j <= i; j++ {
			deck.Add(decks[stock].Pop())
		}
		deck.Top().Flip()
		deck.Expand()
	}

	return &Solitaire{decks: decks}
}

func (s *Solitaire) Init() tea.Cmd {
	return nil
}

func (s *Solitaire) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return s, tea.Quit
		}
	case tea.WindowSizeMsg:
		s.windowHeight = msg.Height
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			if s.mouse.Type != tea.MouseLeft {
				s.mouse = msg
			}
		case tea.MouseRelease:
			if s.mouse.Type == tea.MouseLeft && msg.X == s.mouse.X && msg.Y == s.mouse.Y {
				height := lipgloss.Height(s.View())
				if height > s.maxHeight {
					s.maxHeight = height
				}
				y := msg.Y - (s.windowHeight - s.maxHeight)
				s.click(msg.X, y)
			}
			s.mouse = msg
		}
	}

	return s, nil
}

func (s *Solitaire) click(x, y int) {
	n := len(s.decks) / 2

	for i, deck := range s.decks {
		xi := (i % n) * pkg.Width
		yi := (i / n) * pkg.Height

		if ok, j := deck.IsClicked(x-xi, y-yi); ok {
			switch deckTypes[i] {
			case stock:
				if deck.Size() > 0 {
					s.draw(3, deck, s.decks[waste])
				} else {
					s.draw(s.decks[waste].Size(), s.decks[waste], deck)
				}
			case waste:
				if deck.Size() > 0 {
					if s.selected != nil && s.selected.deck != i {
						s.toggleSelect(nil)
					}
					s.toggleSelect(&index{deck: i, card: deck.Size() - 1})
				}
			case foundation:
				if s.selected != nil && s.selected.deck != i {
					if !s.move(&index{deck: i}) && deck.Size() > 0 {
						s.toggleSelect(nil)
						s.toggleSelect(&index{deck: i, card: deck.Size() - 1})
					}
				} else if deck.Size() > 0 {
					if s.selected != nil && s.selected.deck != i {
						s.toggleSelect(nil)
					}
					s.toggleSelect(&index{deck: i, card: deck.Size() - 1})
				}
			case tableau:
				if j == deck.Size()-1 && !deck.Top().IsVisible {
					if s.selected != nil {
						s.toggleSelect(s.selected)
					}
					deck.Top().Flip()
				} else if s.selected != nil && s.selected.deck != i {
					ok := s.move(&index{deck: i, card: j})
					if !ok {
						s.toggleSelect(nil)
						s.toggleSelect(&index{deck: i, card: j})
					}
				} else if deck.Size() > 0 && deck.Get(j).IsVisible {
					if s.selected != nil && s.selected.deck == i && s.selected.card != j {
						s.toggleSelect(&index{deck: i, card: j})
					}
					s.toggleSelect(&index{deck: i, card: j})
				}
			}

			break
		}
	}
}

func (s *Solitaire) draw(n int, from, to *pkg.Deck) {
	if s.selected != nil {
		s.toggleSelect(s.selected)
	}

	for i := 0; i < n; i++ {
		if card := from.Pop(); card != nil {
			card.Flip()
			to.Add(card)
		}
	}
}

func (s *Solitaire) move(to *index) bool {
	toDeck := s.decks[to.deck]
	fromDeck := s.decks[s.selected.deck]
	fromCards := fromDeck.GetFrom(s.selected.card)

	switch deckTypes[to.deck] {
	case foundation:
		if s.selected.card == fromDeck.Size()-1 && toDeck.Size() == 0 && fromDeck.Top().Value == 0 || toDeck.Size() > 0 && fromDeck.Top().Value == toDeck.Top().Value+1 && fromDeck.Top().Suit == toDeck.Top().Suit {
			s.toggleSelect(s.selected)
			toDeck.Add(fromDeck.Pop())
			s.selected = nil
			return true
		}
	case tableau:
		if toDeck.Size() == 0 && fromCards[0].Value == 12 || toDeck.Size() > 0 && fromCards[0].Value+1 == toDeck.Top().Value && fromCards[0].Color() != toDeck.Top().Color() {
			idx := s.selected.card
			s.toggleSelect(s.selected)
			toDeck.Add(fromDeck.PopFrom(idx)...)
			s.selected = nil
			return true
		}
	}

	return false
}

func (s *Solitaire) toggleSelect(selected *index) {
	if s.selected != nil {
		for _, card := range s.decks[s.selected.deck].GetFrom(s.selected.card) {
			card.IsSelected = false
		}
		s.selected = nil
	} else if s.decks[selected.deck].Size() > 0 {
		s.selected = selected
		for _, card := range s.decks[s.selected.deck].GetFrom(s.selected.card) {
			card.IsSelected = true
		}
	}
}

func (s *Solitaire) View() string {
	n := len(s.decks) / 2

	var view string
	for i := 0; i < 2; i++ {
		row := make([]string, n)
		for j := range row {
			row[j] = s.decks[i*n+j].View()
		}
		view += lipgloss.JoinHorizontal(lipgloss.Top, row...) + "\n"
	}

	return view
}
