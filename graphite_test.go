package graphite

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"testing"
	"time"
)

func init() {
	timeNow = func() time.Time {
		t, _ := time.Parse("2006-01-02 15:04:05", "2019-01-20 20:20:20")
		return t
	}
}

// Below init function
func TestSentMetricsOverTCP(t *testing.T) {
	expected1 := "prefix.test 1234.123400 1548015620\n"
	expected2 := "prefix.test2 12345.123450 1548015620\n"
	receivedMessage := make(chan string, 1)
	serverStarted := make(chan string, 1)
	port := rand.Intn(40000) + 10000

	go func() {
		// start tcp server
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			t.Fatal(err)
		}
		serverStarted <- "started"
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			defer conn.Close()

			buf, err := ioutil.ReadAll(conn)
			if err != nil {
				t.Fatal(err)
			}
			receivedMessage <- string(buf[:])
		}
	}()
	<-serverStarted

	// create graphite client and sent metrics in separate gorutine
	graphiteClient := NewClient("127.0.0.1", port, "prefix", "tcp")
	metricsMap := map[string]float64{
		"test":  1234.1234,
		"test2": 12345.12345,
	}
	err := graphiteClient.SendData(metricsMap)

	if msg := <-receivedMessage; msg != expected1+expected2 && msg != expected2+expected1 {
		t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, expected1+expected2)
	}

	if err != nil {
		t.Fatalf("Unexpected error sending metrics:\nGot:\t\t%s\nExpected:\tnil\n", err)
	}
}

func TestSentMetricsWithTimeStampOverTCP(t *testing.T) {
	expected1 := "prefix.test 1234.123400 1234\n"
	expected2 := "prefix.test2 12345.123450 1234\n"
	receivedMessage := make(chan string, 1)
	serverStarted := make(chan string, 1)
	port := rand.Intn(40000) + 10000

	go func() {
		// start tcp server
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			t.Fatal(err)
		}
		serverStarted <- "started"
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			defer conn.Close()

			buf, err := ioutil.ReadAll(conn)
			if err != nil {
				t.Fatal(err)
			}
			receivedMessage <- string(buf[:])
		}
	}()
	<-serverStarted

	// create graphite client and sent metrics in separate gorutine
	graphiteClient := NewClient("127.0.0.1", port, "prefix", "tcp")
	metricsMap := map[string]float64{
		"test":  1234.1234,
		"test2": 12345.12345,
	}
	err := graphiteClient.SendDataWithTimeStamp(metricsMap, 1234)

	if msg := <-receivedMessage; msg != expected1+expected2 && msg != expected2+expected1 {
		t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, expected1+expected2)
	}

	if err != nil {
		t.Fatalf("Unexpected error sending metrics:\nGot:\t\t%s\nExpected:\tnil\n", err)
	}
}
func TestSentMetricsWithTimeStampOverUDP(t *testing.T) {
	expectedMessages := []string{"prefix.test 1234.123400 1234\n", "prefix.test2 12345.123450 1234\n"}
	receivedMessage := make(chan string, 2)
	serverStarted := make(chan string, 1)
	port := rand.Intn(40000) + 10000

	go func() {
		// start UDP server
		listener, err := net.ListenPacket("udp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			t.Fatal(err)
		}
		serverStarted <- "started"
		defer listener.Close()
		for {
			buf := make([]byte, 1024)
			n, _, err := listener.ReadFrom(buf)
			if err != nil {
				t.Fatal(err)
			}
			receivedMessage <- string(buf[:n])
		}
	}()
	<-serverStarted

	// create graphite Client and sent two metrics
	graphiteClient := NewClient("127.0.0.1", port, "prefix", "udp")
	metricsMap := map[string]float64{
		"test":  1234.1234,
		"test2": 12345.12345,
	}
	err := graphiteClient.SendDataWithTimeStamp(metricsMap, 1234)

	for i := 1; i <= 2; i++ {
		msg := <-receivedMessage
		found := false
		for _, n := range expectedMessages {
			if n == msg {
				found = true
			}
		}
		if !found {
			t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, expectedMessages)
		}
	}
	if err != nil {
		t.Fatalf("Unexpected error sending metrics:\nGot:\t\t%s\nExpected:\tnil\n", err)
	}
}

func TestSentMetricsOverUDP(t *testing.T) {
	expectedMessages := []string{"prefix.test 1234.123400 1548015620\n", "prefix.test2 12345.123450 1548015620\n"}
	receivedMessage := make(chan string, 2)
	serverStarted := make(chan string, 1)
	port := rand.Intn(40000) + 10000

	go func() {
		// start UDP server
		listener, err := net.ListenPacket("udp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			t.Fatal(err)
		}
		serverStarted <- "started"
		defer listener.Close()
		for {
			buf := make([]byte, 1024)
			n, _, err := listener.ReadFrom(buf)
			if err != nil {
				t.Fatal(err)
			}
			receivedMessage <- string(buf[:n])
		}
	}()
	<-serverStarted

	// create graphite Client and sent two metrics
	graphiteClient := NewClient("127.0.0.1", port, "prefix", "udp")
	metricsMap := map[string]float64{
		"test":  1234.1234,
		"test2": 12345.12345,
	}
	err := graphiteClient.SendData(metricsMap)

	for i := 1; i <= 2; i++ {
		msg := <-receivedMessage
		found := false
		for _, n := range expectedMessages {
			if n == msg {
				found = true
			}
		}
		if !found {
			t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, expectedMessages)
		}
	}
	if err != nil {
		t.Fatalf("Unexpected error sending metrics:\nGot:\t\t%s\nExpected:\tnil\n", err)
	}
}

func TestSendFailure(t *testing.T) {
	// create bad graphite Client, expects error establishing connection
	graphiteClient := NewClient("bad_host", -1, "", "")
	err := graphiteClient.SendData(map[string]float64{})

	if err == nil {
		t.Fatal("Unexpected error sending metrics:\nGot:\t\tnil\nExpected:\tnot-nil\n")
	}
}
