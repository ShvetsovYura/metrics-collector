package metric

import "strconv"

type Gauge float64
type Counter int64

type Metric interface {
	ToString() string
}

func (g Gauge) ToString() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

func (c Counter) ToString() string {
	return strconv.FormatInt(int64(c), 10)
}
