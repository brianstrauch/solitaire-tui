package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/brianstrauch/solitaire-tui/internal/solitaire"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	p := tea.NewProgram(solitaire.New(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
