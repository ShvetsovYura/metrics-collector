package types

import (
	"fmt"
	"strconv"
)

type Gauge float64
type Counter int64

func (g Gauge) ToString() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

func (c Counter) ToString() string {
	return fmt.Sprintf("%d", c)
}

type Sender interface {
	Send(string, string)
}

type Stringer interface {
	ToString() string
}

type Stored interface {
	UpdateGauge(name string, val float64)
	UpdateCounter(name string, val int64)
	GetVal(name string) (Stringer, error)
	ToList() []string
}
