package collector

import "github.com/arsonistgopher/jkafkaexporter/rpc"

// RPCCollector collects metrics from JunOS using rpc.Client
type RPCCollector interface {

	// Collect collects metrics from JunOS
	Collect(client rpc.Client, ch chan<- string, label string) error
}
