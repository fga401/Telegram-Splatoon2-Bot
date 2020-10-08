package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	ret := getSplatoonNextUpdateTime(time.Now())
	fmt.Println(ret)
	ret = getSplatoonNextUpdateTime(time.Now().Add(time.Hour))
	fmt.Println(ret)
	ret = getSplatoonNextUpdateTime(getSplatoonNextUpdateTime(time.Now()))
	fmt.Println(ret)
	assert.Nil(t,nil)
}
