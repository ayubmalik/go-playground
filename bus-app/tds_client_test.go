package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMama(t *testing.T) {
	s := "hello"

	assert.Equal(t, "hello world", s)
}

func TestOrigins(t *testing.T) {
	client := TdsClient{
		key:       "some api key",
		carrierId: 777,
		url:       "some url",
	}

	origins, err := client.Origins()
	assert.Nil(t, err, "error was not nil")
	assert.NotEmpty(t, origins)
}
