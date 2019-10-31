// Simple Graphite client
package graphite

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

var timeNow = time.Now

// Client - struct with Graphite connection settings
type Client struct {
	host     string
	port     int
	prefix   string
	protocol string
}

// NewClient - returns new Client
func NewClient(Host string, Port int, Prefix string, Protocol string) *Client {
	return &Client{
		host:     Host,
		port:     Port,
		prefix:   Prefix,
		protocol: Protocol,
	}
}

// SendData - creates new connection to Graphite server  and pushes batch of metrics in this single connection. Default connect timeout is set to 3s.
//
// SendData receives as argument []map[string]int64 where string is metric name, float64 is metric value, example:
//   map[string]float64{"test": 1234.1234}
func (g *Client) SendData(data []map[string]float64) error {
	dataToSent := g.prepareData(data)
	conn, err := net.DialTimeout(g.protocol, g.host+":"+strconv.Itoa(g.port), time.Second*3)
	if err != nil {
		return err
	}
	defer conn.Close()
	for _, str := range dataToSent {
		_, _ = conn.Write([]byte(str))
	}
	return nil
}

func (g *Client) prepareData(data []map[string]float64) []string {
	dataToGraphite := make([]string, 0)
	for _, metric := range data {
		for metricName, metricVal := range metric {
			dataToGraphite = append(dataToGraphite, fmt.Sprintf("%s.%s %f %d\n", g.prefix, metricName, metricVal, timeNow().Unix()))
		}
	}
	return dataToGraphite
}
