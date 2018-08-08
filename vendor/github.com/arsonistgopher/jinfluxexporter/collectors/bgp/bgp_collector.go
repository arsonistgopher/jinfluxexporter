package bgp

import (
	"time"

	"github.com/arsonistgopher/jinfluxexporter/collector"
	channels "github.com/arsonistgopher/jinfluxexporter/rootchannels"
	"github.com/arsonistgopher/jinfluxexporter/rpc"
)

type bgpCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &bgpCollector{}
}

// Collect collects metrics from JunOS
func (c *bgpCollector) Collect(client rpc.Client, ch chan<- channels.InfluxDBMeasurement, label string, measurement string) error {
	sessions, err := c.bgpSessions(client)
	if err != nil {
		return err
	}

	for _, s := range sessions {

		up := 0
		if s.Up {
			up = 1
		}

		tagset := make(map[string]string)
		tagset["host"] = label

		fieldset := map[string]interface{}{
			"up":               uint64(up),
			"receivedPrefixes": uint64(s.ReceivedPrefixes),
			"acceptedPrefixes": uint64(s.AcceptedPrefixes),
			"rejectedPrefixes": uint64(s.RejectedPrefixes),
			"activePrefixes":   uint64(s.ActivePrefixes),
			"inputMessages":    uint64(s.InputMessages),
			"outputMessages":   uint64(s.OutputMessages),
			"flaps":            uint64(s.Flaps),
		}

		ch <- channels.InfluxDBMeasurement{Measurement: measurement, TagSet: tagset, FieldSet: fieldset, TimeStamp: time.Now()}
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
