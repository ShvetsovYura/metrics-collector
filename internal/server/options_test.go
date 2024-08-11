package server

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadOptions(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	basePath := path.Join(cwd, "..", "..", "testdata")
	pathToConfig := path.Join(basePath, "test-server-config.json")

	want := &Options{
		EndpointAddr:    "localhost:6789",
		Restore:         true,
		StoreInterval:   time.Duration(600 * time.Second),
		CryptoKey:       "hoho.pem",
		LogLevel:        "debug",
		FileStoragePath: "/tmp/metrics-db.json",
	}
	errSetEnv := os.Setenv("CONFIG", pathToConfig)
	assert.NoError(t, errSetEnv)
	opt := ReadOptions()
	assert.Equal(t, *want, *opt)

}
