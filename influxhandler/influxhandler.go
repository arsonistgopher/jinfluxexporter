package influxhandler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/arsonistgopher/jinfluxexporter/junoscollector"
	channels "github.com/arsonistgopher/jinfluxexporter/rootchannels"
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
func StartInflux(me string, ic Config, jc *junoscollector.JunosCollector, done chan struct{}, wg *sync.WaitGroup, timeout int) (chan channels.InfluxDBMeasurement, error) {

	// Create the buffer depth to match the number of collectors
	responsechan := make(chan channels.InfluxDBMeasurement, jc.Len())

	go func(me string, ic Config, jc *junoscollector.JunosCollector, done chan struct{}, wg *sync.WaitGroup, responsechan chan channels.InfluxDBMeasurement) {
		ticker := time.NewTicker(ic.InfluxdbExport)
		// To keep influxDB connection alive
		influxping := time.NewTicker(time.Duration(5 * time.Second))
		influxdbdeath := make(chan struct{}, 1)

		go func(responsechan chan channels.InfluxDBMeasurement, ic Config, done chan struct{}) {

			bp, err := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  ic.InfluxdbDB,
				Precision: "s",
			})

			if err != nil {
				log.Fatal(err)
			}

			transportTimeout := time.Duration(timeout+1) * time.Second
			// Create a new HTTPClient
			httpc, err := client.NewHTTPClient(client.HTTPConfig{
				Addr: ic.InfluxdbHost,
				// Here we set the transport timeout to be that of the influxexport period + 1 second, thus hopefully keeping things alive
				Timeout: transportTimeout,
			})
			if err != nil {
				log.Fatal(err)
			}

			for {
				select {
				case <-done:
					httpc.Close()
					if err != nil {
						log.Fatal(err)
					}
					wg.Done()
					return
				case r := <-responsechan:

					pt, err := client.NewPoint(r.Measurement, r.TagSet, r.FieldSet, r.TimeStamp)
					if err != nil {
						fmt.Print(err)
					}
					bp.AddPoint(pt)

					// Write the batch
					if err := httpc.Write(bp); err != nil {
						fmt.Print(err)
					}

				case <-influxping.C:
					// Do a InfluxDB ping every five seconds
					_, _, err := httpc.Ping(time.Second * 5)

					// If we get an error here, a TODO() would be to create a new client session and test it before crapping out
					if err != nil {
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
