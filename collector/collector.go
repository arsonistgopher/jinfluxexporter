package collector

import (
	"github.com/arsonistgopher/jkafkaexporter/internal/channels"
	"github.com/arsonistgopher/jkafkaexporter/rpc"
)

// RPCCollector collects metrics from JunOS using rpc.Client
type RPCCollector interface {

	// Collect collects metrics from JunOS
	Collect(client rpc.Client, ch chan<- channels.Response, label string, topic string) error
}
