package solitaire

import (
	"fmt"
	"strings"

	"github.com/brianstrauch/solitaire-tui/pkg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type deckType int

const (
	stock      deckType = 0
	waste      deckType = 1
	foundation deckType = 2
	tableau    deckType = 6
)

var deckTypes = []deckType{
	stock,
	waste,
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

var deckLocations = []cell{
	{0, 0},
	{6, 0},
	{18, 0},
	{24, 0},
	{30, 0},
	{36, 0},
	{0, 5},
	{6, 5},
	{12, 5},
	{18, 5},
	{24, 5},
	{30, 5},
	{36, 5},
}

type Solitaire struct {
	message string

	decks    []*pkg.Deck
	selected *index

	mouse        tea.MouseMsg
	windowHeight int
	maxHeight    int
}

type cell struct {
	x int
	y int
}

type index struct {
	deck int
	card int
}

func New() *Solitaire {
	decks := make([]*pkg.Deck, 13)

	decks[stock] = pkg.NewFullDeck()
	for i := 1; i < len(decks); i++ {
		decks[i] = pkg.NewEmptyDeck()
	}

	for i := 0; i < len(decks)-int(tableau); i++ {
		deck := decks[int(tableau)+i]
		for j := 0; j <= i; j++ {
			deck.Add(decks[stock].Pop())
		}
		deck.Top().Flip()
		deck.Expand()
	}

	return &Solitaire{
		message: "Solitaire",
		decks:   decks,
	}
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
	s.message = fmt.Sprintf("(%d, %d)", x, y)

	for i, deck := range s.decks {
		loc := deckLocations[i]
		if ok, j := deck.IsClicked(x-loc.x, y-loc.y); ok {
			switch deckTypes[i] {
			case stock:
				if deck.Size() > 0 {
					s.draw(3, deck, s.decks[waste])
				} else {
					s.draw(s.decks[waste].Size(), s.decks[waste], deck)
				}
			case waste:
				if deck.Size() > 0 {
					s.toggleSelect(&index{deck: i, card: deck.Size() - 1})
				}
			case foundation:
				if s.selected != nil && s.selected.deck != i {
					ok := s.move(&index{deck: i})
					if !ok {
						s.toggleSelect(&index{deck: i, card: deck.Size() - 1})
					}
				} else if deck.Size() > 0 {
					s.toggleSelect(&index{deck: i, card: deck.Size() - 1})
				}
			case tableau:
				s.message = fmt.Sprintf("%d %d", j, deck.Size()-1)
				if j == deck.Size()-1 && !deck.Top().IsVisible {
					if s.selected != nil {
						s.toggleSelect(s.selected)
					}
					deck.Top().Flip()
				} else if s.selected != nil && s.selected.deck != i {
					ok := s.move(&index{deck: i, card: j})
					if !ok {
						s.toggleSelect(&index{deck: i, card: j})
						s.toggleSelect(&index{deck: i, card: j})
					}
				} else if deck.Get(j).IsVisible {
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
			s.message += card.String()
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
		s.message = fmt.Sprintf("%d %d", toDeck.Size(), fromCards[0].Value)
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
	} else {
		s.selected = selected
		for _, card := range s.decks[s.selected.deck].GetFrom(s.selected.card) {
			card.IsSelected = true
		}
	}
}

func (s *Solitaire) View() string {
	view := lipgloss.JoinHorizontal(lipgloss.Top,
		s.decks[0].View(),
		s.decks[1].View(),
		strings.Repeat(" ", 6),
		s.decks[2].View(),
		s.decks[3].View(),
		s.decks[4].View(),
		s.decks[5].View(),
	) + "\n"

	view += lipgloss.JoinHorizontal(lipgloss.Top,
		s.decks[6].View(),
		s.decks[7].View(),
		s.decks[8].View(),
		s.decks[9].View(),
		s.decks[10].View(),
		s.decks[11].View(),
		s.decks[12].View(),
	) + "\n"

	return view
}
