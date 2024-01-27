package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDump(t *testing.T) {
	path := "tt.txt"
	defer func() {
		os.Remove(path)
	}()
	fs := NewFileStorage(path)
	var g map[string]float64 = make(map[string]float64, 10)
	g["Alloc"] = 44.1
	err := fs.Dump(g, 10)

	t.Run("test dump", func(t *testing.T) {
		assert.NoError(t, err)
		assert.FileExists(t, path)
	})
}
