package agent

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgent_collectMetricsGenerator(t *testing.T) {
	type fields struct {
		mx      sync.RWMutex
		metrics map[string]MetricItem
		options *AgentOptions
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
				mx:      sync.RWMutex{},
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
				mx:      tt.fields.mx,
				metrics: tt.fields.metrics,
				options: tt.fields.options,
			}

			got := a.collectMetricsGenerator(tt.args.ctx)
			var items = make([]MetricItem, 0, 30)
			for m := range got {
				items = append(items, m)
			}
			assert.Equal(t, tt.want, len(items))
		})
	}
}
