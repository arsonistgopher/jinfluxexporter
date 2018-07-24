# jKafkaexporter

[![GoDoc](https://godoc.org/github.com/arsonistgopher/jkafkaexporter?status.svg)](https://godoc.org/github.com/arsonistgopher/go-netconf/jkafkaexporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/arsonistgopher/jkafkaexporter)](https://goreportcard.com/report/github.com/arsonistgopher/jkafkaexporter)
[![Build Status](https://travis-ci.org/arsonistgopher/jkafkaexporter.png)](https://travis-ci.org/arsonistgopher/jkafkaexporter)

This is an application for Junos that drives NETCONF in order to retrieve data that is transformed into JSON and then placed on to a Kafka bus.

> **Note:** The idea here is that it is `easy` to create a new collector and re-build the binary.

## Release

Currently beta version. Seems to behave itself!

## Features
* Simple user land application 

## Install
* Requires Go 1.4 or later
* `go get github.com/arsonistgopher/jkafkaexporter`

## Documentation
You can view full API documentation at GoDoc: http://godoc.org/github.com/arsonistgopher/jkafkaexporter

## Usage

`./jkafkaexporter -identity vmx01 -kafkaperiod 1 -kafkaport 9092 -kafkahost localhost -password Passw0rd -username bob -sshport 22 -target 10.42.0.132`

You can also put an ampersand after the invocation to push the application in to the background.

Also, this application can also work with ssh-keys using the '-sshkey` switch which the argument is the fully qualified path to the ssh-private key. Ensure the public key has been installed on to the Junos device in question.

Currently, collectors include: alarm, environment, interfaces and routing-engine. With relative ease you should be able to create a custom collector based on the existing patterns in the aforementioned collectors. Include the import in `main.go` at both the package import section and relevant code section. See examples below:

```go
    import (
    ...
	// Add new collectors here
    "github.com/arsonistgopher/jkafkaexporter/something"
    )

    ...
    
    // And also add new collectors here...
	...
	c.Add("something", something.NewCollector())
```

In the last line of config, the "something" is also the Kafka topic.

## Attribution

Some code patterns were borrowed from [Daniel Czerwonk's Junos Prometheus exporter](https://github.com/czerwonk/junos_exporter) and as such his name is correctly present in the license. 

## Contributing

Welcoming contributions with arms open. Steps to contribute:

1.  Fork this repo using GitHub in to your own account
2.  Create changes on master and create tests ideally
3.  Document changes in comments and/or commit statements
4.  Add your name to the license. If the year exists in the license, append your name to the last person
5.  Create pull request
6.  Maintainer (currently just me) will comment and/or request changes
7.  Feel good and enjoy the kudos tokens (note> this is not a type of crypto-currency)

## Plans

I intend to create a library of collectors. Right now these are somewhat limited. Please create a collector using any of the existing ones as a base pattern and push back here!

## License

MIT