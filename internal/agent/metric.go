package agent

import (
	"encoding/json"
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal/types"
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
	// fmt.Println("start send gauge", mName)
	link := fmt.Sprintf("http://%s/update/", baseURL)
	val := float64(g)
	obj := types.Metrics{
		ID:    mName,
		MType: "gauge",
		Value: &val,
	}
	data, _ := json.Marshal(obj)

	util.SendRequest(link, "application/json", data)
}

func (c counter) Send(mName string, baseURL string) {
	// fmt.Println("start send counter")

	link := fmt.Sprintf("http://%s/update/", baseURL)
	val := int64(c)
	data, _ := json.Marshal(types.Metrics{
		ID:    "PollCounter",
		MType: "counter",
		Delta: &val,
	})
	util.SendRequest(link, "application/json", data)

}
