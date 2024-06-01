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
	mx      sync.RWMutex
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
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go a.runCollectMetrics(ctx, wg)
	go a.runSendMetrics(ctx, wg)

	wg.Wait()
	logger.Log.Info("end agent app")
}

func (a *Agent) runCollectMetrics(ctx context.Context, wg *sync.WaitGroup) {
	collectTicker := time.NewTicker(time.Duration(a.options.PoolInterval) * time.Second)
	defer collectTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case <-collectTicker.C:
			processWaiter := &sync.WaitGroup{}
			processWaiter.Add(1)
			metricsCh := a.collectMetricsGenerator()
			addMetricsCh := a.collectAdditionalMetricsGenerator()
			allMetricsCh := multiplexChannels(ctx, metricsCh, addMetricsCh)
			go a.processMetrics(processWaiter, allMetricsCh)
			wg.Wait()
		}
	}
}

func (a *Agent) runSendMetrics(ctx context.Context, wg *sync.WaitGroup) {
	var toSend = make(chan MetricItem, 100)
	sendTicker := time.NewTicker(time.Duration(a.options.ReportInterval) * time.Second)
	defer sendTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(toSend)
			wg.Done()
			return
		case <-sendTicker.C:
			logger.Log.Info("start send")
			a.mx.RLock()
			for _, v := range a.metrics {
				toSend <- v
			}
			a.mx.RUnlock()
			var workers int
			if a.options.RateLimit == 0 {
				workers = len(a.metrics)
			} else {
				workers = a.options.RateLimit
			}
			for w := 0; w < workers; w++ {
				go a.senderWorker(toSend)
			}
			logger.Log.Info("end send")
		}
	}
}

func (a *Agent) senderWorker(items <-chan MetricItem) {
	for m := range items {
		link := "http://" + a.options.EndpointAddr + "/update/"
		var data []byte
		if m.MType == GaugeTypeName {
			data, _ = json.Marshal(Metric{
				ID:    m.ID,
				MType: m.MType,
				// если не буду ссылаться на поле другой структуры, то всегда по-ссылке
				// будет последнее значение gauge и для всех одинаковое
				// поэтому пришлось сделать еще одну структуру в начале файла,
				// но с полями-значениями, а не ссылками
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

		sendMetric(data, link, DefaultContentType, a.options.Key)
	}
}

func (a *Agent) processMetrics(wg *sync.WaitGroup, metricsCh <-chan MetricItem) {
	defer wg.Done()
	a.mx.Lock()
	for m := range metricsCh {
		a.metrics[m.ID] = m
	}

	var newVal int64
	if v, ok := a.metrics[CounterFieldName]; ok {
		newVal = v.Delta + 1
	}
	a.metrics[CounterFieldName] = MetricItem{
		ID:    CounterFieldName,
		MType: CounterTypeName,
		Delta: newVal,
	}
	a.mx.Unlock()
}

func makeGaugeMetricItem(name string, val float64) MetricItem {
	return MetricItem{ID: name, MType: GaugeTypeName, Value: val}
}

func (a *Agent) collectMetricsGenerator() chan MetricItem {
	outCh := make(chan MetricItem)

	go func() {
		var rtm runtime.MemStats
		runtime.ReadMemStats(&rtm)
		defer close(outCh)
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
	}()
	return outCh

}

func (a *Agent) collectAdditionalMetricsGenerator() chan MetricItem {
	var outCh = make(chan MetricItem)

	go func() {
		m, _ := mem.VirtualMemory()
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

func multiplexChannels(ctx context.Context, channels ...chan MetricItem) chan MetricItem {
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

	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	return resultCh
}
