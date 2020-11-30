package main

import (
	"bigw-voting/ui"
	"bigw-voting/util"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
)

var receivedPeerShares map[string]int
var receivedPeerOutputs map[string]int

// listener is the function primarily responsible for listening to and
// responding to messages.
func listener(v *Voter) {
	for {
		msg := <-v.Peer.Messages

		msgSplit := strings.Split(string(msg), " ")
		if msgSplit[0] == "VotepackVerify" {
			localHash := sha256.Sum256(votepack.Export())

			// Only respond to peer if there is a difference in votepack
			if msgSplit[1] != string(localHash[:]) {
				err := v.Peer.SendMessage([]byte("VotepackInvalid"))
				if err != nil {
					util.Errorf("Unable to send message to %v, %v\n", v.Peer.PeerAddress.IP.String(), err)
				}
			}

			continue
		}

		if msgSplit[0] == "VotepackInvalid" {
			ui.Stop()
			panic(fmt.Errorf("peer %v is using a different votepack", v.Peer.PeerAddress.IP.String()))
		}

		if msgSplit[0] == "StatusUpdate" {
			v.Status = fmt.Sprint(msgSplit[1:])
			util.Infof("Peer %v is now %v\n", v.Peer.PeerAddress.IP.String(), v.Status)

			var dontBegin bool
			for _, v := range allVoters {
				if v.Status != "Voting Complete" {
					// Don't begin BGW if peers have not finished voting
					dontBegin = true
					break
				}
			}

			if dontBegin {
				continue
			}

			// All peers have voted, proceed with BGW
			RunBGW()
		}

		if msgSplit[0] == "YourShare" {
			conv, err := strconv.Atoi(msgSplit[1])
			if err != nil {
				util.Errorln("Unable to parse share recieved from peer")
			}

			receivedPeerShares[v.Peer.PeerAddress.IP.String()] = conv
		}

		if msgSplit[0] == "MyOutput" {
			conv, err := strconv.Atoi(msgSplit[1])
			if err != nil {
				util.Errorln("Unable to parse output recieved from peer")
			}

			receivedPeerOutputs[v.Peer.PeerAddress.IP.String()] = conv
		}
	}
}
