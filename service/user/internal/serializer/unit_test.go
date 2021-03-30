package serializer

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"telegram-splatoon2-bot/service/user/database"
)

func TestUserID(t *testing.T) {
	testcases := []database.UserID{
		-123,
		0,
		123,
		math.MaxInt64,
		math.MinInt64,
	}
	for _, expected := range testcases {
		data := FromID(expected)
		actual := ToID(data)
		require.Equal(t, expected, actual, "ID should be the same after serialization and deserialization.")
	}
}

func TestStatus(t *testing.T) {
	testcases := []database.Status{
		database.Status{
			UserID:       0,
			SessionToken: "",
			IKSM:         "",
			Language:     "",
			Timezone:     0,
			LastBattle:   "",
			LastSalmon:   "",
		},
		database.Status{
			UserID:       123456789,
			SessionToken: "abc.defghi.jklm.nopqrstuvw.xyz",
			IKSM:         "0000000000000000000000000000000000000000",
			Language:     "en",
			Timezone:     720,
			LastBattle:   "123456",
			LastSalmon:   "123456",
		},
	}
	for _, expected := range testcases {
		data := FromStatus(expected)
		actual := ToStatus(data)
		require.Equal(t, expected, actual, "Status should be the same after serialization and deserialization.")
	}
}

func TestAccounts(t *testing.T) {
	testcases := [][]database.Account{
		[]database.Account{},
		[]database.Account{
			{
				UserID:       0,
				SessionToken: "",
				Tag:          "",
			},
			{
				UserID:       123456789,
				SessionToken: "abc.defghi.jklm.nopqrstuvw.xyz",
				Tag:          "sadasfasda:15820asd",
			},
		},
	}
	for _, expected := range testcases {
		data := FromAccounts(expected)
		actual := ToAccounts(data)
		require.Equal(t, expected, actual, "Accounts should be the same after serialization and deserialization.")
	}
}
