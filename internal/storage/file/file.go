package file

import (
	"bytes"
	"encoding/json"
	"os"
)

type DumpItem struct {
	Gauges map[string]float64 `json:"gauges"`
	Cnt    int64              `json:"counter"`
}

type FileStorage struct {
	path string
}

func NewFileStorage(pathToFile string) *FileStorage {
	return &FileStorage{path: pathToFile}
}

func (fs *FileStorage) Dump(gauges map[string]float64, counter int64) error {
	f, err := os.OpenFile(fs.path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return nil
	}
	defer f.Close()

	di := DumpItem{Gauges: gauges, Cnt: counter}

	data, err := json.MarshalIndent(di, "", "  ")
	if err != nil {
		return err
	}
	f.Write(data)

	return nil
}

func (fs *FileStorage) Restore() (map[string]float64, int64, error) {
	var buf bytes.Buffer
	f, err := os.OpenFile(fs.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, 0, err
	}

	_, err1 := buf.ReadFrom(f)
	if err1 != nil {
		return nil, 0, err1
	}
	di := DumpItem{}
	err2 := json.Unmarshal(buf.Bytes(), &di)
	if err2 != nil {
		return nil, 0, err2
	}

	return di.Gauges, di.Cnt, nil

}
