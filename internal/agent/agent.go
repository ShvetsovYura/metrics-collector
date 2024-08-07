package agent

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
)

// MetricItem: универсальная структура для данных для хранения единицы метрики
type MetricItem struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m MetricItem) MarshalJSON() ([]byte, error) {
	type MetricAlias MetricItem

	if m.MType == CounterTypeName {
		aliasValue := struct {
			MetricAlias
			Delta *int64 `json:"delta,omitempty"`
		}{
			MetricAlias: MetricAlias(m),
			Delta:       &m.Delta,
		}
		return json.Marshal(aliasValue)

	}
	if m.MType == GaugeTypeName {
		aliasValue := struct {
			MetricAlias
			ValuePtr *float64 `json:"value,omitempty"`
		}{
			MetricAlias: MetricAlias(m),
			ValuePtr:    &m.Value,
		}

		return json.Marshal(aliasValue)
	}

	return nil, errors.New("нет корретного типа метрики")
}

type Sender interface {
	Send(data []byte) error
}

// Agent: структура, для работы с метриками
type Agent struct {
	mx      sync.RWMutex
	metrics map[string]MetricItem
	options *Options
	sender  Sender
}

// NewAgent: инициализация нового экземляра агента сбора метрик
func NewAgent(metricsCount int, metricSender Sender, options *Options) *Agent {
	return &Agent{
		metrics: make(map[string]MetricItem, metricsCount),
		options: options,
		sender:  metricSender,
	}
}

// Run: запуск агента сбора метрик
func (a *Agent) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go a.runCollectMetrics(ctx, wg)
	go a.runSendMetrics(ctx, wg)

	logger.Log.Info("процессы агента запущены")
	wg.Wait()
}

func (a *Agent) runCollectMetrics(ctx context.Context, wg *sync.WaitGroup) {
	collectTicker := time.NewTicker(a.options.PollInterval)

	defer func() {
		collectTicker.Stop()
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("остановка сбора метрик...")
			return
		case <-collectTicker.C:
			logger.Log.Debug("collect metric")

			metricsCh := a.collectMetricsGenerator()
			addMetricsCh := a.collectAdditionalMetricsGenerator()
			mxCh := multiplexChannels(ctx, metricsCh, addMetricsCh)
			// если убрать блокировку - будет падать
			a.mx.Lock()
			for m := range mxCh {
				a.metrics[m.ID] = m
			}
			a.incrementCounter()
			a.mx.Unlock()
		}
	}
}

func (a *Agent) incrementCounter() {
	// TODO: Не нравится нижеследующий блок, подумать над изменением
	// дефолтное значение типа = 0
	var newVal int64
	// если такая метрика counter с этим именем уже есть в коллекции
	// то увеличиваем кол-во
	if v, ok := a.metrics[CounterFieldName]; ok {
		newVal = v.Delta + 1
	}
	// передаем 0 для новой counter-матрики или
	// записываем увеличенное значение для сущестующей
	a.metrics[CounterFieldName] = MetricItem{
		ID:    CounterFieldName,
		MType: CounterTypeName,
		Delta: newVal,
	}
}

func (a *Agent) runSendMetrics(ctx context.Context, wg *sync.WaitGroup) {
	var metricsBufferCh = make(chan MetricItem, 10)

	sendTicker := time.NewTicker(a.options.ReportInterval)

	defer func() {
		sendTicker.Stop()
		close(metricsBufferCh)
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("завершение процесса отправки...")
			return
		case <-sendTicker.C:
			logger.Log.Debug("start send")
			// если бы это было здесь, то залочило бы с размером канала = 10
			// a.mx.RLock()

			// for _, v := range a.metrics {
			// 	metricsBufferCh <- v
			// }
			// a.mx.RUnlock()

			var workers int

			if a.options.RateLimit == 0 {
				workers = len(a.metrics)
			} else {
				workers = a.options.RateLimit
			}

			for w := 0; w < workers; w++ {
				go a.senderWorker(metricsBufferCh)
			}
			a.mx.RLock()

			for _, v := range a.metrics {
				metricsBufferCh <- v
			}
			a.mx.RUnlock()
			logger.Log.Debug("end send")
		}
	}
}

func (a *Agent) senderWorker(metricsCh <-chan MetricItem) {
	for m := range metricsCh {
		data, err := json.Marshal(m)
		if err != nil {
			logger.Log.Error(err)
		} else {
			if err := a.sender.Send(data); err != nil {
				logger.Log.Warnf("не удалось отправить метрику: %s", data)
			}
		}
	}
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

func makeGaugeMetricItem(name string, val float64) MetricItem {
	return MetricItem{ID: name, MType: GaugeTypeName, Value: val}
}

func (a *Agent) collectMetricsGenerator() chan MetricItem {
	outCh := make(chan MetricItem)

	go func() {
		defer close(outCh)

		var rtm runtime.MemStats

		runtime.ReadMemStats(&rtm)
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
