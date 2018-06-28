package rpc

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/Juniper/go-netconf/netconf"
)

// Client sends commands to JunOS and parses results
type Client struct {
	s *netconf.Session
}

// Create sets up the connection
func Create() (*Client, error) {

	var err error
	c := &Client{}
	c.s, err = netconf.DialJunos()

	if err != nil {
		log.Fatal(err)
		return c, err
	}

	return c, nil
}

// Close closes the connection
func (c *Client) Close() error {
	c.s.Close()
	return nil
}

// RunCommandAndParse runs a command on JunOS and unmarshals the XML result
func (c *Client) RunCommandAndParse(rpcenv string, obj interface{}) error {
	var err error

	reply, err := c.s.Exec(netconf.RawMethod(rpcenv))
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
