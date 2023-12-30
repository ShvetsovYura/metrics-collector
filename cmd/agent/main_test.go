package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeLink(t *testing.T) {
	tests := []struct {
		name    string
		mName   string
		val     any
		wantErr bool
		want    string
	}{
		{
			name:    "correct gauge metric",
			mName:   "Alloc",
			val:     gauge(123.456),
			wantErr: false,
			want:    "http://localhost:8080/update/gauge/Alloc/123.456",
		}, {
			name:    "incorrect gauge metric name",
			mName:   "Abracadabra",
			val:     gauge(0.335),
			wantErr: true,
			want:    "",
		}, {
			name:    "incorrect gauge metric value",
			mName:   "Alloc",
			val:     "abracadabra",
			wantErr: true,
			want:    "",
		}, {
			name:    "nil gauge metric value",
			mName:   "Alloc",
			val:     nil,
			wantErr: true,
			want:    "",
		},
		{
			name:    "correct counter value ",
			mName:   "PollCount",
			val:     counter(123),
			wantErr: false,
			want:    "http://localhost:8080/update/counter/PollCount/123",
		},
		{
			name:    "incorrect counter value",
			mName:   "PollCount",
			val:     gauge(0.123),
			wantErr: true,
			want:    "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := makeLink(test.mName, test.val)
			if test.wantErr == true {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want, res)
		})
	}
}

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
			res := contains(values, test.val)
			assert.Equal(t, test.want, res)
		})
	}
}
