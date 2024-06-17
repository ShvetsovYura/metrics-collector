package storage_test

import (
	"testing"

	"github.com/ShvetsovYura/metrics-collector/internal/models"
)

func TestGauge_GetRawValue(t *testing.T) {
	wanted := []float64{123.456, 0, -123.456}
	tests := []struct {
		name string
		g    models.Gauge
		want *float64
	}{
		{
			name: "get gauge raw value",
			g:    models.Gauge(wanted[0]),
			want: &wanted[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.GetRawValue(); *got != *tt.want {
				t.Errorf("Gauge.GetRawValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCounter_GetRawValue(t *testing.T) {
	wanted := []int64{123, 0, 100500}

	tests := []struct {
		name string
		c    models.Counter
		want *int64
	}{
		{
			name: "correct counter value",
			c:    models.Counter(wanted[0]),
			want: &wanted[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetRawValue(); *got != *tt.want {
				t.Errorf("Counter.GetRawValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
