package agent

import (
	"context"
	"slices"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgent_collectMetricsGenerator(t *testing.T) {
	type fields struct {
		metrics map[string]MetricItem
		options *Options
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "success get metrics",
			fields: fields{

				metrics: make(map[string]MetricItem),
			},
			args: args{
				ctx: context.Background(),
			},
			want: 28,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				mx:      sync.RWMutex{},
				metrics: tt.fields.metrics,
				options: tt.fields.options,
			}

			got := a.collectMetricsGenerator()

			var items = make([]MetricItem, 0, 30)

			for m := range got {
				items = append(items, m)
			}

			assert.Equal(t, tt.want, len(items))
		})
	}
}

func TestAgent_processMetrics(t *testing.T) {
	type fields struct {
		metrics map[string]MetricItem
		options *Options
	}

	type args struct {
		metricsCh chan MetricItem
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]MetricItem
	}{
		{
			name: "base proccess metric",
			fields: fields{
				metrics: make(map[string]MetricItem, 10),
				options: nil,
			},
			args: args{
				metricsCh: make(chan MetricItem, 7),
			},
			want: map[string]MetricItem{
				"gauge1":    {ID: "gauge1", MType: "gauge", Delta: 0, Value: 2.34},
				"gauge0":    {ID: "gauge0", MType: "gauge", Delta: 0, Value: 1.23},
				"gauge2":    {ID: "gauge2", MType: "gauge", Delta: 0, Value: 3.45},
				"gauge3":    {ID: "gauge3", MType: "gauge", Delta: 0, Value: 5.67},
				"PollCount": {ID: "PollCount", MType: "counter", Delta: 0, Value: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			a := &Agent{
				mx:      sync.RWMutex{},
				metrics: tt.fields.metrics,
				options: tt.fields.options,
			}
			// wg := &sync.WaitGroup{}
			// wg.Add(1)

			for i, v := range []float64{1.23, 2.34, 3.45, 5.67} {
				tt.args.metricsCh <- MetricItem{
					ID:    "gauge" + strconv.Itoa(i),
					MType: "gauge",
					Value: v,
				}
			}
			a.processMetrics(tt.args.metricsCh)
			close(tt.args.metricsCh)
			// wg.Wait()
			assert.Equal(t, 5, len(tt.fields.metrics))
			assert.EqualValues(t, tt.want, tt.fields.metrics)
		})
	}
}

func Benchmark_multiplexMetrics(b *testing.B) {
	b.Run("Run with init metircs store zero-sized", func(b *testing.B) {
		a := NewAgent(0, nil)
		ctx := context.Background()

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wg := &sync.WaitGroup{}
			wg.Add(1)

			metricsCh := a.collectMetricsGenerator()
			addMetricsCh := a.collectAdditionalMetricsGenerator()
			allMetricsCh := multiplexChannels(ctx, metricsCh, addMetricsCh)

			go a.processMetrics(allMetricsCh)

			wg.Wait()
		}
	})

	b.Run("Run with presizing metircs store", func(b *testing.B) {
		a := NewAgent(100, nil)
		ctx := context.Background()

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wg := &sync.WaitGroup{}
			wg.Add(1)

			metricsCh := a.collectMetricsGenerator()
			addMetricsCh := a.collectAdditionalMetricsGenerator()
			allMetricsCh := multiplexChannels(ctx, metricsCh, addMetricsCh)

			go a.processMetrics(allMetricsCh)

			wg.Wait()
		}
	})
}

func TestAgent_collectAdditionalMetricsGenerator(t *testing.T) {
	type fields struct {
		metrics map[string]MetricItem
		options *Options
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "simple test collect addiditional metrics",
			fields: fields{
				metrics: make(map[string]MetricItem),
				options: nil,
			},
			want: []string{"TotalMemory", "FreeMemory"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				mx:      sync.RWMutex{},
				metrics: tt.fields.metrics,
				options: tt.fields.options,
			}
			metricsCh := a.collectAdditionalMetricsGenerator()

			var resultMetrics = make([]MetricItem, 0, 10)
			for v := range metricsCh {
				resultMetrics = append(resultMetrics, v)
			}
			for _, f := range tt.want {
				isContains := slices.ContainsFunc(resultMetrics, func(m MetricItem) bool {
					return m.ID == f
				})
				assert.True(t, isContains)
			}
			assert.LessOrEqual(t, len(tt.want), len(resultMetrics))

		})
	}
}
