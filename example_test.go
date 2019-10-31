package graphite_test

import (
	"log"

	"github.com/jtaczanowski/go-graphite-client"
)

func Example() {
	graphiteClient := graphite.NewClient("localhost", 2003, "metrics.prefix", "tcp")

	// metrics map
	metricsMap := map[string]float64{
		"test_metric":  1234.1234,
		"test_metric2": 12345.12345,
	}

	// graphiteClient.SendData(data map[string]float64) error - this method receives a map of metrics as an argument
	if err := graphiteClient.SendData(metricsMap); err != nil {
		log.Printf("Error sending metrics: %v", err)
	}
}
