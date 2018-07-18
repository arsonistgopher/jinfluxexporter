# jKafkaexporter

[![GoDoc](https://godoc.org/github.com/arsonistgopher/jkafkaexporter?status.svg)](https://godoc.org/github.com/arsonistgopher/go-netconf/jkafkaexporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/arsonistgopher/jkafkaexporter)](https://goreportcard.com/report/github.com/arsonistgopher/jkafkaexporter)
[![Build Status](https://travis-ci.org/arsonistgopher/jkafkaexporter.png)](https://travis-ci.org/arsonistgopher/jkafkaexporter)

This is an application for Junos that drives NETCONF in order to retrieve data that is transformed into JSON and then placed on to a Kafka bus.

> **Note:** The idea here is that it is `easy` to create a new collector and re-build the binary.

## Features
* Simple user land application 

## Install
* Requires Go 1.4 or later
* `go get github.com/arsonistgopher/jkafkaexporter`

## Documentation
You can view full API documentation at GoDoc: http://godoc.org/github.com/arsonistgopher/jkafkaexporter

## Attribution

Some code patterns were borrowed from [Daniel Czerwonk's Junos Prometheus exporter](https://github.com/czerwonk/junos_exporter) and as such his name is correctly present in the license. 

## License

MIT

