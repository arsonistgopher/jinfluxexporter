package collector

import (
	channels "github.com/arsonistgopher/jinfluxexporter/rootchannels"
	"github.com/arsonistgopher/jinfluxexporter/rpc"
)

// RPCCollector collects metrics from JunOS using rpc.Client
type RPCCollector interface {

	// Collect collects metrics from JunOS
	Collect(client rpc.Client, ch chan<- channels.InfluxDBMeasurement, label string, topic string) error
}
