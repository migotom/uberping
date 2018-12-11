package netcat

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type Netcat struct {
	ipaddr *net.IPAddr
	ip     string
	port   string

	// Timeout specifies tiemout for connection establishing.
	Timeout time.Duration

	connectionTries        int
	connectionsEstablished int
	rtt                    time.Duration
	err                    error

	// OnFinish is called when Netcat exits
	OnFinish func(*Statistics)
}

// Statistics represent the stats of a Netcat
type Statistics struct {
	// ConnectionTries specifies number of connection tries
	ConnectionTries int

	// ConnectionsEstablished secifies number of established connections
	ConnectionsEstablished int

	// ConnectionError specifies connection error.
	ConnectionError error

	// ConnectionLoss is the percentage of connections lost.
	ConnectionLoss float64

	// IPAddr is the address of the host being probed.
	IPAddr *net.IPAddr

	// Addr is the string address of the host being probed.
	Addr string

	// Port is service port number used to establish connection.
	Port string

	// Rtt is round-trip time of dialing to host.
	Rtt time.Duration
}

func NewNetcat(host, port string) (*Netcat, error) {
	ip, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return nil, err
	}

	if p, err := strconv.Atoi(port); err != nil || p <= 0 || p > 65535 {
		return nil, fmt.Errorf("wrong port number: %v", p)
	}

	return &Netcat{
		ipaddr: ip,
		ip:     ip.String(),
		port:   port,
	}, nil
}

func (n *Netcat) Run() {
	n.connectionTries++

	// TODO refactor this to nonblocking version with tries count > 1
	connT := time.Now()
	connection, err := net.DialTimeout("tcp", n.ip+":"+n.port, n.Timeout)

	if err == nil {
		n.connectionsEstablished++
		n.rtt = time.Since(connT)
		connection.Close()
	} else {
		n.err = err
	}

	handler := n.OnFinish
	if handler != nil {
		s := n.Statistics()
		handler(s)
	}
}

func (n *Netcat) Statistics() *Statistics {
	loss := float64(n.connectionTries-n.connectionsEstablished) / float64(n.connectionTries) * 100
	s := Statistics{
		ConnectionTries:        n.connectionTries,
		ConnectionsEstablished: n.connectionsEstablished,
		IPAddr:                 n.ipaddr,
		Addr:                   n.ip,
		Port:                   n.port,
		Rtt:                    n.rtt,
		ConnectionError:        n.err,
		ConnectionLoss:         loss,
	}
	return &s
}
