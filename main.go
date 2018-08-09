package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	//_ "net/http/pprof"

	"github.com/arsonistgopher/jinfluxexporter/influxhandler"
	"github.com/arsonistgopher/jinfluxexporter/junoscollector"
	"golang.org/x/crypto/ssh"

	// Add new collectors here
	"github.com/arsonistgopher/jinfluxexporter/collectors/alarm"
	"github.com/arsonistgopher/jinfluxexporter/collectors/bgp"
	"github.com/arsonistgopher/jinfluxexporter/collectors/environment"
	"github.com/arsonistgopher/jinfluxexporter/collectors/interfacediagnostics"
	"github.com/arsonistgopher/jinfluxexporter/collectors/interfaces"
	"github.com/arsonistgopher/jinfluxexporter/collectors/route"
	"github.com/arsonistgopher/jinfluxexporter/collectors/routingengine"
)

const version string = "00.01.00"

// Beta release 00.01.00

var (
	showVersion  = flag.Bool("version", false, "Print version information.")
	influxExport = flag.Int("influxperiod", 30, "Number of seconds in-between InfluxDB exports")
	influxHost   = flag.String("influxhost", "http://127.0.0.1:8086", "Host string in form http(s)://IP:PORT")
	influxDB     = flag.String("influxdb", "junos", "Database name")
	identity     = flag.String("identity", "vmx", "Identity of device targeted for data collection")
	username     = flag.String("username", "influx", "Username for NETCONF SSH connection")
	password     = flag.String("password", "Passw0rd", "Password for NETCONF SSH connection")
	port         = flag.Int("sshport", 22, "Port for NETCONF SSH connection")
	target       = flag.String("target", "127.0.0.1", "Host IP or FQDN of NETCONF server")
	sshkey       = flag.String("sshkey", "./id_rsa.pub", "Fully qualified path to SSH private key")
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

	runtime.GOMAXPROCS(2)
	// Parse the flags
	flag.Parse()

	wg := &sync.WaitGroup{}

	// Setup influxdeath channel
	influxdeath := make(chan struct{}, 1)
	period := time.Duration(int64(*influxExport) * int64(time.Second))

	// Build Influx config from command line arguments
	iconfig := influxhandler.Config{
		InfluxdbExport: period,
		InfluxdbHost:   *influxHost,
		InfluxdbDB:     *influxDB,
	}

	// Create an sshconfig empty type so we can conditionally populate it depending on the passed in SSH config
	var sshconfig *ssh.ClientConfig

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
	// Collector name is also the measurement name
	c := junoscollector.NewJunosCollector(sshconfig, *port, *target)
	c.Add("alarm", alarm.NewCollector(""))
	c.Add("interfaces", interfaces.NewCollector())
	c.Add("routing-engine", routingengine.NewCollector())
	c.Add("environment", environment.NewCollector())
	c.Add("route", route.NewCollector())
	c.Add("bgp", bgp.NewCollector())
	c.Add("interfacediagnostics", interfacediagnostics.NewCollector())

	// Add one to WaitGroup
	wg.Add(1)

	// Start Influx GR that will consume the collector and transmit info to the topic
	_, err := influxhandler.StartInflux(*identity, iconfig, c, influxdeath, wg, *influxExport)

	if err != nil {
		fmt.Printf("Error starting Influx handler: %s", err)
	}

	// Loop here now and wait for death signals
	// Create signal channel and register signals of interest
	sigs := make(chan os.Signal, 3)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// These lines are for "GOPS". Comment them out if you do not want to debug.
	//if err := agent.Listen(agent.Options{}); err != nil {
	//	log.Fatal(err)
	//}
	// End of "GOPS"

	// go func() {
	//	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	// Create signal listener loop GR
	for {

		select {
		case c := <-sigs:
			// fmt.Println("DEBUG: Received signal of some sort...")

			if c == syscall.SIGINT || c == syscall.SIGTERM || c == syscall.SIGKILL {

				influxdeath <- struct{}{}
				wg.Wait()
				os.Exit(0)
			}
		}
	}
}
