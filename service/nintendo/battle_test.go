package nintendo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	iksm     = "82b55de7e27a253e70154c5b0b9888a9d14f1da3"
	timezone = 480
)

func TestGetAllBattleResults(t *testing.T) {
	prepareTest()
	ret, err := GetAllBattleResults(iksm, timezone, acceptLang)
	assert.Nil(t, err)
	println(ret)
}
