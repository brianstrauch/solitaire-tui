package pkg

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	values = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	suits  = []string{"♠", "♦", "♥", "♣"}
)

const (
	width  = 6
	height = 5
)

type Card struct {
	Value      int
	Suit       int
	IsVisible  bool
	IsSelected bool
}

func NewCard(value, suit int) *Card {
	return &Card{
		Value: value,
		Suit:  suit,
	}
}

func (c *Card) View() string {
	if !c.IsVisible {
		return viewCard("╱", "", c.IsSelected)
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Color()))
	return viewCard(" ", style.Render(c.String()), c.IsSelected)
}

func (c *Card) Flip() {
	c.IsVisible = !c.IsVisible
}

func (c *Card) Color() string {
	if c.Suit == 1 || c.Suit == 2 {
		return "#FF0000"
	} else {
		return "#000000"
	}
}

func (c *Card) String() string {
	return values[c.Value] + suits[c.Suit]
}

func viewCard(design, shorthand string, isSelected bool) string {
	style := lipgloss.NewStyle()
	if isSelected {
		style = style.Foreground(lipgloss.Color("#FFFF00"))
	}

	padding := strings.Repeat("─", width-2-lipgloss.Width(shorthand))

	view := style.Render("╭") + shorthand + style.Render(padding+"╮") + "\n"
	for i := 1; i < height-1; i++ {
		view += style.Render("│"+strings.Repeat(design, width-2)+"│") + "\n"
	}
	view += style.Render("╰"+padding) + shorthand + style.Render("╯")

	return view
}
