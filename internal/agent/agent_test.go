package agent

import (
	"context"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
}

func (c *MockHTTPClient) Send(data []byte) error {
	return nil
}

func TestAgent_collectMetricsGenerator(t *testing.T) {
	type fields struct {
		options    *Options
		httpClient Sender
		count      int
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
				httpClient: &MockHTTPClient{},
				count:      10,
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
			collection := NewMetricCollector(tt.fields.count)
			a := NewAgent(collection, tt.fields.httpClient, tt.fields.options)
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

func TestAgent_collectAdditionalMetricsGenerator(t *testing.T) {
	type fields struct {
		options    *Options
		httpClient Sender
		count      int
	}
	tests := []struct {
		name   string
		fields fields
		want   struct {
			count int
			keys  []string
		}
	}{
		{
			name: "simple test collect addiditional metrics",
			fields: fields{
				httpClient: &MockHTTPClient{},
				options:    nil,
				count:      10,
			},
			want: struct {
				count int
				keys  []string
			}{
				count: 3,
				keys: []string{
					"TotalMemory",
					"FreeMemory",
					"CPUutilization1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := make([]MetricItem, 0)
			c := NewMetricCollector(tt.fields.count)
			a := NewAgent(c, tt.fields.httpClient, tt.fields.options)
			got := a.collectAdditionalMetricsGenerator()
			for item := range got {
				metrics = append(metrics, item)
			}
			assert.LessOrEqual(t, tt.want.count, len(metrics))

			for _, v := range tt.want.keys {
				isContains := slices.ContainsFunc(metrics, func(m MetricItem) bool {
					return v == m.ID
				})
				assert.True(t, isContains)
			}

		})
	}
}

func TestAgent_runCollectMetrics(t *testing.T) {
	tests := []struct {
		name string
		want struct {
			count int
		}
	}{
		{
			name: "test ",
			want: struct{ count int }{
				count: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pollInterval, _ := time.ParseDuration("200ms")
			opt := Options{
				PollInterval: pollInterval,
			}
			c := NewMetricCollector(10)
			s := &MockHTTPClient{}
			a := NewAgent(c, s, &opt)
			ctx, cancel := context.WithCancel(context.Background())
			wg := sync.WaitGroup{}
			wg.Add(1)

			go a.runCollectMetrics(ctx, &wg)
			go func() {
				time.Sleep(time.Duration(1 * time.Second))
				cancel()

			}()
			wg.Wait()
			assert.LessOrEqual(t, tt.want.count, a.collection.Count())

		})
	}
}
