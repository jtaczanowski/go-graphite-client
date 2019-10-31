# go-graphite-client [![Build Status](https://travis-ci.org/jtaczanowski/go-graphite-client.png?branch=master)](https://travis-ci.org/jtaczanowski/go-graphite-client) [![Coverage Status](https://coveralls.io/repos/github/jtaczanowski/go-graphite-client/badge.svg?branch=master)](https://coveralls.io/github/jtaczanowski/go-graphite-client?branch=master)
go-graphite-client - Simple Golang Graphite client which allows sending batches of metrics in single connection.

The optimal use of the library is to collect all the metrics in ```[]map[string]float64``` and after that pass it to SendData() method. SendData() method creates new connection to Graphite server and pushes all metric in **single** connection.

Example usage (also present in `example_text.go`)
```go
package main

import (
	"log"

	graphite "github.com/jtaczanowski/go-graphite-client"
)

func main() {
	graphiteClient := graphite.NewClient("localhost", 2003, "metrics.prefix", "tcp")

	// metrics
	exampleMetric1 := map[string]float64{"test_metric": 1234.1234}
	exampleMetric2 := map[string]float64{"test_metric2": 12345.12345}
	// list of the metrics
	metricsToSend := []map[string]float64{exampleMetric1, exampleMetric2}

	// graphiteClient.SendData(data []map[string]float64) error - this method receives a list of metrics as an argument
	// 
	if err := graphiteClient.SendData(metricsToSend); err != nil {
		log.Printf("Error sending metrics: %v", err)
	}
}
```
