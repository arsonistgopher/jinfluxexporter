package collector

import (
	"github.com/arsonistgopher/jinfluxdbexporter/internal/channels"
	"github.com/arsonistgopher/jinfluxdbexporter/rpc"
)

// RPCCollector collects metrics from JunOS using rpc.Client
type RPCCollector interface {

	// Collect collects metrics from JunOS
	Collect(client rpc.Client, ch chan<- channels.InfluxDBMeasurement, label string, topic string) error
}
