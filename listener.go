package main

import (
	"bigw-voting/p2p"
	"bigw-voting/util"
	"crypto/sha256"
	"strings"
)

// listener is the function primarily responsible for listening to and
// responding to messages.
func listener(p *p2p.Peer) {
	for {
		msg := <-p.Messages

		msgSplit := strings.Split(string(msg), " ")
		if msgSplit[0] == "VotepackVerify" {
			localHash := sha256.Sum256(votepack.Export())

			// Only respond to peer if there is a difference in votepack
			if msgSplit[1] != string(localHash[:]) {
				err := p.SendMessage([]byte("VotepackInvalid"))
				if err != nil {
					util.Errorf("Unable to send message to %v, %v\n", p.PeerAddress.IP.String(), err)
				}
			}

			continue
		}

		if msgSplit[0] == "VotepackInvalid" {
			util.Errorf("Peer %v is using a different votepack\n", p.PeerAddress.IP.String())
			continue
		}
	}
}
