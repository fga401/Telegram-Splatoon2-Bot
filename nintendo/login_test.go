package nintendo

import "testing"
import "github.com/stretchr/testify/assert"

func TestGetSessionToken(t *testing.T) {
	err := getSessionToken()
	assert.Nil(t, err)
}
