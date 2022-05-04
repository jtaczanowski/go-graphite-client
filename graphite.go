// Package graphite provides simple Graphite client
package graphite

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

var timeNow = time.Now

const defaultTimeout = 3 * time.Second

// Client - struct with Graphite connection settings.
type Client struct {
	host     string
	port     int
	prefix   string
	protocol string
	timeOut  time.Duration
}

// NewClient - returns new Client with default connection timeout set to 3s.
func NewClient(Host string, Port int, Prefix string, Protocol string) *Client {
	return &Client{
		host:     Host,
		port:     Port,
		prefix:   Prefix,
		protocol: Protocol,
		timeOut:  defaultTimeout,
	}
}

// SendData - creates new connection to Graphite server and pushes
// provided batch of metrics in this single connection, thread-safe.
//
// Returns error in case of problems establishing, sending data or closing the connection
// (which should not be a problem with such short-lived connections).
//
// SendData receives as argument map[string]int64 where string is metric name,
// float64 is metric value, example:
//   map[string]float64{"test": 1234.1234, "test": 1234.1234}
func (g *Client) SendData(data map[string]float64) error {
	conn, err := net.DialTimeout(g.protocol, g.host+":"+strconv.Itoa(g.port), g.timeOut)
	if err != nil {
		return err
	}

	dataToSent := g.prepareData(data, timeNow().Unix())
	for _, str := range dataToSent {
		_, err = conn.Write([]byte(str))
		if err != nil {
			return err
		}
	}
	// it's safe to close connection here because
	// we are not exiting the function elsewhere after connection is open
	return conn.Close()
}

// SendDataWithTimeStamp - creates new connection to Graphite server and pushes
// provided batch of metrics in this single connection, thread-safe.
//
// Returns error in case of problems establishing, sending data or closing the connection
// (which should not be a problem with such short-lived connections).
//
// SendData receives as first argument map[string]int64 where string is metric name,
// float64 is metric value, example:
//   map[string]float64{"test": 1234.1234, "test": 1234.1234}
// and as a second argument Unix timestamp with which the metrics will be sent, example:
// timeNow().Unix()
func (g *Client) SendDataWithTimeStamp(data map[string]float64, timestamp int64) error {
	conn, err := net.DialTimeout(g.protocol, g.host+":"+strconv.Itoa(g.port), g.timeOut)
	if err != nil {
		return err
	}

	dataToSent := g.prepareData(data, timestamp)
	for _, str := range dataToSent {
		_, err = conn.Write([]byte(str))
		if err != nil {
			return err
		}
	}
	// it's safe to close connection here because
	// we are not exiting the function elsewhere after connection is open
	return conn.Close()
}

func (g *Client) prepareData(data map[string]float64, timestamp int64) []string {
	dataToGraphite := make([]string, 0)
	for metricName, metricVal := range data {
		dataToGraphite = append(dataToGraphite, fmt.Sprintf("%s.%s %f %d\n", g.prefix, metricName, metricVal, timestamp))
	}
	return dataToGraphite
}
