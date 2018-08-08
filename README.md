# jinfluxdbexporter

[![GoDoc](https://godoc.org/github.com/arsonistgopher/jinfluxdbexporter?status.svg)](https://godoc.org/github.com/arsonistgopher/go-netconf/jinfluxdbexporter)
[![Build Status](https://travis-ci.org/arsonistgopher/jinfluxdbexporter.svg?branch=master)](https://travis-ci.org/arsonistgopher/jinfluxdbexporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/arsonistgopher/jinfluxdbexporter)](https://goreportcard.com/report/github.com/arsonistgopher/jinfluxdbexporter)

This is an application for Junos that drives NETCONF in order to retrieve operational data and said data is inserted as measurements in to an InfluxDB database. Tags and fields are derived automatically from the measurements taken by this application. Your best bet is to use [Chronograf](https://www.influxdata.com/time-series-platform/chronograf/) to view the data being generated.

> **Note:** The idea here is that it is `easy` to create a new collector and re-build the binary.

## Release

Currently beta version. Seems to behave itself!

InfluxDB options are fairly limited, but the idea here is if we start with the basics and add features as they're required, the project will receive enhancements naturally.

## Features
* Simple user land application 

## Install
* Requires Go 1.4 or later
* `go get github.com/arsonistgopher/jinfluxexporter`

## Documentation
You can view full API documentation at GoDoc: http://godoc.org/github.com/arsonistgopher/jinfluxdbexporter.

## Usage

./jinfluxexporter -identity=vmx01 -influxdb=junos -influxhost=http://IP:PORT -influxperiod=5 -username=USERNAME -password=JUNOSPASSWORD -target=JUNOSHOST

You can also put an ampersand after the invocation to push the application in to the background. If you want to run this in the background or as a service, my advice here is to put it in a Docker container (see the Dockerfile).

Also, this application can also work with ssh-keys using the '-sshkey` switch which the argument is the fully qualified path to the ssh-private key. Ensure the public key has been installed on to the Junos device in question.

With relative ease you should be able to create a custom collector based on the existing patterns in the aforementioned collectors. Include the import in `main.go` at both the package import section and relevant code section. See examples below:

```go
    import (
    ...
	// Add new collectors here
    "github.com/arsonistgopher/jinfluxexporter/collectors/something"
    )

    ...
    
    // And also add new collectors here...
	...
	c.Add("something", something.NewCollector())
```

## Docker

In order to build the Docker image, follow the steps below. Feel free to modify as you see fit.

```bash
docker build -f Dockerfile . -t arsonistgopher/jinfluxexporter
docker run -d --name bob arsonistgopher/jinfluxexporter
```

You can also check Docker's logs
```bash
docker logs bob
```

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