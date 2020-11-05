package p2p

import (
	"bigw-voting/ui"
	"net"
	"strings"
	"time"
)

var unconnectedPeers []*net.UDPAddr
var recievedPeers []string

func listener() {
	for {
		buf := make([]byte, 2048)
		n, replyTo, err := port.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}

		// Is the person a peer?
		p := getPeer(replyTo)
		if p == nil {
			var shouldNotAdd bool
			for _, unconnectedP := range unconnectedPeers {
				if unconnectedP.String() == replyTo.String() {
					shouldNotAdd = true
					break
				}
			}

			if !shouldNotAdd {
				ui.Infof("Adding %v to unconnected list\n", replyTo.String())
				unconnectedPeers = append(unconnectedPeers, replyTo)
			}
		}

		// "Unconnected Peers" returns a list of unconnected peers
		if string(buf[:n]) == "Unconnected Peers" {
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

				ui.Infof("Recieved a connection attempt from %v\n", replyTo.String())
			} else {
				ui.Warnln("Recieved connection attempt from unknown peer, not responding")
			}
			continue
		}

		// Check to see whether it our connection has been established
		if string(buf[:n]) == "Established" {
			p.Established = true
			ui.Infof("Established new connection with %v\n", p.PeerAddress.String())
			continue
		}

		// Check if it is a list of peers being recieved from an intermediate
		split := strings.Split(string(buf[:n]), " ")
		if split[0] == "Unconnected" && split[1] == "Peers" {
			recievedPeers = split[2 : len(split)-1]
			continue
		}

		// Start parsing normal packets
		if p == nil {
			ui.Warnln("Recieved packet from unknown peer")
			continue
		}

		msg := new(Message)
		msg.Deserialize(buf[:n])
		if msg.Ack {
			ui.Infof("Received ack for seq. number: %v\n", msg.SequenceNumber)
			for k, unacked := range p.unackedMessages {
				if unacked.SequenceNumber == msg.SequenceNumber {
					p.MaxRTT = 2 * time.Now().Sub(unacked.sentAt)
					p.unackedMessages = append(p.unackedMessages[:k], p.unackedMessages[k+1:]...)
				}
			}

			continue
		}

		ui.Infof("Acking seq. number: %v\n", msg.SequenceNumber)
		_, err = port.WriteToUDP((&Message{Data: []byte{}, SequenceNumber: msg.SequenceNumber, Ack: true}).Serialize(), replyTo)
		if err != nil {
			panic(err)
		}

		p.latestSeqNumber = msg.SequenceNumber

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
