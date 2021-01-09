package p2p

import (
	"bigw-voting/util"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// Peer represents a connection to a peer
type Peer struct {
	PeerAddress *net.UDPAddr
	Messages    chan []byte
	Established bool
	MaxRTT      time.Duration

	peerSeqNumber   int
	localSeqNumber  int
	unackedMessages []*Message
}

var peers []*Peer
var connectingTo []string

// StartConnection initiates a connection to a new peer. It requires an intermediate as well as an IP
func StartConnection(intermediate string, newPeerIPString string) (*Peer, error) {
	for _, v := range connectingTo {
		if v == newPeerIPString {
			return nil, fmt.Errorf("Already connecting to %v", newPeerIPString)
		}
	}

	util.Infoln("Beginning new connection to", newPeerIPString)
	connectingTo = append(connectingTo, newPeerIPString)

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

		for _, received := range receivedPeers {
			n, err := net.ResolveUDPAddr("udp4", received)
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

			time.Sleep(1 * time.Second)
		}

		if !newPeer.Established && i == 4 {
			return nil, errors.New("timeout on connection attempt")
		}
	}

	GossipPeers(intermediate)

	// Wait for all gossips to complete to ensure no overlapping sequence numbers
	// time.Sleep(2 * time.Second)

	newPeerCallback(newPeer)

	return newPeer, nil
}

// SendMessage sends a message to a peer
func (p *Peer) SendMessage(msg []byte) error {
	newMessage := p.NewMessage(msg, false, false)
	p.unackedMessages = append(p.unackedMessages, newMessage)
	time.AfterFunc(p.MaxRTT, p.retransmission)

	util.Infof("Sending message with seq. number: %v\n", p.localSeqNumber)
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

			// Peer should be disconnected
			if p.MaxRTT > 8*time.Second {
				for k, v := range peers {
					if v.PeerAddress.IP.String() == p.PeerAddress.String() {
						peers = append(peers[:k], peers[k+1:]...)
						break
					}
				}
			}

			// Retransmit the packet and add to p.unacked with new sentAt
			packet.sentAt = time.Now()
			time.AfterFunc(p.MaxRTT, p.retransmission)

			newPacket := packet.Serialize()
			util.Warnf("Retransmitting seq. number: %v\n", p.localSeqNumber)
			_, err := port.WriteToUDP(newPacket, p.PeerAddress)
			if err != nil {
				panic(err)
			}
		}
	}
}

// BroadcastMessage sends a message to all peers, telling them to pass it on as well
// Note: Acks for broadcasts are only received from peers currently connected to this peer
// Note: Peer receiving broadcasts should be able to handle duplicates
func BroadcastMessage(msg []byte, maxBounces int8) error {
	// Do not broadcast message which has reached maxBounces
	if maxBounces < 0 {
		return nil
	}

	for _, p := range peers {
		if !p.Established {
			continue
		}

		newMessage := p.NewMessage(msg, false, true)
		newMessage.MaxBounces = maxBounces

		util.Infof("Sending broadcast to %v with seq. number: %v\n", p.PeerAddress.IP.String(), newMessage.SequenceNumber)
		_, err := port.WriteToUDP(newMessage.Serialize(), p.PeerAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

// GossipPeers broadcasts our list of peers
func GossipPeers(intermediate string) {
	// TODO(SuperKooks) These ports are hard-coded
	if intermediate == "127.0.0.1:42069" || intermediate == "localhost:42069" {
		intermediate = externalIP + ":42069"
	}

	peersList := "Gossip " + intermediate + " "
	for _, p := range peers {
		if p.Established {
			peersList += p.PeerAddress.IP.String() + " "
		}
	}

	// Remove unecessary space after last peer
	peersList = strings.TrimRight(peersList, " ")

	util.Warnln("Sending:", peersList)

	err := BroadcastMessage([]byte(peersList), 1)
	if err != nil {
		util.Errorln(err)
	}
}

// GetAllPeers returns the list of currently connected peers
func GetAllPeers() []*Peer {
	return peers
}

// GetAllPeerIPs returns a list of all IP addresses of the currently connected peers
func GetAllPeerIPs() []string {
	var out []string
	for _, v := range peers {
		if v.Established {
			out = append(out, v.PeerAddress.IP.String())
		}
	}

	return out
}
