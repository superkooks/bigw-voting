package p2p

import (
	"bigw-voting/util"
	"net"
	"strings"
	"time"
)

var unconnectedPeers []*net.UDPAddr
var receivedPeers []string

func listener() {
	for {
		buf := make([]byte, 2048)
		n, replyTo, err := port.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}

		// Is the person a peer?
		p := getPeer(replyTo)

		// "Unconnected Peers" returns a list of unconnected peers
		if string(buf[:n]) == "Unconnected Peers" {
			var shouldNotAdd bool
			for _, unconnectedP := range unconnectedPeers {
				if unconnectedP.String() == replyTo.String() {
					shouldNotAdd = true
					break
				}
			}

			if !shouldNotAdd {
				util.Infof("Adding %v to unconnected list\n", replyTo.String())
				unconnectedPeers = append(unconnectedPeers, replyTo)
			}

			out := "Unconnected Peers "
			for _, p := range unconnectedPeers {
				out += p.String() + " "
			}

			_, err := port.WriteToUDP([]byte(out), replyTo)
			if err != nil {
				panic(err)
			}
			continue
		}

		// Check for connection setup messages
		if string(buf[:n]) == "Connection Attempt" {
			// Don't respond to peers we haven't started a connection with
			if p != nil {
				_, err = port.WriteToUDP([]byte("Established"), replyTo)
				if err != nil {
					panic(err)
				}

				util.Infof("Received a connection attempt from %v\n", replyTo.String())
			} else {
				util.Warnln("Received connection attempt from unknown peer, not responding")
			}
			continue
		}

		// Check to see whether it our connection has been established
		if string(buf[:n]) == "Established" {
			p.Established = true
			util.Infof("Established new connection with %v\n", p.PeerAddress.String())

			for k, v := range unconnectedPeers {
				if p.PeerAddress.String() == v.String() {
					unconnectedPeers = append(unconnectedPeers[:k], unconnectedPeers[k+1:]...)
				}
			}

			continue
		}

		// Check if it is a list of peers being received from an intermediate
		split := strings.Split(string(buf[:n]), " ")
		if len(split) > 1 {
			if split[0] == "Unconnected" && split[1] == "Peers" {
				receivedPeers = split[2 : len(split)-1]
				continue
			}
		}

		// Start parsing normal packets
		if p == nil {
			util.Warnln("Received packet from unknown peer")
			continue
		}

		msg := new(Message)
		msg.Deserialize(buf[:n])
		if msg.Ack {
			util.Infof("Received ack for seq. number: %v\n", msg.SequenceNumber)
			for k, unacked := range p.unackedMessages {
				if unacked.SequenceNumber == msg.SequenceNumber {
					p.MaxRTT = 2 * time.Now().Sub(unacked.sentAt)
					p.unackedMessages = append(p.unackedMessages[:k], p.unackedMessages[k+1:]...)
				}
			}

			continue
		}

		// Discard duplicate packets
		if msg.SequenceNumber <= p.peerSeqNumber {
			util.Warnf("Discarding duplicate packet %v\n", msg.SequenceNumber)
			util.Errorf("AHHHHHHHHHHHHH: %v\n", string(msg.Data))
			continue
		}

		// Don't ack broadcast packets
		if !msg.Broadcast {
			util.Infof("Acking seq. number: %v\n", msg.SequenceNumber)
			_, err = port.WriteToUDP((&Message{Data: []byte{}, SequenceNumber: msg.SequenceNumber, Ack: true}).Serialize(), replyTo)
			if err != nil {
				panic(err)
			}

			p.peerSeqNumber = msg.SequenceNumber
		} else {
			util.Infoln("Received broadcast packet, passing on to other peers")
			util.Warnf("broadcast content: %v\n", string(msg.Data))
			p.peerSeqNumber = msg.SequenceNumber
			err := BroadcastMessage(msg.Data, msg.MaxBounces-1)
			if err != nil {
				util.Errorln(err)
				continue
			}

			util.Warnln("finished broadcast")
		}

		msgSplit := strings.Split(string(msg.Data), " ")

		// Check if this is a Peer Gossip
		if msgSplit[0] == "Gossip" {
			util.Errorf("gossip message: %v\n", msgSplit)

			intermediate := msgSplit[1]

			util.Errorf("gossiped peers? %v\n", msgSplit[2:])
			for _, gossiped := range msgSplit[2:] {
				util.Errorf("now checking %v\n", gossiped)

				if gossiped == " " || gossiped == "" {
					continue
				}

				var found bool
				for _, connected := range GetAllPeerIPs() {
					if connected == gossiped || externalIP == gossiped {
						found = true
						break
					}
				}

				// We are not connected to this peer, we should connect
				if !found {
					var shouldSkip bool
					for _, v := range connectingTo {
						if v == gossiped {
							util.Infof("Gossipped peer %v is already being connected to\n", gossiped)
							shouldSkip = true
							break
						}
					}

					if shouldSkip {
						continue
					}

					util.Errorf("connecting to %v\n", gossiped)

					go func(gossiped string, intermediate string, externalIP string) {
						_, err := StartConnection(intermediate, gossiped)
						if err != nil {
							util.Errorln(err)
						}
					}(gossiped, intermediate, externalIP)
				}
			}

			// Don't pass on peer gossip to the message listener outside this module
			continue
		}

		// Pass on the message
		p.Messages <- msg.Data
	}
}

func getPeer(addr *net.UDPAddr) *Peer {
	for _, p := range peers {
		if addr.String() == p.PeerAddress.String() {
			return p
		}
	}

	return nil
}
