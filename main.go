package main

import (
	"bigw-voting/commands"
	"bigw-voting/p2p"
	"bigw-voting/ui"
	"bigw-voting/util"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	upnp "github.com/huin/goupnp/dcps/internetgateway2"
)

var votepack *Votepack

func main() {
	parseCommandline()
	commands.RegisterAll()

	go ui.Start()
	defer ui.Stop()

	time.Sleep(100 * time.Millisecond)
	votepack = NewVotepackFromFile(flagVotepackFilename)
	ui.NewVote(votepack.Candidates, ui.SubmitVotes)

	// Find local IP for BGW as well as for UPNP mapping
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var localIP string
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if !util.IsPublicIP(v.IP.String()) && v.IP.To4() != nil {
					localIP = v.IP.String()
					break
				}
			case *net.IPAddr:
				if !util.IsPublicIP(v.IP.String()) && v.IP.To4() != nil {
					localIP = v.IP.String()
					break
				}
			}
		}
	}

	var externalIP string
	if !flagNoUPNP {
		clients, _, err := upnp.NewWANIPConnection1Clients()
		if err != nil {
			panic(err)
		}

		if len(clients) > 1 {
			ui.Stop()
			panic("detected multiple gateway devices")
		}

		if len(clients) < 1 {
			util.Warnln("Did not detect any gateway devices, if you are behind a NAT, you cannot act as an intermediate")
		}

		if len(clients) == 1 {
			client := clients[0]

			util.Infof("Using local IP %v for port mapping\n", localIP)

			// Check for an entry before creating one
			intPort, _, _, _, _, err := client.GetSpecificPortMappingEntry("", 42069, "udp")
			if intPort != 42069 {
				util.Infoln("Creating new port mapping")

				// Create a new port mapping allowing all remotes to connect to us on port 42069 for 30 minutes
				err = client.AddPortMapping("", 42069, "udp", 42069, localIP, true, "BIGW Voting", 1800)
				if err != nil {
					panic(err)
				}
			}

			util.Infoln("Port mapping is established")

			// Get external IP
			externalIP, err = client.GetExternalIPAddress()
			if err != nil {
				panic(err)
			}
			util.Infof("Starting intermediate server at external IP: %v:42069\n", externalIP)
		}
	}

	if !util.IsPublicIP(externalIP) {
		var extIP string
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				panic(err)
			}

			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					if util.IsPublicIP(v.IP.String()) {
						extIP = v.IP.String()
						break
					}

				case *net.IPAddr:
					if util.IsPublicIP(v.IP.String()) {
						extIP = v.IP.String()
						break
					}
				}
			}
		}

		externalIP = extIP
	}

	p2p.Setup(externalIP, NewPeerCallback)

	_, err = p2p.StartConnection(fmt.Sprintf("%v:%v", flagIntermediateIP, flagIntermediatePort), flagPeerIP)
	if err != nil {
		ui.Stop()
		panic(err)
	}

	// Wait for Ctrl-C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

// NewPeerCallback serves as the callback for when new peers are connected.
// It verifies the votepack is consistent.
func NewPeerCallback(p *p2p.Peer) {
	// Start a goroutine (one per peer) for a listener
	go listener(p)

	util.Infoln("Verifying votepack with new peer")
	hash := sha256.Sum256(votepack.Export())
	err := p.SendMessage(append([]byte("VotepackVerify "), hash[:]...))
	if err != nil {
		util.Errorf("Unable to send message to %v, %v\n", p.PeerAddress.IP.String(), err)
	}

	ui.AddPeerToList(p.PeerAddress.IP.String(), "Votepack Verified")
}
