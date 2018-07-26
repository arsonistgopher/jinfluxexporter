package route

import (
	"fmt"
	"time"

	"github.com/arsonistgopher/jinfluxdbexporter/collector"
	"github.com/arsonistgopher/jinfluxdbexporter/internal/channels"
	"github.com/arsonistgopher/jinfluxdbexporter/rpc"
)

type routeCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &routeCollector{}
}

// Collect collects metrics from JunOS
func (c *routeCollector) Collect(client rpc.Client, ch chan<- channels.InfluxDBMeasurement, label string, measurement string) error {
	items, err := c.tableItems(client)

	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	for _, s := range items.Tables {

		tagset := make(map[string]string)
		tagset["host"] = label
		tagset["table"] = s.Name

		fieldset := map[string]interface{}{
			"totalRoutes":  uint64(s.TotalRoutes),
			"activeRoutes": uint64(s.ActiveRoutes),
		}

		ch <- channels.InfluxDBMeasurement{Measurement: measurement, TagSet: tagset, FieldSet: fieldset, TimeStamp: time.Now()}

		for _, p := range s.Protocols {

			tagset := make(map[string]string)
			tagset["host"] = label
			tagset["table"] = s.Name
			tagset["protocol"] = p.Name

			fieldset := map[string]interface{}{
				"protocolName":         p.Name,
				"protocolRoutes":       uint64(p.Routes),
				"protocolActiveRoutes": uint64(p.ActiveRoutes),
			}

			ch <- channels.InfluxDBMeasurement{Measurement: measurement, TagSet: tagset, FieldSet: fieldset, TimeStamp: time.Now()}
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
