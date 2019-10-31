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
