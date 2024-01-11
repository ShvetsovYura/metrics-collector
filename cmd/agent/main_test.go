package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCounter(t *testing.T) {
	m := NewMetrics(30)
	increaseCounter(m)
	increaseCounter(m)

	assert.Equal(t, counter(2), m["PollCounter"])
}
