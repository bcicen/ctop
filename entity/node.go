package entity

import (
	"github.com/bcicen/ctop/models"
	"github.com/bcicen/ctop/connector/collector"
)

type Node struct {
	models.Metrics
	Meta
	Id        string
	collector collector.Collector
}

func NewNode(id string, collector collector.Collector) *Node {
	return &Node{
		Metrics:   models.NewMetrics(),
		Meta:      NewMeta(id),
		Id:        id,
		collector: collector,
	}
}

func (s *Node) SetState(val string) {
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

func (s *Node) Logs() collector.LogCollector {
	return s.collector.Logs()
}
