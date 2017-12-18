package entity

import (
	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/models"
	"github.com/bcicen/ctop/cwidgets"
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
		s.collector.Start(s.Id)
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

func (s *Service) GetMetaEntity() Meta {
	return s.Meta
}

func (s *Service) GetId() string {
	return s.Id
}

func (s *Service) GetMetrics() models.Metrics{
	return s.Metrics
}

func (s *Service) GetMeta(v string) string {
	return s.Meta.GetMeta(v)
}

func (s *Service) SetUpdater(update cwidgets.WidgetUpdater) {
	s.Meta.SetUpdater(update)
}
