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

	// create graphite client and sent metrics in separate gorutine
	go func() {
		graphiteClient := NewClient("localhost", 2003, "prefix", "tcp")
		listMetrics := make([]map[string]float64, 0)

		listMetrics = append(listMetrics, map[string]float64{"test": 1234.1234})
		listMetrics = append(listMetrics, map[string]float64{"test2": 12345.12345})
		graphiteClient.SendData(listMetrics)
	}()

	// start tcp server
	listener, err := net.Listen("tcp", ":2003")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			return
		}
		defer conn.Close()

		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			t.Fatal(err)
		}

		//fmt.Println(string(buf[:]))
		if msg := string(buf[:]); msg != exceptedMessage {
			t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, exceptedMessage)
		}
		return
	}
}

func TestSentMetricsOverUDP(t *testing.T) {
	messageExcepted := "prefix.test 1234.123400 1548015620\nprefix.test2 12345.123450 1548015620\n"
	messageReceived := ""

	// start UDP server
	listener, err := net.ListenPacket("udp", ":2003")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			buf := make([]byte, 1024)
			n, _, err := listener.ReadFrom(buf)
			if err != nil {
				t.Fatal(err)
			}

			//fmt.Println(string(buf[:n]))
			messageReceived = messageReceived + string(buf[:n])
		}
	}()

	// create graphite Client and sent two metrics
	graphiteClient := NewClient("localhost", 2003, "prefix", "udp")
	listMetrics := make([]map[string]float64, 0)
	listMetrics = append(listMetrics, map[string]float64{"test": 1234.1234})
	listMetrics = append(listMetrics, map[string]float64{"test2": 12345.12345})
	graphiteClient.SendData(listMetrics)

	time.Sleep(time.Millisecond * 100)
	if messageExcepted != messageReceived {
		t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", messageReceived, messageExcepted)
	}
}
