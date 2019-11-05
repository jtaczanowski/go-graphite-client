# go-graphite-client [![Build Status](https://travis-ci.org/jtaczanowski/go-graphite-client.png?branch=master)](https://travis-ci.org/jtaczanowski/go-graphite-client) [![Coverage Status](https://coveralls.io/repos/github/jtaczanowski/go-graphite-client/badge.svg?branch=master)](https://coveralls.io/github/jtaczanowski/go-graphite-client?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/jtaczanowski/go-graphite-client)](https://goreportcard.com/report/github.com/jtaczanowski/go-graphite-client)

go-graphite-client - Simple Golang Graphite client which allows sending batches of metrics in single connection.

The optimal use of the library is to collect set of metrics for current minute in single ```map[string]float64```
and after that pass it to `Client.SendData()` method.
`Client.SendData()` method creates a new connection to Graphite server **every time it's called**
and pushes all metric trough it.

Example usage (taken from `example_text.go`)
```go
package main

import (
	"log"

	graphite "github.com/jtaczanowski/go-graphite-client"
)

func Example() {
	graphiteClient := graphite.NewClient("localhost", 2003, "metrics.prefix", "tcp")

	// metrics map
	metricsMap := map[string]float64{
		"test_metric":  1234.1234,
		"test_metric2": 12345.12345,
	}

	// append metrics from function which returns map[string]float64 as well
	for k, v := range metricsGenerator() {
		metricsMap[k] = v
	}

	// graphiteClient.SendData(data map[string]float64) error - this method expects a map of metrics as an argument
	if err := graphiteClient.SendData(metricsMap); err != nil {
		log.Printf("Error sending metrics: %v", err)
	}
}

func metricsGenerator() map[string]float64 {
	return map[string]float64{
		"test_metric4": 3.14159265359,
	}
}

```
