package influxhandler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/arsonistgopher/jinfluxdbexporter/internal/channels"
	"github.com/arsonistgopher/jinfluxdbexporter/junoscollector"
	client "github.com/influxdata/influxdb/client/v2"
)

// Config holds the influxdb info
type Config struct {
	InfluxdbExport time.Duration // Number of seconds in-between influxdb exports
	InfluxdbHost   string        // Host IP or FQDN of influxdb bus
	InfluxdbPort   int           // Port that influxdb is running on
	InfluxdbDB     string        // Database
}

// StartInflux is a GR that accepts a channel.
func StartInflux(me string, ic Config, jc *junoscollector.JunosCollector, done chan struct{}, wg *sync.WaitGroup) (chan channels.InfluxDBMeasurement, error) {

	// Create the buffer depth to match the number of collectors
	responsechan := make(chan channels.InfluxDBMeasurement, jc.Len())

	go func(me string, ic Config, jc *junoscollector.JunosCollector, done chan struct{}, wg *sync.WaitGroup, responsechan chan channels.InfluxDBMeasurement) {
		ticker := time.NewTicker(ic.InfluxdbExport)
		influxdbdeath := make(chan struct{}, 1)

		go func(responsechan chan channels.InfluxDBMeasurement, ic Config, done chan struct{}) {

			bp, err := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  ic.InfluxdbDB,
				Precision: "s",
			})

			if err != nil {
				log.Fatal(err)
			}

			for {
				select {
				case <-done:
					wg.Done()
					return
				case r := <-responsechan:

					// Create a new HTTPClient
					httpc, err := client.NewHTTPClient(client.HTTPConfig{
						Addr: ic.InfluxdbHost,
					})
					if err != nil {
						log.Fatal(err)
					}
					defer httpc.Close()

					pt, err := client.NewPoint(r.Measurement, r.TagSet, r.FieldSet, r.TimeStamp)
					if err != nil {
						fmt.Print(err)
					}
					bp.AddPoint(pt)

					// Write the batch
					if err := httpc.Write(bp); err != nil {
						fmt.Print(err)
					}
				}
			}
		}(responsechan, ic, influxdbdeath)

		for {
			select {
			case <-done:
				// Get's here, we're done
				influxdbdeath <- struct{}{}
				wg.Done()
				return
			case <-ticker.C:
				// For each collector item, collect and dump
				jc.Collect(responsechan, me)
			}
		}
	}(me, ic, jc, done, wg, responsechan)

	return responsechan, nil
}
