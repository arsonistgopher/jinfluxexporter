package kafka

import (
	"fmt"
	"sync"
	"time"

	"github.com/arsonistgopher/jkafkaexporter/junoscollector"
)

// Config holds the Kafka info
type Config struct {
	KafkaExport time.Duration // Number of seconds inbetween kafka exports
	KafkaHost   string        // Host IP or FQDN of kafka bus
	KafkaPort   int           // Port that kafka is running on
	KafkaTopic  string        // Topic for kafka export
}

// StartKafka is a GR that accepts a channel.
func StartKafka(me string, kc Config, jc *junoscollector.JunosCollector, done chan bool, wg *sync.WaitGroup) error {

	go func(me string, kc Config, jc *junoscollector.JunosCollector, done chan bool, wg *sync.WaitGroup) {
		for {
			ticker := time.NewTicker(kc.KafkaExport)
			responsechan := make(chan string, 10)

			select {
			case <-done:
				fmt.Println("DEBUG: Waiting for collector to exit()")
				wg.Done()
				return
			case <-ticker.C:
				// For each collector item, collect and dump
				jc.Collect(responsechan, me)
			case r := <-responsechan:
				// For now print TODO: Send to Kafka client
				fmt.Print(r)
			}
		}
	}(me, kc, jc, done, wg)

	return nil
}
