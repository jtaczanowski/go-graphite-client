package main

import (
	"log"

	"github.com/jtaczanowski/go-graphite-client"
)

func Example() {
	graphiteClient := graphite.NewClient("localhost", 2003, "metrics.prefix", "tcp")

	// metrics
	exampleMetric1 := map[string]float64{"test_metric": 1234.1234}
	exampleMetric2 := map[string]float64{"test_metric2": 12345.12345}
	// list of the metrics
	metricsToSend := []map[string]float64{exampleMetric1, exampleMetric2}

	// graphiteClient.SendData(data []map[string]float64) error - this method receives a list of metrics as an argument
	if err := graphiteClient.SendData(metricsToSend); err != nil {
		log.Printf("Error sending metrics: %v", err)
	}
}
