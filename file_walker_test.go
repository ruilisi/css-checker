package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileWalker(t *testing.T) {
	files, err := WalkMatch("tests", "*.css", []string{})
	assert.NoError(t, err)
	assert.Equal(t, len(files), 2)
}
