package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCounter(t *testing.T) {
	m := NewMetrics()
	m.SetCounter()
	m.SetCounter()

	assert.Equal(t, counter(2), m["PollCounter"])
}
