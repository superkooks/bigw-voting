package p2p

import (
	"bigw-voting/ui"
	"net"
)

var port *net.UDPConn

// Setup starts the connection helper server and prepares
// for peer to peer interactions
func Setup() {
	localAddr, err := net.ResolveUDPAddr("udp4", ":42069")
	if err != nil {
		panic(err)
	}

	port, err = net.ListenUDP("udp", localAddr)
	if err != nil {
		panic(err)
	}

	ui.Infoln("Starting peer-to-peer listener")
	go listener()
}
