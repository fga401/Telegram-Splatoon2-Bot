package database

import (
	"testing"

	"github.com/stretchr/testify/require"
	"telegram-splatoon2-bot/driver/database"
)

func TestTokens(t *testing.T) {
	set := make(map[database.Token]struct{})
	for _, d := range statement {
		set[d.Token] = struct{}{}
	}
	require.Equal(t, len(set), len(statement), "All tokens are different.")
}
