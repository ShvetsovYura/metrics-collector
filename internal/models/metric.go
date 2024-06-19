package models

import "strconv"

type Gauge float64 // алиас для метрики gauge
type Counter int64 // алиас для метрики counter

// ToString, возвращает строковое представление метрики.
func (g Gauge) ToString() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

// GetRawValue, получает указатель исходное значение метрики.
func (g Gauge) GetRawValue() *float64 {
	val := float64(g)
	return &val
}

// ToString, возвращает строковое представление метрики.
func (c Counter) ToString() string {
	return strconv.FormatInt(int64(c), 10)
}

// GetRawValue, получает указатель исходное значение метрики.
func (c Counter) GetRawValue() *int64 {
	val := int64(c)
	return &val
}
