package graphite

import (
	"io/ioutil"
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
	exceptedMessage := "prefix.test 1234.123400 1548015620\nprefix.test2 12345.123450 1548015620\n"
	receivedMessage := make(chan string, 1)
	serverStarted := make(chan string, 1)

	go func() {
		// start tcp server
		listener, err := net.Listen("tcp", "127.0.0.1:2003")
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
	graphiteClient := NewClient("127.0.0.1", 2003, "prefix", "tcp")
	metricsMap := map[string]float64{
		"test":  1234.1234,
		"test2": 12345.12345,
	}
	graphiteClient.SendData(metricsMap)

	if msg := <-receivedMessage; msg != exceptedMessage {
		t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, exceptedMessage)
	}
}

func TestSentMetricsOverUDP(t *testing.T) {
	exceptedMessage1 := "prefix.test 1234.123400 1548015620\n"
	exceptedMessage2 := "prefix.test2 12345.123450 1548015620\n"
	receivedMessage := make(chan string, 2)
	serverStarted := make(chan string, 1)

	go func() {
		// start UDP server
		listener, err := net.ListenPacket("udp", "127.0.0.1:2003")
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
	graphiteClient := NewClient("127.0.0.1", 2003, "prefix", "udp")
	metricsMap := map[string]float64{
		"test":  1234.1234,
		"test2": 12345.12345,
	}
	graphiteClient.SendData(metricsMap)

	if msg := <-receivedMessage; msg != exceptedMessage1 {
		t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, exceptedMessage1)
	}
	if msg := <-receivedMessage; msg != exceptedMessage2 {
		t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, exceptedMessage2)
	}
}
