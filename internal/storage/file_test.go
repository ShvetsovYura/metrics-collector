package storage

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ShvetsovYura/metrics-collector/internal/models"
	"github.com/ShvetsovYura/metrics-collector/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDump(t *testing.T) {
	mem := NewMemory(40)
	path := "tt.txt"

	defer func() {
		err := os.Remove(path)
		if err != nil {
			fmt.Printf("Не удается удалить файл, %s", err.Error())
		}
	}()

	fs := NewFile(path, mem, false, 0)

	var (
		g = make(map[string]float64, 10)
		c = make(map[string]int64, 2)
	)

	g["Alloc"] = 44.1
	g["OtherMetric"] = 123
	c["PollCount"] = 10
	err := fs.Dump(g, c)

	t.Run("test dump & resotre", func(t *testing.T) {
		assert.NoError(t, err)
		assert.FileExists(t, path)

		gauges, counter, err := fs.RestoreNow()

		assert.NoError(t, err)

		assert.Equal(t, g, gauges)
		assert.Equal(t, c, counter)
	})
}

func TestFileStorage_ExtractGauges(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockMemoryStore(ctrl)

	type fields struct {
		path        string
		immediately bool
		memStorage  MemoryStore
	}

	type args struct {
		ctx     context.Context
		mockOut map[string]models.Gauge
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]float64
	}{
		{
			name: "non zero values",
			fields: fields{
				path:        "test.txt",
				immediately: false,
				memStorage:  m,
			},
			args: args{
				ctx: context.Background(),
				mockOut: map[string]models.Gauge{
					"memsize": models.Gauge(1234.45),
					"oprate":  models.Gauge(0.3566),
					"other":   models.Gauge(-134.3),
				},
			},
			want: map[string]float64{
				"memsize": 1234.45,
				"oprate":  0.3566,
				"other":   -134.3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().GetGauges(tt.args.ctx).Return(tt.args.mockOut)
			fs := &File{
				path:        tt.fields.path,
				immediately: tt.fields.immediately,
				memStorage:  tt.fields.memStorage,
			}
			result := fs.ExtractGauges(tt.args.ctx)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestFileStorage_ExtractCounters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockMemoryStore(ctrl)

	type fields struct {
		path        string
		immediately bool
		memStorage  MemoryStore
	}

	type args struct {
		ctx     context.Context
		mockOut map[string]models.Counter
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]int64
	}{
		{
			name: "non zero values",
			fields: fields{
				path:        "test.txt",
				immediately: false,
				memStorage:  m,
			},
			args: args{
				ctx: context.Background(),
				mockOut: map[string]models.Counter{
					"counter":      152,
					"zero_counter": 0,
				},
			},
			want: map[string]int64{
				"counter":      152,
				"zero_counter": 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().GetCounters(tt.args.ctx).Return(tt.args.mockOut)
			fs := &File{
				path:        tt.fields.path,
				immediately: tt.fields.immediately,
				memStorage:  tt.fields.memStorage,
			}
			result := fs.ExtractCounters(tt.args.ctx)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestFile_Restore(t *testing.T) {
	type fields struct {
		path        string
		immediately bool
		memStorage  MemoryStore
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    map[string]float64
	}{
		{
			name: "basic restore from file",
			fields: fields{
				path:        "test.txt",
				immediately: false,
				memStorage:  NewMemory(10),
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
			want: map[string]float64{
				"gauge1": 94.1234,
			},
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &File{
				path:        tt.fields.path,
				immediately: tt.fields.immediately,
				memStorage:  tt.fields.memStorage,
			}
			for k, v := range tt.want {
				err := fs.memStorage.SetGauge(ctx, k, v)
				assert.NoError(t, err)
			}
			saverErr := fs.Save()
			assert.NoError(t, saverErr)
			fsToRestore := &File{
				path:        tt.fields.path,
				immediately: tt.fields.immediately,
				memStorage:  NewMemory(10),
			}
			if err := fsToRestore.Restore(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("File.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
			for k, v := range tt.want {
				g, err := fsToRestore.memStorage.GetGauge(ctx, k)
				assert.NoError(t, err)
				assert.Equal(t, v, float64(g))
			}

		})
	}
}
