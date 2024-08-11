package agent

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
	pathToConfig := path.Join(basePath, "test-agent-config.json")

	want := &Options{
		PollInterval:   time.Duration(1 * time.Second),
		ReportInterval: time.Duration(300 * time.Millisecond),
		CryptoKey:      "abracadabra.pem",
		EndpointAddr:   "localhost:9876",
		LogLevel:       "debug",
	}
	errSetEnv := os.Setenv("CONFIG", pathToConfig)
	assert.NoError(t, errSetEnv)
	opt := ReadOptions()
	assert.Equal(t, *want, *opt)

}
