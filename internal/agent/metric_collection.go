package agent

type MetricsCollection struct {
	metrics map[string]MetricItem
}

func (mc *MetricsCollection) SetItem(m MetricItem) {
	mc.metrics[m.ID] = m
}

func NewMetricCollector(initCountMetrics int) *MetricsCollection {
	return &MetricsCollection{metrics: make(map[string]MetricItem, initCountMetrics)}
}

func (mc *MetricsCollection) IncrementCounter() {
	// дефолтное значение типа = 0
	var newCounterVal int64
	// если такая метрика counter с этим именем уже существует - инкремет
	if v, ok := mc.metrics[CounterFieldName]; ok {
		newCounterVal = v.Delta + 1
	}
	// передаем 0 для новой counter-матрики или
	// записываем увеличенное значение для сущестующей
	mc.metrics[CounterFieldName] = MetricItem{
		ID:    CounterFieldName,
		MType: CounterTypeName,
		Delta: newCounterVal,
	}
}

// Count получить кол-во текущих метрик в коллекции
func (mc *MetricsCollection) Count() int {
	return len(mc.metrics)
}

// Items возращает функцию-итератор по элементам коллекции
func (mc *MetricsCollection) Items() func() (MetricItem, bool) {
	cursor := 0
	var keys = make([]string, 0, len(mc.metrics))
	for k := range mc.metrics {
		keys = append(keys, k)
	}

	fn := func() (MetricItem, bool) {
		if len(keys) < 1 {
			return MetricItem{}, false
		}
		if cursor <= len(mc.metrics) {
			cursor++
		}
		return mc.metrics[keys[cursor-1]], cursor < len(mc.metrics)
	}
	return fn
}
