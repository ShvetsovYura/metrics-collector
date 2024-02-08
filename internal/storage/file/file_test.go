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
	fs := NewFileStorage(path, 40, false, 0)
	var g = make(map[string]float64, 10)
	var c = make(map[string]int64, 2)
	g["Alloc"] = 44.1
	g["OtherMetric"] = 123
	c["PollCount"] = 10
	err := fs.Dump(g, c)

	t.Run("test dump & resotre", func(t *testing.T) {
		assert.NoError(t, err)
		assert.FileExists(t, path)
		gauges, counter, err := fs.RestoreNow()

		assert.NoError(t, err)

		assert.Equal(t, g, gauges)
		assert.Equal(t, c, counter)
	})
}
