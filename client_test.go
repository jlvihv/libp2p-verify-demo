package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_client(t *testing.T) {
	err := client(8080)
	assert.NoError(t, err)
}
