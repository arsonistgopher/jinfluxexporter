package interfacediagnostics

import (
	"strconv"
	"time"

	"github.com/arsonistgopher/jinfluxexporter/collector"
	channels "github.com/arsonistgopher/jinfluxexporter/rootchannels"
	"github.com/arsonistgopher/jinfluxexporter/rpc"
)

type interfaceDiagnosticsCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &interfaceDiagnosticsCollector{}
}

// Collect collects metrics from JunOS
func (c *interfaceDiagnosticsCollector) Collect(client rpc.Client, ch chan<- channels.InfluxDBMeasurement, label string, measurement string) error {
	diagnostics, err := c.interfaceDiagnostics(client)
	if err != nil {
		return err
	}

	for _, d := range diagnostics {

		tagset := make(map[string]string)
		tagset["host"] = label

		fieldset := map[string]interface{}{
			"LaserBiasCurrent":    uint64(d.LaserBiasCurrent),
			"LaserOutputPower":    uint64(d.LaserOutputPower),
			"LaserOutputPowerDbm": uint64(d.LaserOutputPowerDbm),
			"ModuleTemperature":   uint64(d.ModuleTemperature),
		}

		ch <- channels.InfluxDBMeasurement{Measurement: measurement, TagSet: tagset, FieldSet: fieldset, TimeStamp: time.Now()}

		if d.ModuleVoltage > 0 {
			fieldset := map[string]interface{}{
				"ModuleVoltage":              uint64(d.ModuleVoltage),
				"RxSignalAvgOpticalPower":    uint64(d.RxSignalAvgOpticalPower),
				"RxSignalAvgOpticalPowerDbm": uint64(d.RxSignalAvgOpticalPowerDbm),
			}
			ch <- channels.InfluxDBMeasurement{Measurement: measurement, TagSet: tagset, FieldSet: fieldset, TimeStamp: time.Now()}
		} else {
			fieldset := map[string]interface{}{
				"LaserRxOpticalPower":    uint64(d.LaserRxOpticalPower),
				"LaserRxOpticalPowerDbm": uint64(d.LaserRxOpticalPowerDbm),
			}
			ch <- channels.InfluxDBMeasurement{Measurement: measurement, TagSet: tagset, FieldSet: fieldset, TimeStamp: time.Now()}
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
