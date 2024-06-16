package agent

import (
	"context"
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
		metricsCh <-chan MetricItem
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				mx:      sync.RWMutex{},
				metrics: tt.fields.metrics,
				options: tt.fields.options,
			}
			wg := &sync.WaitGroup{}
			wg.Add(1)
			a.processMetrics(wg, tt.args.metricsCh)
			wg.Wait()
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
			go a.processMetrics(wg, allMetricsCh)
			wg.Wait()
			// fmt.Println(a.metrics)
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
			go a.processMetrics(wg, allMetricsCh)
			wg.Wait()
			// fmt.Println(a.metr	ics)
		}
	})

}
