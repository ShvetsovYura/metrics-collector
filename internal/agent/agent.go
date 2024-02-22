package agent

import (
	"context"
	"encoding/json"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricItem struct {
	ID    string
	MType string
	Delta int64
	Value float64
}

type Agent struct {
	metrics map[string]MetricItem
	options *AgentOptions
}

func NewAgent(metricsCount int, options *AgentOptions) *Agent {
	return &Agent{
		metrics: make(map[string]MetricItem, metricsCount),
		options: options,
	}
}

func (a *Agent) Run(ctx context.Context) {
	sendTicker := time.NewTicker(time.Duration(a.options.ReportInterval) * time.Second)
	defer sendTicker.Stop()
	collectTicker := time.NewTicker(time.Duration(a.options.PoolInterval) * time.Second)
	defer collectTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-collectTicker.C:
			metricsCh := a.collectMetricsGenerator(ctx)
			addMetricsCh := a.collentAdditionalMetricsGenerator(ctx)
			go multipllexChannels(ctx, metricsCh, addMetricsCh)
			go a.processMetrics(metricsCh)
		case <-sendTicker.C:
			go a.sendMetrics(ctx)
		}
	}
}

func (a *Agent) sendMetrics(ctx context.Context) {
	logger.Log.Info("start send")
	var toSend = make(chan MetricItem, len(a.metrics))
	for _, v := range a.metrics {
		toSend <- v
	}
	if a.options.RateLimit == 0 {
		for w := 0; w < len(a.metrics); w++ {
			go a.senderWorker(toSend)
		}
	} else {
		for w := 0; w < a.options.RateLimit; w++ {
			go a.senderWorker(toSend)
		}

	}

	close(toSend)
	logger.Log.Info("end send")
}

func (a *Agent) senderWorker(items <-chan MetricItem) {
	logger.Log.Info("__worker")
	for m := range items {
		link := "http://" + a.options.EndpointAddr + "/update/"
		var data []byte
		if m.MType == GaugeTypeName {
			data, _ = json.Marshal(Metric{
				ID:    m.ID,
				MType: m.MType,
				Value: &m.Value,
			})
		}
		if m.MType == CounterTypeName {
			data, _ = json.Marshal(Metric{
				ID:    m.ID,
				MType: m.MType,
				Delta: &m.Delta,
			})
		}

		logger.Log.Info(string(data))
		sendMetric(data, link, DefaultContentType, a.options.Key)
	}
}

func (a Agent) sendMetricsBatch() error {
	logger.Log.Info("Strart send batch metrics")
	metricsBatch := make([]Metric, 0, len(a.metrics))
	for _, m := range a.metrics {
		var m_ Metric
		if m.MType == GaugeTypeName {
			m_ = Metric{
				ID:    m.ID,
				MType: m.MType,
				Value: &m.Value,
			}
		}
		if m.MType == CounterTypeName {
			m_ = Metric{
				ID:    m.ID,
				MType: m.MType,
				Delta: &m.Delta,
			}
		}
		metricsBatch = append(metricsBatch, m_)
	}
	link := "http://" + a.options.EndpointAddr + "/updates/"
	data, err := json.Marshal(metricsBatch)
	if err != nil {
		return err
	}
	sendMetric(data, link, DefaultContentType, a.options.Key)
	return nil
}

func (a *Agent) processMetrics(metricsCh <-chan MetricItem) {
	logger.Log.Info("start process metrics")
	for m := range metricsCh {
		a.metrics[m.ID] = m
	}
	var newVal int64 = 0
	if v, ok := a.metrics[CounterFieldName]; ok {
		newVal = v.Delta + 1
	}
	a.metrics[CounterFieldName] = MetricItem{
		ID:    CounterFieldName,
		MType: CounterTypeName,
		Delta: newVal,
		Value: -1,
	}
	logger.Log.Infof("end process metrics %v %v", len(metricsCh), len(a.metrics))

}
func makeGaugeMetricItem(name string, val float64) MetricItem {
	return MetricItem{ID: name, MType: GaugeTypeName, Value: val, Delta: -1}
}

func (a *Agent) collectMetricsGenerator(ctx context.Context) chan MetricItem {
	outCh := make(chan MetricItem)

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	go func() {
		defer close(outCh)
		logger.Log.Info("start collect metrics")
		outCh <- makeGaugeMetricItem("HeapSys", float64(rtm.HeapSys))
		outCh <- makeGaugeMetricItem("Alloc", float64(rtm.Alloc))
		outCh <- makeGaugeMetricItem("BuckHashSys", float64(rtm.BuckHashSys))
		outCh <- makeGaugeMetricItem("Frees", float64(rtm.Frees))
		outCh <- makeGaugeMetricItem("GCCPUFraction", rtm.GCCPUFraction)
		outCh <- makeGaugeMetricItem("GCSys", float64(rtm.GCSys))
		outCh <- makeGaugeMetricItem("HeapAlloc", float64(rtm.HeapAlloc))
		outCh <- makeGaugeMetricItem("HeapIdle", float64(rtm.HeapIdle))
		outCh <- makeGaugeMetricItem("HeapInuse", float64(rtm.HeapInuse))
		outCh <- makeGaugeMetricItem("HeapObjects", float64(rtm.HeapObjects))
		outCh <- makeGaugeMetricItem("HeapReleased", float64(rtm.HeapReleased))
		outCh <- makeGaugeMetricItem("LastGC", float64(rtm.LastGC))
		outCh <- makeGaugeMetricItem("Lookups", float64(rtm.Lookups))
		outCh <- makeGaugeMetricItem("MCacheInuse", float64(rtm.MCacheInuse))
		outCh <- makeGaugeMetricItem("MCacheSys", float64(rtm.MCacheSys))
		outCh <- makeGaugeMetricItem("MSpanInuse", float64(rtm.MSpanInuse))
		outCh <- makeGaugeMetricItem("MSpanSys", float64(rtm.MSpanSys))
		outCh <- makeGaugeMetricItem("Mallocs", float64(rtm.Mallocs))
		outCh <- makeGaugeMetricItem("NextGC", float64(rtm.NextGC))
		outCh <- makeGaugeMetricItem("NumForcedGC", float64(rtm.NumForcedGC))
		outCh <- makeGaugeMetricItem("NumGC", float64(rtm.NumGC))
		outCh <- makeGaugeMetricItem("OtherSys", float64(rtm.OtherSys))
		outCh <- makeGaugeMetricItem("PauseTotalNs", float64(rtm.PauseTotalNs))
		outCh <- makeGaugeMetricItem("StackInuse", float64(rtm.StackInuse))
		outCh <- makeGaugeMetricItem("StackSys", float64(rtm.StackSys))
		outCh <- makeGaugeMetricItem("Sys", float64(rtm.Sys))
		outCh <- makeGaugeMetricItem("TotalAlloc", float64(rtm.TotalAlloc))
		outCh <- makeGaugeMetricItem("RandomValue", rand.Float64())
		logger.Log.Info("end collect metrics")
	}()
	return outCh

}

func (a *Agent) collentAdditionalMetricsGenerator(ctx context.Context) chan MetricItem {
	var outCh = make(chan MetricItem)

	m, _ := mem.VirtualMemory()
	go func() {
		defer close(outCh)
		outCh <- makeGaugeMetricItem("TotalMemory", float64(m.Total))
		outCh <- makeGaugeMetricItem("FreeMemory", float64(m.Free))

		cpuUtilizations, _ := cpu.Percent(0, true)
		for i, c := range cpuUtilizations {
			outCh <- makeGaugeMetricItem("CPUutilization"+strconv.Itoa(i), c)
		}
	}()
	return outCh
}

func multipllexChannels(ctx context.Context, channels ...chan MetricItem) chan MetricItem {
	resultCh := make(chan MetricItem)
	wg := &sync.WaitGroup{}

	for _, ch := range channels {
		wg.Add(1)
		chClosure := ch
		go func() {
			defer wg.Done()
			for item := range chClosure {
				select {
				case <-ctx.Done():
					return
				case resultCh <- item:
				}
			}
		}()

		go func() {
			wg.Wait()
			if _, ok := <-resultCh; ok {
				close(resultCh)
			}
		}()

	}
	return resultCh
}
