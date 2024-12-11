package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlug(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		slugString := "20st century boy 155"
		slug := NewSlug(slugString)

		require.NotEmpty(t, slug)
		require.Equal(t, "20st-century-boy-155", slug)
	})
}
