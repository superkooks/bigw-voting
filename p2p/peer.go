package p2p

import (
	"bigw-voting/util"
	"errors"
	"fmt"
	"net"
	"time"
)

// Peer represents a connection to a peer
type Peer struct {
	PeerAddress *net.UDPAddr
	Messages    chan []byte
	Established bool
	MaxRTT      time.Duration

	latestSeqNumber int
	unackedMessages []*Message
}

var peers []*Peer

// StartConnection initiates a connection to a new peer. It reqbigwres an intermediate as well as an IP
func StartConnection(intermediate string, newPeerIPString string) (*Peer, error) {
	util.Infoln("Beginning new connection")

	intermediateAddr, err := net.ResolveUDPAddr("udp4", intermediate)
	if err != nil {
		return nil, fmt.Errorf("could not resolve intermediate udp: %v", err)
	}
	util.Infoln("Connecting to intermediate", intermediateAddr.String())

	newPeerIP, err := net.ResolveIPAddr("ip4", newPeerIPString)
	if err != nil {
		return nil, fmt.Errorf("could not resolve new peer ip: %v", err)
	}

	newPeer := &Peer{
		Messages: make(chan []byte),
		MaxRTT:   time.Second,
	}
	peers = append(peers, newPeer)

	// Retrieve external UDP address from the intermediate
	for {
		port.WriteToUDP([]byte("Unconnected Peers"), intermediateAddr)

		if newPeerIP.IP.String() == intermediateAddr.IP.String() {
			newPeer.PeerAddress = intermediateAddr
			break
		}

		for _, recieved := range recievedPeers {
			n, err := net.ResolveUDPAddr("udp4", recieved)
			if err != nil {
				return nil, err
			}

			if n.IP.String() == newPeerIP.IP.String() {
				newPeer.PeerAddress = n
				break
			}
		}

		if newPeer.PeerAddress != nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	// Send connection attempts to the external UDP address of the new peer
	for i := 0; i < 5; i++ {
		if !newPeer.Established {
			util.Infof("Sending connection attempt #%v to %v\n", i, newPeer.PeerAddress.String())

			_, err := port.WriteToUDP([]byte("Connection Attempt"), newPeer.PeerAddress)
			if err != nil {
				return nil, err
			}

			time.Sleep(3 * time.Second)
		}

		if !newPeer.Established && i == 4 {
			return nil, errors.New("timeout on connection attempt")
		}
	}

	return newPeer, nil
}

// SendMessage sends a message to a peer
func (p *Peer) SendMessage(msg []byte) error {
	newMessage := p.NewMessage(msg, false, false)
	p.unackedMessages = append(p.unackedMessages, newMessage)
	time.AfterFunc(p.MaxRTT, p.retransmission)

	util.Infof("Sending message with seq. number: %v\n", p.latestSeqNumber)
	_, err := port.WriteToUDP(newMessage.Serialize(), p.PeerAddress)
	if err != nil {
		return err
	}

	return nil
}

// retransmission is responsible for retransmitting lost packets
func (p *Peer) retransmission() {
	for _, packet := range p.unackedMessages {
		// Should this packet have been acked by now?
		if time.Now().After(packet.sentAt.Add(p.MaxRTT)) {
			// Increase the peer's MaxRTT after losing packet
			p.MaxRTT = time.Duration(float64(p.MaxRTT.Milliseconds())*1.75) * time.Millisecond

			// Retransmit the packet and add to p.unacked with new sentAt
			packet.sentAt = time.Now()
			time.AfterFunc(p.MaxRTT, p.retransmission)

			newPacket := packet.Serialize()
			util.Warnf("Retransmitting seq. number: %v\n", p.latestSeqNumber)
			_, err := port.WriteToUDP(newPacket, p.PeerAddress)
			if err != nil {
				panic(err)
			}
		}
	}
}

// BroadcastMessage sends a message to all peers, telling them to pass it on as well
// Note: Broadcast messages are not acked.
// Note: Peer receiving broadcasts should be able to handle duplicates
func BroadcastMessage(msg []byte) error {
	for _, p := range peers {
		newMessage := p.NewMessage(msg, false, true)
		p.unackedMessages = append(p.unackedMessages, newMessage)

		_, err := port.WriteToUDP(newMessage.Serialize(), p.PeerAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

// GossipPeers broadcasts our list of peers
func GossipPeers(intermediate string) {
	peersList := intermediate + " "
	for _, p := range peers {
		peersList += p.PeerAddress.String() + " "
	}

	BroadcastMessage([]byte(peersList))
}
