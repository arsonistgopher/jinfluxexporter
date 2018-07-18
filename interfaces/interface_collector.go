package interfaces

import (
	"fmt"
	"strings"

	"github.com/arsonistgopher/jkafkaexporter/collector"
	"github.com/arsonistgopher/jkafkaexporter/rpc"
)

// Collector collects interface metrics
type interfaceCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	c := &interfaceCollector{}
	return c
}

// Collect collects metrics from JunOS
func (c *interfaceCollector) Collect(client rpc.Client, ch chan<- string, label string) error {
	stats, err := c.interfaceStats(client)
	if err != nil {
		return err
	}

	for _, s := range stats {
		c.collectForInterface(s, ch, label)
	}

	return nil
}

func (c *interfaceCollector) interfaceStats(client rpc.Client) ([]*InterfaceStats, error) {
	x := &InterfaceRpc{}
	err := rpc.RunCommandAndParse(client, `<get-interface-information><level>statistics</level><level-extra>detail</level-extra></get-interface-information>`, &x)
	if err != nil {
		return nil, err
	}

	stats := make([]*InterfaceStats, 0)
	for _, phy := range x.Interfaces {
		s := &InterfaceStats{
			IsPhysical:     true,
			Name:           strings.Replace(phy.Name, "\n", "", -1),
			AdminStatus:    phy.AdminStatus == "up",
			OperStatus:     phy.OperStatus == "up",
			ErrorStatus:    !(phy.AdminStatus == phy.OperStatus),
			Description:    strings.Replace(phy.Description, "\n", "", -1),
			Mac:            strings.Replace(phy.MacAddress, "\n", "", -1),
			ReceiveDrops:   float64(phy.InputErrors.Drops),
			ReceiveErrors:  float64(phy.InputErrors.Errors),
			ReceiveBytes:   float64(phy.Stats.InputBytes),
			TransmitDrops:  float64(phy.OutputErrors.Drops),
			TransmitErrors: float64(phy.OutputErrors.Errors),
			TransmitBytes:  float64(phy.Stats.OutputBytes),
		}

		stats = append(stats, s)

		for _, log := range phy.LogicalInterfaces {
			sl := &InterfaceStats{
				IsPhysical:    false,
				Name:          log.Name,
				Description:   log.Description,
				Mac:           phy.MacAddress,
				ReceiveBytes:  float64(log.Stats.InputBytes),
				TransmitBytes: float64(log.Stats.OutputBytes),
			}

			stats = append(stats, sl)
		}
	}

	return stats, nil
}

func (*interfaceCollector) collectForInterface(s *InterfaceStats, ch chan<- string, label string) {

	if s.IsPhysical {
		adminUp := 0
		if s.AdminStatus {
			adminUp = 1
		}
		operUp := 0
		if s.OperStatus {
			operUp = 1
		}
		err := 0
		if s.ErrorStatus {
			err = 1
		}

		jsonResponse := "{Node: %s, Iface: %s, IfaceDescription: %s, IfaceMAC: %s, ReceivedBytes: %f, " +
			"TransmitBytes: %f, AdminState: %d, OpState: %d, Err: %d, TXErr: %f, TXDrop: %f, RXError: %f, RXDrops: %f}"

		ch <- fmt.Sprintf(jsonResponse, label, s.Name, s.Description, s.Mac, s.ReceiveBytes, s.TransmitBytes,
			adminUp, operUp, err, s.TransmitErrors, s.TransmitDrops, s.ReceiveErrors, s.ReceiveDrops)
	}
}
