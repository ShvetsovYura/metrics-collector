package agent

import (
	"encoding/json"
	"fmt"
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
	var del *int64
	var val *float64
	if m.MType == CounterTypeName {
		del = &m.Delta
	}
	if m.MType == GaugeTypeName {
		val = &m.Value
	}

	aliasValue := struct {
		MetricAlias
		DeltaPtr *int64   `json:"delta,omitempty"`
		ValuePtr *float64 `json:"value,omitempty"`
	}{
		MetricAlias: MetricAlias(m),
		DeltaPtr:    del,
		ValuePtr:    val,
	}
	jsonValue, err := json.Marshal(aliasValue)
	if err != nil {
		return nil, fmt.Errorf("ошбика при сериализации json %w", err)
	}
	return jsonValue, nil
}
