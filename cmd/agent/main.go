package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const baseUrl string = "http://localhost:8080/update"

type gauge float64
type counter int64

var gaugeFields = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
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
	"RandomValue",
}

var counterFields = []string{
	"PollCount",
}

type Metric struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge

	PollCount   counter
	RandomValue gauge
}

func main() {
	var m Metric
	var elapsed int

	poolInterval := 2
	reportInterval := 10

	for {
		if elapsed > 0 && elapsed%poolInterval == 0 {
			CollectMetrics(&m)
		}

		if elapsed > 0 && elapsed%reportInterval == 0 {
			SendMetrics(&m)
		}

		time.Sleep(time.Duration(1) * (time.Second))
		elapsed++
	}
}

func contains(s []string, v string) bool {
	for _, val := range s {
		if val == v {
			return true
		}
	}
	return false
}

func SendMetrics(m *Metric) {
	v := reflect.ValueOf(m).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fn := v.Type().Field(i).Name
		// ft := v.Field(i).Type().String()
		link, err := makeLink(fn, f.Interface())
		if err == nil {
			sendRequest(link)
		}
	}
}

func makeLink(mName string, v any) (string, error) {
	var mType string
	var val string

	_, ok := v.(gauge)
	val = fmt.Sprint(v)

	isGauge := contains(gaugeFields, mName)
	isCounter := contains(counterFields, mName)
	if ok && isGauge {
		mType = "gauge"
	} else {
		_, ok := v.(counter)
		if ok && isCounter {
			mType = "counter"
		} else {
			return "", errors.New("не удается получить тип метрики")
		}
	}

	if mType != "" {
		return fmt.Sprintf("%s/%s/%s/%s", baseUrl, mType, mName, val), nil
	} else {
		return "", errors.New("не удалось подготовить URL")
	}
}

func sendRequest(link string) {
	r, err := http.Post(link, "text/html", nil)
	if err != nil {
		return
	}
	defer r.Body.Close()
}

func CollectMetrics(m *Metric) {
	var rtm runtime.MemStats

	runtime.ReadMemStats(&rtm)

	m.HeapSys = gauge(rtm.HeapSys)
	m.Alloc = gauge(rtm.Alloc)
	m.BuckHashSys = gauge(rtm.BuckHashSys)
	m.Frees = gauge(rtm.Frees)
	m.GCCPUFraction = gauge(rtm.GCCPUFraction)
	m.GCSys = gauge(rtm.GCSys)
	m.HeapAlloc = gauge(rtm.HeapAlloc)
	m.HeapIdle = gauge(rtm.HeapIdle)
	m.HeapInuse = gauge(rtm.HeapInuse)
	m.HeapObjects = gauge(rtm.HeapObjects)
	m.HeapReleased = gauge(rtm.HeapReleased)
	m.HeapSys = gauge(rtm.HeapSys)
	m.LastGC = gauge(rtm.LastGC)
	m.Lookups = gauge(rtm.Lookups)
	m.MCacheInuse = gauge(rtm.MCacheInuse)
	m.MCacheSys = gauge(rtm.MCacheSys)
	m.MSpanInuse = gauge(rtm.MSpanInuse)
	m.MSpanSys = gauge(rtm.MSpanSys)
	m.Mallocs = gauge(rtm.Mallocs)
	m.NextGC = gauge(rtm.NextGC)
	m.NumForcedGC = gauge(rtm.NumForcedGC)
	m.NumGC = gauge(rtm.NumGC)
	m.OtherSys = gauge(rtm.OtherSys)
	m.PauseTotalNs = gauge(rtm.PauseTotalNs)
	m.StackInuse = gauge(rtm.StackInuse)
	m.StackSys = gauge(rtm.StackSys)
	m.Sys = gauge(rtm.Sys)
	m.TotalAlloc = gauge(rtm.TotalAlloc)
	m.RandomValue = gauge(rand.Float64())
	m.PollCount++

}
