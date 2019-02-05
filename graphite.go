package graphite

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

var timeNow = time.Now

// GraphiteClient - struct with graphite connection settings
type GraphiteClient struct {
	host     string
	port     int
	prefix   string
	protocol string
}

// NewGraphiteClient - returns new NewGraphiteClient
func NewGraphiteClient(Host string, Port int, Prefix string, Protocol string) *GraphiteClient {
	return &GraphiteClient{
		host:     Host,
		port:     Port,
		prefix:   Prefix,
		protocol: Protocol,
	}
}

// SentData - pushes data to graphite server. Default connect timeout is set to 3s
// SentData receives as argument []map[string]int64 where string is metric name, float64 is metric value
// example: map[string]float64{"test": 1234.1234}
func (g *GraphiteClient) SendData(data []map[string]float64) error {
	dataToSent := g.prepareGraphiteData(data)
	conn, err := net.DialTimeout(g.protocol, g.host+":"+strconv.Itoa(g.port), time.Second*3)
	if err != nil {
		return err
	}
	defer conn.Close()
	for _, str := range dataToSent {
		conn.Write([]byte(str))
	}
	return nil
}

func (g *GraphiteClient) prepareGraphiteData(data []map[string]float64) []string {
	dataToGraphite := make([]string, 0)
	for _, metric := range data {
		for metricName, metricVal := range metric {
			dataToGraphite = append(dataToGraphite, fmt.Sprintf("%s.%s %f %d\n", g.prefix, metricName, metricVal, timeNow().Unix()))
		}
	}
	return dataToGraphite
}
