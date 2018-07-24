package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"log"

	"github.com/arsonistgopher/jkafkaexporter/junoscollector"
	"github.com/arsonistgopher/jkafkaexporter/kafka"
	"golang.org/x/crypto/ssh"

	// Add new collectors here
	"github.com/arsonistgopher/jkafkaexporter/collectors/alarm"
	"github.com/arsonistgopher/jkafkaexporter/collectors/environment"
	"github.com/arsonistgopher/jkafkaexporter/collectors/interfaces"
	"github.com/arsonistgopher/jkafkaexporter/collectors/routingengine"

	"github.com/google/gops/agent"
)

const version string = "00.01.00"

// Beta release 00.01.00

var (
	showVersion = flag.Bool("version", false, "Print version information.")
	kafkaExport = flag.Int("kafkaperiod", 30, "Number of seconds inbetween kafka exports")
	kafkaHost   = flag.String("kafkahost", "127.0.0.1", "Host IP or FQDN of kafka bus")
	kafkaPort   = flag.Int("kafkaport", 3000, "Port that kafka is running on")
	identity    = flag.String("identity", "vmx", "Topic for kafka export")
	username    = flag.String("username", "kafka", "Username for kafka NETCONF SSH connection")
	password    = flag.String("password", "kafka", "Password for kafka NETCONF SSH connection")
	port        = flag.Int("sshport", 22, "Port for kafka NETCONF SSH connection")
	target      = flag.String("target", "127.0.0.1", "Host IP or FQDN of NETCONF server")
	sshkey      = flag.String("sshkey", "./id_rsa.pub", "Fully qualified path to SSH private key")
)

// PublicKeyFile parses the SSH private key from a FQ file and returns an AuthMethod
// Function pattern taken from one of Svetlin Ralchev's blog posts
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func main() {
	// Parse the flags
	flag.Parse()

	wg := &sync.WaitGroup{}

	// Setup kafkadeath channel
	kafkadeath := make(chan bool, 1)
	period := time.Duration(int64(*kafkaExport) * int64(time.Second))

	// Build Kafka config from command line arguments
	kconfig := kafka.Config{
		KafkaExport: period,
		KafkaHost:   *kafkaHost,
		KafkaPort:   *kafkaPort,
	}

	// Create an sshconfig empty type so we can conditionally populate it depending on the passed in SSH config
	sshconfig := &ssh.ClientConfig{}

	if *sshkey != "" {
		sshconfig = &ssh.ClientConfig{
			User:            *username,
			Auth:            []ssh.AuthMethod{ssh.Password(*password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	} else {
		sshconfig = &ssh.ClientConfig{
			User: *username,
			Auth: []ssh.AuthMethod{
				PublicKeyFile(*sshkey)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	}

	// And also add new collectors here...
	// Collector name is also the Kafka topic
	c := junoscollector.NewJunosCollector(sshconfig, *port, *target)
	c.Add("alarm", alarm.NewCollector(""))
	c.Add("interfaces", interfaces.NewCollector())
	c.Add("routing-engine", routingengine.NewCollector())
	c.Add("environment", environment.NewCollector())

	// Add one to WaitGroup
	wg.Add(1)

	// Start kafka GR that will consume the collector and transmit info to the topic
	_, err := kafka.StartKafka(*identity, kconfig, c, kafkadeath, wg)

	if err != nil {
		log.Printf("Error starting kafka: %s", err)
	}

	// Loop here now and wait for death signals
	// Create signal channel and register signals of interest
	sigs := make(chan os.Signal, 3)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}

	// Create signal listener loop GR
	for {

		select {
		case c := <-sigs:
			// fmt.Println("DEBUG: Received signal of some sort...")

			if c == syscall.SIGINT || c == syscall.SIGTERM || c == syscall.SIGKILL {

				kafkadeath <- true
				// fmt.Println("DEBUG: Waiting for sync group to be done")
				wg.Wait()
				os.Exit(0)
			}
		}
	}
}
