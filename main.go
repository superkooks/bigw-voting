package main

import (
	"bigw-voting/p2p"
	"bigw-voting/ui"
	"fmt"
	"net"
	"time"

	upnp "github.com/huin/goupnp/dcps/internetgateway2"
)

func main() {
	parseCommandline()

	go ui.Start()

	time.Sleep(100 * time.Millisecond)
	ui.NewVote([]string{"Lenin", "Stalin", "Krushchev", "Brezhnev"}, ui.SubmitVotes)

	if !flagNoUPNP {
		clients, _, err := upnp.NewWANIPConnection1Clients()
		if err != nil {
			panic(err)
		}

		if len(clients) > 1 {
			panic("detected multiple gateway devices")
		}

		if len(clients) < 1 {
			ui.Warnln("Did not detect any gateway devices, if you are behind a NAT, you cannot act as an intermediate")
		}

		if len(clients) == 1 {
			client := clients[0]

			// Find local IP to map to
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
						localIP = v.IP.String()
					case *net.IPAddr:
						localIP = v.IP.String()
					}

					break
				}
			}

			ui.Infof("Using local IP %v for port mapping\n", localIP)

			// Check for an entry before creating one
			intPort, _, _, _, _, err := client.GetSpecificPortMappingEntry("", 42069, "udp")
			if intPort != 42069 {
				ui.Infoln("Creating new port mapping")

				// Create a new port mapping allowing all remotes to connect to us on port 42069
				err = client.AddPortMapping("", 42069, "udp", 42069, localIP, true, "BIGW Voting", 900)
				if err != nil {
					panic(err)
				}
			}

			ui.Infoln("Port mapping is established")

			// Get external IP
			remoteIP, err := client.GetExternalIPAddress()
			if err != nil {
				panic(err)
			}
			ui.Infof("Starting intermediate server at external IP: %v:42069\n", remoteIP)
		}
	}

	p2p.Setup()

	newPeer, err := p2p.StartConnection(fmt.Sprintf("%v:%v", flagIntermediateIP, flagIntermediatePort), flagPeerIP)
	if err != nil {
		panic(err)
	}

	newPeer.SendMessage([]byte("Hello world!"))
	ui.Infoln(string(<-newPeer.Messages))

	finished := make(chan bool)
	<-finished
}
