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

func (n *Node) SetState(val string) {
	n.Meta.SetMeta("state", val)
	// start collector, if needed
	if val == "running" && !n.collector.Running() {
		n.collector.Start()
		//s.Read(s.collector.Stream())
	}
	// stop collector, if needed
	if val != "running" && n.collector.Running() {
		n.collector.Stop()
	}
}

func (n *Node) Logs() collector.LogCollector {
	return n.collector.Logs()
}

func (n *Node) GetMetaEntity() Meta {
	return n.Meta
}

func (n *Node) GetId() string {
	return n.Id
}

func (n *Node) GetMetrics() models.Metrics{
	return n.Metrics
}
