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
	Width  = 6
	Height = 5
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
	color := lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}

	if c.IsSelected {
		color = lipgloss.AdaptiveColor{Light: "#FFFF00", Dark: "#FFFF00"}
	}

	if !c.IsVisible {
		return viewCard("╱", "", color)
	}

	style := lipgloss.NewStyle().Foreground(c.Color())
	return viewCard(" ", style.Render(c.String()), color)
}

func (c *Card) Flip() {
	c.IsVisible = !c.IsVisible
}

func (c *Card) Color() lipgloss.AdaptiveColor {
	if c.Suit == 1 || c.Suit == 2 {
		return lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}
	} else {
		return lipgloss.AdaptiveColor{Light: "#000000", Dark: "#888888"}
	}
}

func (c *Card) String() string {
	return values[c.Value] + suits[c.Suit]
}

func viewCard(design, shorthand string, color lipgloss.AdaptiveColor) string {
	style := lipgloss.NewStyle().Foreground(color)
	padding := strings.Repeat("─", Width-2-lipgloss.Width(shorthand))

	view := style.Render("╭") + shorthand + style.Render(padding+"╮") + "\n"
	for i := 1; i < Height-1; i++ {
		view += style.Render("│"+strings.Repeat(design, Width-2)+"│") + "\n"
	}
	view += style.Render("╰"+padding) + shorthand + style.Render("╯")

	return view
}
