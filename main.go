package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"log"

	"github.com/arsonistgopher/jkafkaexporter/junoscollector"
	"github.com/arsonistgopher/jkafkaexporter/kafka"
)

const version string = "0.0.2"

var (
	showVersion = flag.Bool("version", false, "Print version information.")
	kafkaExport = flag.Int("kafkaperiod", 30, "Number of seconds inbetween kafka exports")
	kafkaHost   = flag.String("kafkastring", "127.0.0.1", "Host IP or FQDN of kafka bus")
	kafkaPort   = flag.Int("kafkaport", 3000, "Port that kafka is running on")
	kafkaTopic  = flag.String("kafkatopic", "nodeX", "Topic for kafka export")
)

func main() {
	flag.Parse()

	wg := &sync.WaitGroup{}

	// Setup kafkadeath channel
	kafkadeath := make(chan bool, 1)
	period := time.Duration(*kafkaExport * int(time.Second))

	// Build Kafka config from command line arguments
	kconfig := kafka.Config{
		KafkaExport: period,
		KafkaHost:   *kafkaHost,
		KafkaPort:   *kafkaPort,
		KafkaTopic:  *kafkaTopic,
	}

	// Create Junos collector system
	c := junoscollector.NewJunosCollector()

	// Start kafka GR that will consume the collector and transmit info to the topic
	err := kafka.StartKafka(kconfig, c, kafkadeath, wg)

	if err != nil {
		log.Printf("Error starting kafka: %s", err)
	}
	wg.Add(1)

	// Loop here now and wait for death signals
	// Create signal channel and register signals of interest
	sigs := make(chan os.Signal, 3)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// Create signal listener loop GR
	for {

		select {
		case c := <-sigs:

			if c == syscall.SIGINT || c == syscall.SIGTERM || c == syscall.SIGKILL {

				kafkadeath <- true
				wg.Wait()
				os.Exit(0)
			}
		}
	}
}
