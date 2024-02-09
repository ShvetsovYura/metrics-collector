package agent

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"golang.org/x/exp/constraints"
)

type Agent struct {
	metrics metrics
	options *AgentOptions
}

func NewAgent(metricsCount int, options *AgentOptions) *Agent {
	return &Agent{
		metrics: NewMetrics(metricsCount),
		options: options,
	}
}

func (a *Agent) Run() {
	collectTicker := time.NewTicker(time.Duration(a.options.PoolInterval) * time.Second)
	defer collectTicker.Stop()
	sendTicker := time.NewTicker(time.Duration(a.options.ReportInterval) * time.Second)
	defer sendTicker.Stop()

	for {
		select {
		case <-collectTicker.C:
			a.collectMetrics()
		case <-sendTicker.C:
			// a.sendMetrics()
			a.sendMetricsBatch()
		}
	}
}

func (a *Agent) sendMetrics() {
	logger.Log.Info("Start send metrics")
	for k, v := range a.metrics {
		v.Send(k, a.options.EndpointAddr)
	}
}
func (a Agent) sendMetricsBatch() {
	logger.Log.Info("Strart send batch metrics")
	a.metrics.SendBatch(a.options.EndpointAddr)
}

func setGauge[Numeric constraints.Float | constraints.Integer](m metrics, name string, v Numeric) {
	m[name] = gauge(v)
}

func (a *Agent) setCounter() error {
	val := a.metrics[counterMetricFieldName]
	a.metrics["PollCount"] = val.(counter) + counter(1)
	return nil
}

func (a *Agent) collectMetrics() {
	var rtm runtime.MemStats

	runtime.ReadMemStats(&rtm)

	setGauge(a.metrics, "HeapSys", rtm.HeapSys)
	setGauge(a.metrics, "Alloc", rtm.Alloc)
	setGauge(a.metrics, "BuckHashSys", rtm.BuckHashSys)
	setGauge(a.metrics, "Frees", rtm.Frees)
	setGauge(a.metrics, "GCCPUFraction", rtm.GCCPUFraction)
	setGauge(a.metrics, "GCSys", rtm.GCSys)
	setGauge(a.metrics, "HeapAlloc", rtm.HeapAlloc)
	setGauge(a.metrics, "HeapIdle", rtm.HeapIdle)
	setGauge(a.metrics, "HeapInuse", rtm.HeapInuse)
	setGauge(a.metrics, "HeapObjects", rtm.HeapObjects)
	setGauge(a.metrics, "HeapReleased", rtm.HeapReleased)
	setGauge(a.metrics, "LastGC", rtm.LastGC)
	setGauge(a.metrics, "Lookups", rtm.Lookups)
	setGauge(a.metrics, "MCacheInuse", rtm.MCacheInuse)
	setGauge(a.metrics, "MCacheSys", rtm.MCacheSys)
	setGauge(a.metrics, "MSpanInuse", rtm.MSpanInuse)
	setGauge(a.metrics, "MSpanSys", rtm.MSpanSys)
	setGauge(a.metrics, "Mallocs", rtm.Mallocs)
	setGauge(a.metrics, "NextGC", rtm.NextGC)
	setGauge(a.metrics, "NumForcedGC", rtm.NumForcedGC)
	setGauge(a.metrics, "NumGC", rtm.NumGC)
	setGauge(a.metrics, "OtherSys", rtm.OtherSys)
	setGauge(a.metrics, "PauseTotalNs", rtm.PauseTotalNs)
	setGauge(a.metrics, "StackInuse", rtm.StackInuse)
	setGauge(a.metrics, "StackSys", rtm.StackSys)
	setGauge(a.metrics, "Sys", rtm.Sys)
	setGauge(a.metrics, "TotalAlloc", rtm.TotalAlloc)
	setGauge(a.metrics, "RandomValue", rand.Float64())
	a.setCounter()
	logger.Log.Info("Success collect metrics")
}
