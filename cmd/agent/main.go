package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/ShvetsovYura/metrics-collector/internal/types"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

type gauge float64
type counter int64

func (g gauge) Send(mName string, baseURL string) {
	link := fmt.Sprintf("http://%s/update/gauge/%s/%f", baseURL, mName, g)
	sendRequest(link)
}

func (c counter) Send(mName string, baseURL string) {
	link := fmt.Sprintf("http://%s/update/counter/%s/%d", baseURL, mName, c)
	sendRequest(link)
}

type Metrics map[string]types.Sender

func NewMetrics() Metrics {
	m := make(map[string]types.Sender, 33)
	return m
}

func (m Metrics) SetGauge(name string, v any) {
	switch t := v.(type) {
	case uint32:
		m[name] = gauge(t)
	case uint64:
		m[name] = gauge(t)
	case float64:
		m[name] = gauge(t)
	}
}

func (m Metrics) SetCounter() {
	val := m["PollCounter"]
	if val == nil {
		m["PollCounter"] = counter(1)
	} else {
		m["PollCounter"] = val.(counter) + counter(1)
	}
}

func main() {
	m := NewMetrics()
	opts := new(util.AgentOptions)
	opts.ParseArgs()
	opts.ParseEnvs()
	var elapsed int

	for {
		if elapsed > 0 {
			if elapsed%opts.PoolInterval == 0 {
				CollectMetrics(&m)
			}
			if elapsed%opts.ReportInterval == 0 {
				SendMetrics(m, opts.EndpointAddr)
			}
		}

		time.Sleep(time.Duration(1) * (time.Second))
		elapsed++
	}
}

func SendMetrics(m Metrics, endpoint string) {
	fmt.Println("start send metrics")
	for k, v := range m {
		v.Send(k, endpoint)
	}
}

func sendRequest(link string) {
	r, err := http.Post(link, "text/html", nil)
	if err != nil {
		return
	}
	defer r.Body.Close()
}

func CollectMetrics(m *Metrics) {
	var rtm runtime.MemStats

	runtime.ReadMemStats(&rtm)

	m.SetGauge("HeapSys", rtm.HeapSys)
	m.SetGauge("Alloc", rtm.Alloc)
	m.SetGauge("BuckHashSys", rtm.BuckHashSys)
	m.SetGauge("Frees", rtm.Frees)
	m.SetGauge("GCCPUFraction", rtm.GCCPUFraction)
	m.SetGauge("GCSys", rtm.GCSys)
	m.SetGauge("HeapAlloc", rtm.HeapAlloc)
	m.SetGauge("HeapIdle", rtm.HeapIdle)
	m.SetGauge("HeapInuse", rtm.HeapInuse)
	m.SetGauge("HeapObjects", rtm.HeapObjects)
	m.SetGauge("HeapReleased", rtm.HeapReleased)
	m.SetGauge("LastGC", rtm.LastGC)
	m.SetGauge("Lookups", rtm.Lookups)
	m.SetGauge("MCacheInuse", rtm.MCacheInuse)
	m.SetGauge("MCacheSys", rtm.MCacheSys)
	m.SetGauge("MSpanInuse", rtm.MSpanInuse)
	m.SetGauge("MSpanSys", rtm.MSpanSys)
	m.SetGauge("Mallocs", rtm.Mallocs)
	m.SetGauge("NextGC", rtm.NextGC)
	m.SetGauge("NumForcedGC", rtm.NumForcedGC)
	m.SetGauge("NumGC", rtm.NumGC)
	m.SetGauge("OtherSys", rtm.OtherSys)
	m.SetGauge("PauseTotalNs", rtm.PauseTotalNs)
	m.SetGauge("StackInuse", rtm.StackInuse)
	m.SetGauge("StackSys", rtm.StackSys)
	m.SetGauge("Sys", rtm.Sys)
	m.SetGauge("TotalAlloc", rtm.TotalAlloc)
	m.SetGauge("RandomValue", rand.Float64())
	m.SetCounter()
}
