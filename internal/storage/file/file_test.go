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
	g["OtherMetric"] = 123
	var c int64 = 10
	err := fs.Dump(g, c)

	t.Run("test dump & resotre", func(t *testing.T) {
		assert.NoError(t, err)
		assert.FileExists(t, path)
		gauges, counter, err := fs.Restore()

		assert.NoError(t, err)

		assert.Equal(t, g, gauges)
		assert.Equal(t, c, counter)
	})
}
