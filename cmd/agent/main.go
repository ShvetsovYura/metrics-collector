package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/exp/constraints"
)

const metiricsCount int = 40

type Sender interface {
	Send(string, string)
}
type Metrics map[string]Sender

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

func NewMetrics(metiricsCount int) Metrics {
	m := make(map[string]Sender, metiricsCount)
	m["PollCounter"] = counter(0)
	return m
}

func setGauge[Numeric constraints.Float | constraints.Integer](metrics Metrics, name string, v Numeric) {
	metrics[name] = gauge(v)
}

func increaseCounter(m Metrics) {
	val := m["PollCounter"]
	m["PollCounter"] = val.(counter) + counter(1)
}

func main() {
	m := NewMetrics(metiricsCount)
	opts := new(AgentOptions)
	opts.ParseArgs()
	if err := opts.ParseEnvs(); err != nil {
		log.Fatal(err.Error())
	}

	var elapsed int

	for {
		if elapsed > 0 {
			if elapsed%opts.PoolInterval == 0 {
				CollectMetrics(m)
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
	log.Println("start send metrics")
	for k, v := range m {
		v.Send(k, endpoint)
	}
}

func sendRequest(link string) error {
	r, err := http.Post(link, "text/html", nil)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

func CollectMetrics(m Metrics) {
	var rtm runtime.MemStats

	runtime.ReadMemStats(&rtm)

	setGauge(m, "HeapSys", rtm.HeapSys)
	setGauge(m, "Alloc", rtm.Alloc)
	setGauge(m, "BuckHashSys", rtm.BuckHashSys)
	setGauge(m, "Frees", rtm.Frees)
	setGauge(m, "GCCPUFraction", rtm.GCCPUFraction)
	setGauge(m, "GCSys", rtm.GCSys)
	setGauge(m, "HeapAlloc", rtm.HeapAlloc)
	setGauge(m, "HeapIdle", rtm.HeapIdle)
	setGauge(m, "HeapInuse", rtm.HeapInuse)
	setGauge(m, "HeapObjects", rtm.HeapObjects)
	setGauge(m, "HeapReleased", rtm.HeapReleased)
	setGauge(m, "LastGC", rtm.LastGC)
	setGauge(m, "Lookups", rtm.Lookups)
	setGauge(m, "MCacheInuse", rtm.MCacheInuse)
	setGauge(m, "MCacheSys", rtm.MCacheSys)
	setGauge(m, "MSpanInuse", rtm.MSpanInuse)
	setGauge(m, "MSpanSys", rtm.MSpanSys)
	setGauge(m, "Mallocs", rtm.Mallocs)
	setGauge(m, "NextGC", rtm.NextGC)
	setGauge(m, "NumForcedGC", rtm.NumForcedGC)
	setGauge(m, "NumGC", rtm.NumGC)
	setGauge(m, "OtherSys", rtm.OtherSys)
	setGauge(m, "PauseTotalNs", rtm.PauseTotalNs)
	setGauge(m, "StackInuse", rtm.StackInuse)
	setGauge(m, "StackSys", rtm.StackSys)
	setGauge(m, "Sys", rtm.Sys)
	setGauge(m, "TotalAlloc", rtm.TotalAlloc)
	setGauge(m, "RandomValue", rand.Float64())
	increaseCounter(m)
}
