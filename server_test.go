package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_server(t *testing.T) {
	err := start()
	assert.NoError(t, err)
}
