package p2p

import (
	"bigw-voting/util"
	"net"
)

var port *net.UDPConn
var externalIP string

// Setup starts the connection helper server and prepares
// for peer to peer interactions
func Setup(ip string) {
	// Store reference to external IP for gossiping peers
	externalIP = ip

	localAddr, err := net.ResolveUDPAddr("udp4", ":42069")
	if err != nil {
		panic(err)
	}

	port, err = net.ListenUDP("udp", localAddr)
	if err != nil {
		panic(err)
	}

	util.Infoln("Starting peer-to-peer listener")
	go listener()
}
