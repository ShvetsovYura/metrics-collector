package agent

import (
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

type Sender interface {
	Send(string, string)
}
type metrics map[string]Sender
type gauge float64
type counter int64

func NewMetrics(initSize int) metrics {
	m := make(map[string]Sender, initSize)
	m["PollCounter"] = counter(0)
	return m
}

func (g gauge) Send(mName string, baseURL string) {
	link := fmt.Sprintf("http://%s/update/gauge/%s/%f", baseURL, mName, g)
	util.SendRequest(link)
}

func (c counter) Send(mName string, baseURL string) {
	link := fmt.Sprintf("http://%s/update/counter/%s/%d", baseURL, mName, c)
	util.SendRequest(link)
}
