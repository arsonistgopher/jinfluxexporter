package kafka

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Shopify/sarama"
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
func StartKafka(me string, kc Config, jc *junoscollector.JunosCollector, done chan bool, wg *sync.WaitGroup) (chan string, error) {

	responsechan := make(chan string, 10)

	go func(me string, kc Config, jc *junoscollector.JunosCollector, done chan bool, wg *sync.WaitGroup, responsechan chan string) {
		ticker := time.NewTicker(kc.KafkaExport)
		kafkadeath := make(chan bool, 1)

		go func(responsechan chan string, kc Config, done chan bool) {

			config := sarama.NewConfig()
			config.Producer.Retry.Max = 5
			config.Producer.RequiredAcks = sarama.WaitForAll
			config.Producer.Return.Successes = true
			dialstring1 := fmt.Sprintf("%s:%d", kc.KafkaHost, kc.KafkaPort)
			brokers := []string{dialstring1}

			pd, err := sarama.NewSyncProducer(brokers, config)
			if err != nil {
				// Should not reach here
				panic(err)
			}

			defer func() {
				if err := pd.Close(); err != nil {
					// Should not reach here
					panic(err)
				}
			}()
			for {
				select {
				case <-done:
					wg.Done()
					return
				case r := <-responsechan:
					strTime := strconv.Itoa(int(time.Now().Unix()))
					msg := &sarama.ProducerMessage{
						Topic: kc.KafkaTopic,
						Key:   sarama.StringEncoder(strTime),
						Value: sarama.StringEncoder(r),
					}

					_, _, err = pd.SendMessage(msg)
					if err != nil {
						panic(err)
					}
				}
			}
		}(responsechan, kc, kafkadeath)

		for {
			select {
			case <-done:
				// Get's here, we're done
				kafkadeath <- true
				wg.Done()
				return
			case <-ticker.C:
				// For each collector item, collect and dump
				jc.Collect(responsechan, me)
			}
		}
	}(me, kc, jc, done, wg, responsechan)

	return responsechan, nil
}
