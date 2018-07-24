package environment

import (
	"fmt"

	"github.com/arsonistgopher/jkafkaexporter/internal/channels"

	"github.com/arsonistgopher/jkafkaexporter/collector"
	"github.com/arsonistgopher/jkafkaexporter/rpc"
)

type environmentCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	c := &environmentCollector{}
	return c
}

// Collect collects metrics from JunOS
func (c *environmentCollector) Collect(client rpc.Client, ch chan<- channels.Response, label string, topic string) error {
	items, err := c.environmentItems(client)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	for _, item := range items {
		// fmt.Println("Ready to transmit environment data over channel...")
		jsonResponse := "{Node: %s, EnvironmentItem: %s, Temperature: %f}"
		ch <- channels.Response{Data: fmt.Sprintf(jsonResponse, label, item.Name, item.Temperature), Topic: topic}
		// fmt.Println("Transmitted data over channel...")
	}

	return nil
}

func (c *environmentCollector) environmentItems(client rpc.Client) ([]*EnvironmentItem, error) {
	x := &EnvironmentRpc{}
	err := rpc.RunCommandAndParse(client, "<get-environment-information/>", &x)
	if err != nil {
		return nil, err
	}

	// remove duplicates
	list := make(map[string]float64)
	for _, item := range x.Items {
		if item.Temperature != nil {
			list[item.Name] = float64(item.Temperature.Value)
		}
	}

	items := make([]*EnvironmentItem, 0)
	for name, value := range list {
		i := &EnvironmentItem{
			Name:        name,
			Temperature: value,
		}
		items = append(items, i)
	}

	return items, nil
}
