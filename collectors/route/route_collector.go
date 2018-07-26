package route

import (
	"fmt"

	"github.com/arsonistgopher/jkafkaexporter/collector"
	"github.com/arsonistgopher/jkafkaexporter/internal/channels"
	"github.com/arsonistgopher/jkafkaexporter/rpc"
)

type routeCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &routeCollector{}
}

// Collect collects metrics from JunOS
func (c *routeCollector) Collect(client rpc.Client, ch chan<- channels.Response, label string, topic string) error {
	items, err := c.tableItems(client)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	for _, s := range items.Tables {
		jsonResponse := "{Node: %s, totalRoutes: %d, activeRoutes: %d, maxRoutes: %d}"
		ch <- channels.Response{Data: fmt.Sprintf(jsonResponse, label, s.TotalRoutes, s.ActiveRoutes, s.MaxRoutes), Topic: topic}

		for _, p := range s.Protocols {
			jsonResponse := "{Node: %s, protocolName: %s, protocolRoutes: %d, protocolActiveRoutes: %d}"
			ch <- channels.Response{Data: fmt.Sprintf(jsonResponse, label, p.Name, p.Routes, p.ActiveRoutes), Topic: topic}
		}
	}

	return nil
}

// Collect collects metrics from JunOS
func (c *routeCollector) tableItems(client rpc.Client) (RouteRpc, error) {
	x := &RouteRpc{}
	err := rpc.RunCommandAndParse(client, "<get-route-summary-information/>", &x)
	if err != nil {
		return *x, err
	}

	return *x, nil
}
