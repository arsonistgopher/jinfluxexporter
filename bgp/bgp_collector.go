package bgp

import (
	"fmt"

	"github.com/arsonistgopher/jkafkaexporter/collector"
	"github.com/arsonistgopher/jkafkaexporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix string = "junos_bgp_seesion_"

var (
	upDesc               *prometheus.Desc
	receivedPrefixesDesc *prometheus.Desc
	acceptedPrefixesDesc *prometheus.Desc
	rejectedPrefixesDesc *prometheus.Desc
	activePrefixesDesc   *prometheus.Desc
	inputMessagesDesc    *prometheus.Desc
	outputMessagesDesc   *prometheus.Desc
	flapsDesc            *prometheus.Desc
)

func init() {
	l := []string{"target", "asn", "ip"}
	upDesc = prometheus.NewDesc(prefix+"up", "Session is up (1 = Established)", l, nil)
	receivedPrefixesDesc = prometheus.NewDesc(prefix+"prefixes_received_count", "Number of received prefixes", l, nil)
	acceptedPrefixesDesc = prometheus.NewDesc(prefix+"prefixes_accepted_count", "Number of accepted prefixes", l, nil)
	rejectedPrefixesDesc = prometheus.NewDesc(prefix+"prefixes_rejected_count", "Number of rejected prefixes", l, nil)
	activePrefixesDesc = prometheus.NewDesc(prefix+"prefixes_active_count", "Number of active prefixes (best route in RIB)", l, nil)
	inputMessagesDesc = prometheus.NewDesc(prefix+"messages_input_count", "Number of received messages", l, nil)
	outputMessagesDesc = prometheus.NewDesc(prefix+"messages_output_count", "Number of transmitted messages", l, nil)
	flapsDesc = prometheus.NewDesc(prefix+"flap_count", "Number of session flaps", l, nil)
}

type bgpCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &bgpCollector{}
}

// Collect collects metrics from JunOS
func (c *bgpCollector) Collect(client *rpc.Client, ch chan<- string, label string) error {
	sessions, err := c.bgpSessions(client)
	if err != nil {
		return err
	}

	for _, s := range sessions {
		c.collectForSession(s, ch, label)
	}

	return nil
}

func (c *bgpCollector) bgpSessions(client *rpc.Client) ([]*BgpSession, error) {
	x := &BgpRpc{}
	err := client.RunCommandAndParse("<get-bgp-summary-information/>", &x)
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

func (*bgpCollector) collectForSession(s *BgpSession, ch chan<- string, label string) {

	up := 0
	if s.Up {
		up = 1
	}

	jsonResponse := "{Node: %s, ASN: %s, IP: %s, Session: {Up: %s, ReceivedPrefixes: %s, AcceptedPrefixes: %s, RejectedPrefixes: %s, ActivePrefixes: %s, InputMessages: %s, OutputMessages: %s, Flaps %s}}"
	ch <- fmt.Sprintf(jsonResponse, label, s.Asn, s.Ip, string(up), s.ReceivedPrefixes, s.AcceptedPrefixes, s.RejectedPrefixes, s.ActivePrefixes, s.InputMessages, s.OutputMessages, s.Flaps)
}
