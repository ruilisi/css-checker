package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConf(t *testing.T) {
	found := make([]bool, 2)
	err := make([]error, 2)
	conf := Params{}
	found[0], err[0] = getConf(&conf, "tests/css-checker.yaml")
	found[1], err[1] = getConf(&conf, "tests/css-checker-notexist.yaml")
	assert.Equal(t, true, found[0])
	assert.NoError(t, err[0])
	assert.Equal(t, false, found[1])
	assert.NoError(t, err[1])
}
func TestHash(t *testing.T) {
	assert.Equal(t, uint64(0x541c3843ef77f983), hash("iudhsgvio6908&&gUikezjjdfl"))

}
func TestMin(t *testing.T) {
	assert.Equal(t, 1, min(1, 2))
	assert.Equal(t, 1, min(2, 1))
}

func TestGetConfResults(t *testing.T) {
	conf := Params{}
	found, err := getConf(&conf, "tests/css-checker.yaml")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, conf.LongScriptsCheck, false)
	assert.Equal(t, conf.ColorsCheck, false)
	assert.Equal(t, conf.Unused, true)
	assert.Equal(t, conf.LongScriptLength, 25)
}
