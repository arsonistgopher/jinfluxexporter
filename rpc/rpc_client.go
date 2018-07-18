package rpc

import (
	"encoding/xml"
	"fmt"
	"log"

	drv "github.com/arsonistgopher/go-netconf/drivers/driver"
	sshdriver "github.com/arsonistgopher/go-netconf/drivers/ssh"
	"golang.org/x/crypto/ssh"
)

// Client sends commands to JunOS and parses results
type Client struct {
	D drv.Driver
}

// Create sets up the connection
func Create(sshconfig *ssh.ClientConfig, target string, port int) (*Client, error) {

	var err error

	c := Client{}

	d := drv.New(sshdriver.New())

	nc := d.(*sshdriver.DriverSSH)

	nc.Host = target
	nc.Port = port

	// Sort yourself out with SSH. Easiest to do that here.
	nc.SSHConfig = sshconfig

	c.D = nc

	err = c.D.Dial()

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
		return &c, err
	}

	return &c, nil
}

// Close closes the connection
func Close(c *Client) error {
	err := c.D.Close()
	return err
}

// RunCommandAndParse runs a command on JunOS and unmarshals the XML result
func RunCommandAndParse(c Client, rpcenv string, obj interface{}) error {
	var err error

	// reply, err := c.s.Exec(netconf.RawMethod(rpcenv))
	reply, err := c.D.SendRaw(rpcenv)

	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	err = xml.Unmarshal([]byte(reply.Data), obj)

	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	return nil
}
