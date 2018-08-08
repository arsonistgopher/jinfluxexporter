package environment

import (
	"fmt"
	"time"

	channels "github.com/arsonistgopher/jinfluxexporter/rootchannels"

	"github.com/arsonistgopher/jinfluxexporter/collector"
	"github.com/arsonistgopher/jinfluxexporter/rpc"
)

type environmentCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	c := &environmentCollector{}
	return c
}

// Collect collects metrics from JunOS
func (c *environmentCollector) Collect(client rpc.Client, ch chan<- channels.InfluxDBMeasurement, label string, measurement string) error {
	items, err := c.environmentItems(client)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	for _, item := range items {

		tagset := make(map[string]string)
		tagset["host"] = label

		fieldset := map[string]interface{}{
			"receivedPrefixes": item.Name,
			"acceptedPrefixes": item.Temperature,
		}

		ch <- channels.InfluxDBMeasurement{Measurement: measurement, TagSet: tagset, FieldSet: fieldset, TimeStamp: time.Now()}
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
