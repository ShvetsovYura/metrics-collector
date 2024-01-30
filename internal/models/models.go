package models

type Metrics struct {
	ID    string   `json:"id" bson:"id"`                           // имя метрики
	MType string   `json:"type" bson:"type"`                       // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty" bson:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" bson:"value,omitempty"` // значение метрики в случае передачи gauge
}
