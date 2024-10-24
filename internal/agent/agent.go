package agent

import (
	"context"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
)

func MakeGaugeMetricItem(name string, val float64) MetricItem {
	return MetricItem{ID: name, MType: GaugeTypeName, Value: val}
}

type Sender interface {
	Send(item MetricItem) error
}
type Setter interface {
	SetItem(m MetricItem)
}
type Incrementer interface {
	IncrementCounter()
}

type Storer interface {
	Setter
	Incrementer
	Count() int
	Items() func() (MetricItem, bool)
}

// Agent: структура, для работы с метриками
type Agent struct {
	mx         sync.RWMutex
	collection Storer
	options    *Options
	sender     Sender
}

// NewAgent: инициализация нового экземляра агента сбора метрик
func NewAgent(metricCollector Storer, metricSender Sender, options *Options) *Agent {
	return &Agent{
		collection: metricCollector,
		options:    options,
		sender:     metricSender,
	}
}

// Run: запуск агента сбора метрик
func (a *Agent) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go a.runCollectMetrics(ctx, wg) // сбор метрик
	go a.runSendMetrics(ctx, wg)    // отправка метрик

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
			// так как в другой горутине читаем
			// из этой мапы
			// a.mx.Lock()
			for m := range mxCh {
				a.collection.SetItem(m)
			}
			a.collection.IncrementCounter()
			// a.mx.Unlock()
		}
	}
}

func (a *Agent) runSendMetrics(ctx context.Context, wg *sync.WaitGroup) {
	var metricsBufferCh = make(chan MetricItem, a.options.RateLimit)

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

			var workers int

			if a.options.RateLimit == 0 {
				if a.collection.Count() < 1 {
					continue
				}
				workers = a.collection.Count()
			} else {
				workers = a.options.RateLimit
			}

			for w := 0; w < workers; w++ {
				go a.senderWorker(metricsBufferCh)
			}
			a.mx.RLock()

			next := a.collection.Items()
			for {
				val, hasNext := next()
				metricsBufferCh <- val
				if !hasNext {
					break
				}
			}

			a.mx.RUnlock()
			logger.Log.Debug("end send")
		}
	}
}

func (a *Agent) senderWorker(metricsCh <-chan MetricItem) {
	for m := range metricsCh {
		if err := a.sender.Send(m); err != nil {
			logger.Log.Warnf("не удалось отправить метрику: %s", m)
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

func (a *Agent) collectMetricsGenerator() chan MetricItem {
	outCh := make(chan MetricItem)

	go func() {
		defer close(outCh)

		var rtm runtime.MemStats

		runtime.ReadMemStats(&rtm)
		outCh <- MakeGaugeMetricItem("HeapSys", float64(rtm.HeapSys))
		outCh <- MakeGaugeMetricItem("Alloc", float64(rtm.Alloc))
		outCh <- MakeGaugeMetricItem("BuckHashSys", float64(rtm.BuckHashSys))
		outCh <- MakeGaugeMetricItem("Frees", float64(rtm.Frees))
		outCh <- MakeGaugeMetricItem("GCCPUFraction", rtm.GCCPUFraction)
		outCh <- MakeGaugeMetricItem("GCSys", float64(rtm.GCSys))
		outCh <- MakeGaugeMetricItem("HeapAlloc", float64(rtm.HeapAlloc))
		outCh <- MakeGaugeMetricItem("HeapIdle", float64(rtm.HeapIdle))
		outCh <- MakeGaugeMetricItem("HeapInuse", float64(rtm.HeapInuse))
		outCh <- MakeGaugeMetricItem("HeapObjects", float64(rtm.HeapObjects))
		outCh <- MakeGaugeMetricItem("HeapReleased", float64(rtm.HeapReleased))
		outCh <- MakeGaugeMetricItem("LastGC", float64(rtm.LastGC))
		outCh <- MakeGaugeMetricItem("Lookups", float64(rtm.Lookups))
		outCh <- MakeGaugeMetricItem("MCacheInuse", float64(rtm.MCacheInuse))
		outCh <- MakeGaugeMetricItem("MCacheSys", float64(rtm.MCacheSys))
		outCh <- MakeGaugeMetricItem("MSpanInuse", float64(rtm.MSpanInuse))
		outCh <- MakeGaugeMetricItem("MSpanSys", float64(rtm.MSpanSys))
		outCh <- MakeGaugeMetricItem("Mallocs", float64(rtm.Mallocs))
		outCh <- MakeGaugeMetricItem("NextGC", float64(rtm.NextGC))
		outCh <- MakeGaugeMetricItem("NumForcedGC", float64(rtm.NumForcedGC))
		outCh <- MakeGaugeMetricItem("NumGC", float64(rtm.NumGC))
		outCh <- MakeGaugeMetricItem("OtherSys", float64(rtm.OtherSys))
		outCh <- MakeGaugeMetricItem("PauseTotalNs", float64(rtm.PauseTotalNs))
		outCh <- MakeGaugeMetricItem("StackInuse", float64(rtm.StackInuse))
		outCh <- MakeGaugeMetricItem("StackSys", float64(rtm.StackSys))
		outCh <- MakeGaugeMetricItem("Sys", float64(rtm.Sys))
		outCh <- MakeGaugeMetricItem("TotalAlloc", float64(rtm.TotalAlloc))
		outCh <- MakeGaugeMetricItem("RandomValue", rand.Float64())
	}()

	return outCh
}

func (a *Agent) collectAdditionalMetricsGenerator() chan MetricItem {
	var outCh = make(chan MetricItem)

	go func() {
		m, _ := mem.VirtualMemory()

		defer close(outCh)
		outCh <- MakeGaugeMetricItem("TotalMemory", float64(m.Total))
		outCh <- MakeGaugeMetricItem("FreeMemory", float64(m.Free))

		cpuUtilizations, _ := cpu.Percent(0, true)
		for i, c := range cpuUtilizations {
			outCh <- MakeGaugeMetricItem("CPUutilization"+strconv.Itoa(i), c)
		}
	}()

	return outCh
}
