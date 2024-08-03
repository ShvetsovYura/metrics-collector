package util

import (
	"os"
	"path"
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

func TestHash(t *testing.T) {
	cases := []struct {
		name       string
		inputBytes []byte
		inputKey   string
		outHash    string
	}{{
		name:       "case1",
		inputBytes: []byte("myinputstring"),
		inputKey:   "private_key",
		outHash:    "3ed15a8efa23d4c73682b9b0e5b953e362b317acbc0d3e242746be625ec03cf7",
	}, {
		name:       "case2",
		inputBytes: []byte(""),
		inputKey:   "",
		outHash:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := Hash(c.inputBytes, c.inputKey)
			assert.Equal(t, c.outHash, result)
		})
	}
}

func TestEncryptData(t *testing.T) {
	testMessage := []byte("this test message")
	cwd, err := os.Getwd()
	assert.NoError(t, err)
	basePath := path.Join(cwd, "..", "..", "testdata")
	publicKeyPath := path.Join(basePath, "public.pem")
	privateKeyPath := path.Join(basePath, "private.pem")

	encryptedMessage, err := EncryptData(testMessage, publicKeyPath)

	assert.NoError(t, err)
	decryptedMessage, err := DecryptData(encryptedMessage, privateKeyPath)
	assert.NoError(t, err)
	assert.Equal(t, testMessage, decryptedMessage)
}
