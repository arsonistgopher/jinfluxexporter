package bgp

import (
	"fmt"

	"github.com/arsonistgopher/jkafkaexporter/collector"
	"github.com/arsonistgopher/jkafkaexporter/internal/channels"
	"github.com/arsonistgopher/jkafkaexporter/rpc"
)

type bgpCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &bgpCollector{}
}

// Collect collects metrics from JunOS
func (c *bgpCollector) Collect(client rpc.Client, ch chan<- channels.Response, label string, topic string) error {
	sessions, err := c.bgpSessions(client)
	if err != nil {
		return err
	}

	for _, s := range sessions {

		up := 0
		if s.Up {
			up = 1
		}

		jsonResponse := "{Node: %s, up: %f, receivedPrefixes: %f, acceptedPrefixes: %f, rejectedPrefixes: %f, activePrefixes: %f, inputMessages: %f, outputMessages: %f, flaps: %f}"
		ch <- channels.Response{Data: fmt.Sprintf(jsonResponse, label, up, s.ReceivedPrefixes, s.AcceptedPrefixes, s.RejectedPrefixes, s.ActivePrefixes, s.InputMessages, s.OutputMessages, s.Flaps), Topic: topic}
	}

	return nil
}

func (c *bgpCollector) bgpSessions(client rpc.Client) ([]*BgpSession, error) {
	x := &BgpRpc{}
	err := rpc.RunCommandAndParse(client, "<get-bgp-summary-information/>", &x)
	if err != nil {
		return nil, err
	}

	sessions := make([]*BgpSession, 0)
	for _, peer := range x.Peers {
		s := &BgpSession{
			Ip:               peer.Ip,
			Up:               peer.State == "Established",
			Asn:              peer.Asn,
			Flaps:            float64(peer.Flaps),
			InputMessages:    float64(peer.InputMessages),
			OutputMessages:   float64(peer.OutputMessages),
			AcceptedPrefixes: float64(peer.Rib.AcceptedPrefixes),
			ActivePrefixes:   float64(peer.Rib.ActivePrefixes),
			ReceivedPrefixes: float64(peer.Rib.ReceivedPrefixes),
			RejectedPrefixes: float64(peer.Rib.RejectedPrefixes),
		}

		sessions = append(sessions, s)
	}

	return sessions, nil
}
