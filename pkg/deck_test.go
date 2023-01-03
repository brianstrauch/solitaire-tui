package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDeck(t *testing.T) {
	deck := NewFullDeck()

	expected := TestDeck(
		"A♠?", "2♠?", "3♠?", "4♠?", "5♠?", "6♠?", "7♠?", "8♠?", "9♠?", "10♠?", "J♠?", "Q♠?", "K♠?",
		"A♦?", "2♦?", "3♦?", "4♦?", "5♦?", "6♦?", "7♦?", "8♦?", "9♦?", "10♦?", "J♦?", "Q♦?", "K♦?",
		"A♥?", "2♥?", "3♥?", "4♥?", "5♥?", "6♥?", "7♥?", "8♥?", "9♥?", "10♥?", "J♥?", "Q♥?", "K♥?",
		"A♣?", "2♣?", "3♣?", "4♣?", "5♣?", "6♣?", "7♣?", "8♣?", "9♣?", "10♣?", "J♣?", "Q♣?", "K♣?",
	)

	require.ElementsMatch(t, expected.cards, deck.cards)
}

func TestShuffle(t *testing.T) {
	deck := TestDeck("A♠", "2♠")
	deck.Shuffle()

	require.ElementsMatch(t, TestDeck("A♠", "2♠").cards, deck.cards)
}
