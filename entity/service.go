package entity

import (
	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/models"
)

type Service struct {
	models.Metrics
	Meta
	Id        string
	collector collector.Collector
}

func NewService(id string, collector collector.Collector) *Service {
	return &Service{
		Metrics:   models.NewMetrics(),
		Meta:      NewMeta(id),
		Id:        id,
		collector: collector,
	}
}

func (s *Service) SetState(val string) {
	s.Meta.SetMeta("state", val)
	// start collector, if needed
	if val == "running" && !s.collector.Running() {
		s.collector.Start()
		//s.Read(s.collector.Stream())
	}
	// stop collector, if needed
	if val != "running" && s.collector.Running() {
		s.collector.Stop()
	}
}

func (s *Service) Logs() collector.LogCollector {
	return s.collector.Logs()
}
