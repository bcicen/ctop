package entity

import (
	"github.com/bcicen/ctop/models"
	"github.com/bcicen/ctop/connector/collector"
)

type Task struct {
	models.Metrics
	Meta
	Id        string
	collector collector.Collector
}

func NewTask(id string, collector collector.Collector) *Task {
	return &Task{
		Metrics:   models.NewMetrics(),
		Meta:      NewMeta(id),
		Id:        id,
		collector: collector,
	}
}

func (t *Task) SetState(val string) {
	t.Meta.SetMeta("state", val)
	// start collector, if needed
	if val == "running" && !t.collector.Running() {
		t.collector.Start()
		//s.Read(s.collector.Stream())
	}
	// stop collector, if needed
	if val != "running" && t.collector.Running() {
		t.collector.Stop()
	}
}

func (t *Task) Logs() collector.LogCollector {
	return t.collector.Logs()
}

func (t *Task) GetMetaEntity() Meta {
	return t.Meta
}

func (t *Task) GetId() string {
	return t.Id
}

func (t *Task) GetMetrics() models.Metrics {
	return t.Metrics
}
