package entity

import (
	"github.com/bcicen/ctop/cwidgets/compact"
	"github.com/bcicen/ctop/cwidgets"
	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/models"
)

type Service struct {
	models.Metrics
	Id        string
	Meta      map[string]string
	Widgets   *compact.Compact
	Display   bool // display this service in compact view
	updater   cwidgets.WidgetUpdater
	collector collector.Collector
}

func NewService(id string, collector collector.Collector) *Service {
	widgets := compact.NewCompact(id)
	return &Service{
		Metrics: 	models.NewMetrics(),
		Id: 	 	id,
		Meta:	 	make(map[string]string),
		Widgets: 	widgets,
		updater: 	widgets,
		collector:	collector,
	}
}

func (s *Service) Logs() collector.LogCollector {
	return s.collector.Logs()
}

func (s *Service) SetMeta(k, v string) {
	s.Meta[k] = v
	s.updater.SetMeta(k, v)
}

func (s *Service) GetMeta(k string) string {
	if v, ok := s.Meta[k]; ok {
		return v
	}
	return ""
}
