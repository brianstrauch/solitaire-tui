package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlip(t *testing.T) {
	card := new(Card)
	card.Flip()

	require.True(t, card.IsVisible)
}
