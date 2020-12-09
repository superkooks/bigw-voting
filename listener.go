package main

import (
	"bigw-voting/commands"
	"bigw-voting/ui"
	"bigw-voting/util"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var receivedPeerShares = make(map[string][]int)
var receivedPeerOutputs = make(map[string]int)

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
			v.Status = strings.Join(msgSplit[1:], " ")
			util.Infof("Peer %v is now %v\n", v.Peer.PeerAddress.IP.String(), v.Status)

			ui.SetStatusOfPeer(v.Peer.PeerAddress.IP.String(), v.Status)

			var dontBegin bool
			for _, v := range allVoters {
				if v.Status != "Voting Complete" {
					util.Infoln("Peers are still voting")

					// Don't begin BGW if peers have not finished voting
					dontBegin = true
					break
				}
			}

			if dontBegin || localStatus != "Voting Complete" {
				continue
			}

			// All peers have voted, proceed with BGW
			go RunBGW()
			continue
		}

		if msgSplit[0] == "YourShares" {
			var shares []int
			err := json.Unmarshal([]byte(strings.Join(msgSplit[1:], " ")), &shares)
			if err != nil {
				util.Errorln("Failed to unmarshal peer shares")
			}

			receivedPeerShares[v.Peer.PeerAddress.IP.String()] = shares
			continue
		}

		if msgSplit[0] == "MyOutput" {
			conv, err := strconv.Atoi(msgSplit[1])
			if err != nil {
				util.Errorln("Unable to parse output recieved from peer")
			}

			util.Errorf("Received output %v from %v\n", conv, v.Peer.PeerAddress.IP.String())
			receivedPeerOutputs[v.Peer.PeerAddress.IP.String()] = conv
			continue
		}

		if msgSplit[0] == "Nick" {
			nick := strings.Join(msgSplit[1:], " ")
			ui.SetNickOfPeer(v.Peer.PeerAddress.IP.String(), nick)
			continue
		}

		if msgSplit[0] == "TrusteeVote" {
			trustee := strings.Join(msgSplit[1:], " ")

			var found bool
			for _, v := range votepack.TrusteeVotes {
				if v == trustee {
					found = true
					break
				}
			}

			if !found {
				util.Errorf("Peer %v has been caught cheating. Attempted to add trustee vote not in votepack.\n", v.Peer.PeerAddress.IP.String())
				continue
			}

			found = false
			for _, v := range commands.AllTrusteeVotes {
				if v == trustee {
					found = true
					break
				}
			}

			if found {
				util.Errorln("Peer %v has been caught cheating. Attempted to add trustee vote for person who already has trustee.\n", v.Peer.PeerAddress.IP.String())
				continue
			}

			commands.AllTrusteeVotes = append(commands.AllTrusteeVotes, trustee)
			util.Infof("Peer %v is trustee for %v\n", v.Peer.PeerAddress.IP.String(), trustee)
			continue
		}
	}
}
