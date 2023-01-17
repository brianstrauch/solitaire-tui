package pkg

import (
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Deck struct {
	cards      []*Card
	isExpanded bool
}

func NewDeck(cards []*Card) *Deck {
	return &Deck{cards: cards}
}

func NewFullDeck() *Deck {
	cards := make([]*Card, len(values)*len(suits))
	for i := range values {
		for j := range suits {
			cards[i*len(suits)+j] = NewCard(i, j)
		}
	}

	deck := NewDeck(cards)
	deck.Shuffle()

	return deck
}

func NewEmptyDeck() *Deck {
	return NewDeck(make([]*Card, 0))
}

func (d *Deck) Shuffle() {
	rand.Shuffle(d.Size(), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Expand() {
	d.isExpanded = true
}

func (d *Deck) View() string {
	if d == nil {
		return strings.Repeat(" ", Width)
	}

	// Outline
	if d.Size() == 0 {
		return viewCard(" ", "", lipgloss.AdaptiveColor{Light: "#EEEEEE", Dark: "#888888"})
	}

	// Expanded cards
	if d.isExpanded {
		var view string
		for i := 0; i < d.Size()-1; i++ {
			view += strings.Split(d.cards[i].View(), "\n")[0] + "\n"
		}
		return view + d.cards[d.Size()-1].View()
	}

	// Top card only
	return d.cards[d.Size()-1].View()
}

func (d *Deck) IsClicked(x, y int) (bool, int) {
	if d == nil {
		return false, 0
	}

	if d.Size() == 0 {
		return x >= 0 && x < Width && y >= 0 && y < Height, 0
	}

	if d.isExpanded {
		for i := d.Size() - 1; i >= 0; i-- {
			if x >= 0 && x < Width && y >= i && y < i+Height {
				return true, i
			}
		}
		return false, 0
	}

	return x >= 0 && x < Width && y >= 0 && y < Height, 0
}

func (d *Deck) Add(cards ...*Card) {
	d.cards = append(d.cards, cards...)
}

func (d *Deck) Top() *Card {
	return d.Get(d.Size() - 1)
}

func (d *Deck) Bottom() *Card {
	return d.Get(0)
}

func (d *Deck) Get(idx int) *Card {
	return d.cards[idx]
}

func (d *Deck) GetFrom(idx int) []*Card {
	return d.cards[idx:]
}

func (d *Deck) Pop() *Card {
	if len(d.cards) > 0 {
		return d.PopFrom(d.Size() - 1)[0]
	}

	return nil
}

func (d *Deck) PopFrom(idx int) []*Card {
	cards := d.cards[idx:]
	d.cards = d.cards[:idx]
	return cards
}

func (d *Deck) Size() int {
	return len(d.cards)
}

// TestDeck is a helper function to simplify testing.
func TestDeck(shorthands ...string) *Deck {
	cards := make([]*Card, len(shorthands))
	for i, shorthand := range shorthands {
		cards[i] = testCard(shorthand)
	}
	return &Deck{cards: cards}
}

func testCard(shorthand string) *Card {
	card := &Card{IsVisible: !strings.HasSuffix(shorthand, "?")}

	for i, value := range values {
		if strings.HasPrefix(shorthand, value) {
			card.Value = i
			break
		}
	}

	for i, suit := range suits {
		if strings.Contains(shorthand, suit) {
			card.Suit = i
			break
		}
	}

	return card
}
