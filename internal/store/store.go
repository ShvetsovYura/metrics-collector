package store

import (
	"github.com/ShvetsovYura/metrics-collector/internal/store/memory"
)

type Store interface {
	SetGauge(name string, val float64) error
	SetCounter() error
}

func NewStore() (Store, error) {
	return memory.NewStore(40), nil
}
