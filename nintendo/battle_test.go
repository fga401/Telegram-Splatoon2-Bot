package nintendo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetResults(t *testing.T) {
	prepareTest()
	_, err := GetResults("72811dcfc39d199333bc466549a80ac0ff06070d", 480,"en-US")
	assert.Nil(t, err)
}

