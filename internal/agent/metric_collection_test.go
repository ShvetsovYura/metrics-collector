package agent

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsCollection_Items(t *testing.T) {
	metrics := map[string]MetricItem{
		"metric1": {ID: "m1", MType: "gauge", Value: 123.45, Delta: 0},
		"metric2": {ID: "m2", MType: "gauge", Value: -0.453234, Delta: 0},
		"metric3": {ID: "m3", MType: "gauge", Value: 123213.45, Delta: 0},
		"metric4": {ID: "m4", MType: "gauge", Value: 0, Delta: 0},
		"metric5": {ID: "m5", MType: "counter", Value: 0, Delta: 10},
	}
	type fields struct {
		metrics map[string]MetricItem
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]MetricItem
	}{
		{
			name: "all elements returning from iterator",
			fields: struct {
				metrics map[string]MetricItem
			}{
				metrics: metrics,
			},
			want: metrics,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &MetricsCollection{
				metrics: tt.fields.metrics,
			}
			var metrics []MetricItem
			next := mc.Items()
			for {
				val, hasNext := next()
				metrics = append(metrics, val)
				if !hasNext {
					break
				}
			}

			assert.Equal(t, len(tt.want), len(metrics))

			for _, v := range tt.want {
				isContain := slices.Contains(metrics, v)
				assert.True(t, isContain)
			}

			// if got := mc.Items(); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("MetricsCollection.Items() = %v, want %v", got, tt.want)
			// }
		})
	}
}
