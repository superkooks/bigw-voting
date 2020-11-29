package p2p

import (
	"bigw-voting/util"
	"net"
)

var port *net.UDPConn
var externalIP string
var newPeerCallback func(*Peer)

// Setup starts the connection helper server and prepares
// for peer to peer interactions
func Setup(ip string, callback func(*Peer)) {
	// Store reference to external IP for gossiping peers
	externalIP = ip

	// Set the callback to be called when a new peer is connected
	newPeerCallback = callback

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

// SetNewPeerCallback sets the callback to be called when a
// new peer is connected
func SetNewPeerCallback(callback func(*Peer)) {
	newPeerCallback = callback
}
