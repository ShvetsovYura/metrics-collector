package utilu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	values := []string{"one", "two", "four", "five"}
	tests := []struct {
		name string
		val  string
		want bool
	}{{
		name: "contains in", val: "one", want: true,
	}, {
		name: "not contains", val: "three", want: false,
	}, {
		name: "not contains", val: "ten", want: false,
	}, {
		name: "contains", val: "five", want: true,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := Contains(values, test.val)
			assert.Equal(t, test.want, res)
		})
	}
}
