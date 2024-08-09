package agent

import (
	"context"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHttpClient struct {
}

func (c *MockHttpClient) Send(data []byte) error {
	return nil
}

func TestAgent_collectMetricsGenerator(t *testing.T) {
	type fields struct {
		metrics          map[string]MetricItem
		options          *Options
		httpClient       Sender
		metricCollection Storer
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   struct {
			count int
			keys  []string
		}
	}{
		{
			name: "success get metrics",
			fields: fields{
				metrics:          make(map[string]MetricItem),
				httpClient:       &MockHttpClient{},
				metricCollection: NewMetricCollector(10),
			},
			args: args{
				ctx: context.Background(),
			},
			want: struct {
				count int
				keys  []string
			}{
				count: 28,
				keys: []string{"Alloc",
					"BuckHashSys",
					"Frees",
					"GCCPUFraction",
					"GCSys",
					"HeapAlloc",
					"HeapIdle",
					"HeapInuse",
					"HeapObjects",
					"HeapReleased",
					"LastGC",
					"Lookups",
					"MCacheInuse",
					"MCacheSys",
					"MSpanInuse",
					"MSpanSys",
					"Mallocs",
					"NextGC",
					"NumForcedGC",
					"NumGC",
					"OtherSys",
					"PauseTotalNs",
					"StackInuse",
					"StackSys",
					"Sys",
					"TotalAlloc",
					"RandomValue"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAgent(tt.fields.metricCollection, tt.fields.httpClient, tt.fields.options)
			metrics := make([]MetricItem, 0)
			got := a.collectMetricsGenerator()
			for item := range got {
				metrics = append(metrics, item)
			}
			assert.Equal(t, tt.want.count, len(metrics))

			for _, v := range tt.want.keys {
				isContains := slices.ContainsFunc(metrics, func(m MetricItem) bool {
					return v == m.ID
				})
				assert.True(t, isContains)
			}
		})
	}
}

// func TestAgent_processMetrics(t *testing.T) {
// 	type fields struct {
// 		metrics map[string]MetricItem
// 		options *Options
// 	}

// 	type args struct {
// 		metricsCh chan MetricItem
// 	}

// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   map[string]MetricItem
// 	}{
// 		{
// 			name: "base proccess metric",
// 			fields: fields{
// 				metrics: make(map[string]MetricItem, 10),
// 				options: nil,
// 			},
// 			args: args{
// 				metricsCh: make(chan MetricItem, 7),
// 			},
// 			want: map[string]MetricItem{
// 				"gauge1":    {ID: "gauge1", MType: "gauge", Delta: 0, Value: 2.34},
// 				"gauge0":    {ID: "gauge0", MType: "gauge", Delta: 0, Value: 1.23},
// 				"gauge2":    {ID: "gauge2", MType: "gauge", Delta: 0, Value: 3.45},
// 				"gauge3":    {ID: "gauge3", MType: "gauge", Delta: 0, Value: 5.67},
// 				"PollCount": {ID: "PollCount", MType: "counter", Delta: 0, Value: 0},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(_ *testing.T) {
// 			a := &Agent{
// 				mx:      sync.RWMutex{},
// 				metrics: tt.fields.metrics,
// 				options: tt.fields.options,
// 			}
// 			// wg := &sync.WaitGroup{}
// 			// wg.Add(1)

// 			for i, v := range []float64{1.23, 2.34, 3.45, 5.67} {
// 				tt.args.metricsCh <- MetricItem{
// 					ID:    "gauge" + strconv.Itoa(i),
// 					MType: "gauge",
// 					Value: v,
// 				}
// 			}
// 			a.processMetrics(tt.args.metricsCh)
// 			close(tt.args.metricsCh)
// 			// wg.Wait()
// 			assert.Equal(t, 5, len(tt.fields.metrics))
// 			assert.EqualValues(t, tt.want, tt.fields.metrics)
// 		})
// 	}
// }

// func Benchmark_multiplexMetrics(b *testing.B) {
// 	b.Run("Run with init metircs store zero-sized", func(b *testing.B) {
// 		a := NewAgent(0, nil)
// 		ctx := context.Background()

// 		b.ResetTimer()

// 		for i := 0; i < b.N; i++ {
// 			wg := &sync.WaitGroup{}
// 			wg.Add(1)

// 			metricsCh := a.collectMetricsGenerator()
// 			addMetricsCh := a.collectAdditionalMetricsGenerator()
// 			allMetricsCh := multiplexChannels(ctx, metricsCh, addMetricsCh)

// 			go a.processMetrics(allMetricsCh)

// 			wg.Wait()
// 		}
// 	})

// 	b.Run("Run with presizing metircs store", func(b *testing.B) {
// 		a := NewAgent(100, nil)
// 		ctx := context.Background()

// 		b.ResetTimer()

// 		for i := 0; i < b.N; i++ {
// 			wg := &sync.WaitGroup{}
// 			wg.Add(1)

// 			metricsCh := a.collectMetricsGenerator()
// 			addMetricsCh := a.collectAdditionalMetricsGenerator()
// 			allMetricsCh := multiplexChannels(ctx, metricsCh, addMetricsCh)

// 			go a.processMetrics(allMetricsCh)

// 			wg.Wait()
// 		}
// 	})
// }

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
