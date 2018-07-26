package interfacediagnostics

import (
	"fmt"
	"strconv"

	"github.com/arsonistgopher/jkafkaexporter/collector"
	"github.com/arsonistgopher/jkafkaexporter/internal/channels"
	"github.com/arsonistgopher/jkafkaexporter/rpc"
)

type interfaceDiagnosticsCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &interfaceDiagnosticsCollector{}
}

// Collect collects metrics from JunOS
func (c *interfaceDiagnosticsCollector) Collect(client rpc.Client, ch chan<- channels.Response, label string, topic string) error {
	diagnostics, err := c.interfaceDiagnostics(client)
	if err != nil {
		return err
	}

	for _, d := range diagnostics {

		jsonResponse1 := "{Node: %s, laserBiasCurrent: %f, laserOutputPower: %f, laserOutputPowerDbm: %f, moduleTemperature: %f}"
		jsonResponse2 := "{Node: %s, moduleVoltage: %f, rxSignalAvgOpticalPower: %f, rxSignalAvgOpticalPowerDbm: %f}"
		jsonResponse3 := "{Node: %s, laserRxOpticalPower: %f, laserRxOpticalPowerDbm: %f}"

		ch <- channels.Response{Data: fmt.Sprintf(jsonResponse1, label, d.LaserBiasCurrent, d.LaserOutputPower, d.LaserOutputPowerDbm, d.ModuleTemperature), Topic: topic}

		if d.ModuleVoltage > 0 {
			ch <- channels.Response{Data: fmt.Sprintf(jsonResponse2, label, d.ModuleVoltage, d.RxSignalAvgOpticalPower, d.RxSignalAvgOpticalPowerDbm), Topic: topic}
		} else {
			ch <- channels.Response{Data: fmt.Sprintf(jsonResponse3, label, d.LaserRxOpticalPower, d.LaserRxOpticalPowerDbm), Topic: topic}

		}
	}

	return nil
}

func (c *interfaceDiagnosticsCollector) interfaceDiagnostics(client rpc.Client) ([]*InterfaceDiagnostics, error) {
	x := &InterfaceDiagnosticsRpc{}
	err := rpc.RunCommandAndParse(client, "<get-interface-optics-diagnostics-information/>", &x)
	if err != nil {
		return nil, err
	}

	diagnostics := make([]*InterfaceDiagnostics, 0)
	for _, diag := range x.Diagnostics {
		if diag.Diagnostics.NA == "N/A" {
			continue
		}
		d := &InterfaceDiagnostics{
			Name:              diag.Name,
			LaserBiasCurrent:  float64(diag.Diagnostics.LaserBiasCurrent),
			LaserOutputPower:  float64(diag.Diagnostics.LaserOutputPower),
			ModuleTemperature: float64(diag.Diagnostics.ModuleTemperature.Value),
		}
		f, err := strconv.ParseFloat(diag.Diagnostics.LaserOutputPowerDbm, 64)
		if err == nil {
			d.LaserOutputPowerDbm = f
		}

		if diag.Diagnostics.ModuleVoltage > 0 {
			d.ModuleVoltage = float64(diag.Diagnostics.ModuleVoltage)
			d.RxSignalAvgOpticalPower = float64(diag.Diagnostics.RxSignalAvgOpticalPower)
			f, err = strconv.ParseFloat(diag.Diagnostics.RxSignalAvgOpticalPowerDbm, 64)
			if err == nil {
				d.RxSignalAvgOpticalPowerDbm = f
			}
		} else {
			d.LaserRxOpticalPower = float64(diag.Diagnostics.LaserRxOpticalPower)
			f, err = strconv.ParseFloat(diag.Diagnostics.LaserRxOpticalPowerDbm, 64)
			if err == nil {
				d.LaserRxOpticalPowerDbm = f
			}
		}

		diagnostics = append(diagnostics, d)
	}

	return diagnostics, nil
}
