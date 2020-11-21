package main

import (
	"bigw-voting/bgw"
	"bigw-voting/shamir"
	"fmt"
)

func main() {
	// parseCommandline()
	// commands.RegisterAll()

	// go ui.Start()

	// time.Sleep(100 * time.Millisecond)
	// ui.NewVote([]string{"Lenin", "Stalin", "Krushchev", "Brezhnev"}, ui.SubmitVotes)

	// // Find local IP for BGW as well as for UPNP mapping
	// ifaces, err := net.Interfaces()
	// if err != nil {
	// 	panic(err)
	// }

	// var localIP string
	// for _, i := range ifaces {
	// 	addrs, err := i.Addrs()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	for _, addr := range addrs {
	// 		switch v := addr.(type) {
	// 		case *net.IPNet:
	// 			localIP = v.IP.String()
	// 		case *net.IPAddr:
	// 			localIP = v.IP.String()
	// 		}

	// 		break
	// 	}
	// }

	// externalIP := localIP

	// if !flagNoUPNP {
	// 	clients, _, err := upnp.NewWANIPConnection1Clients()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	if len(clients) > 1 {
	// 		panic("detected multiple gateway devices")
	// 	}

	// 	if len(clients) < 1 {
	// 		util.Warnln("Did not detect any gateway devices, if you are behind a NAT, you cannot act as an intermediate")
	// 	}

	// 	if len(clients) == 1 {
	// 		client := clients[0]

	// 		util.Infof("Using local IP %v for port mapping\n", localIP)

	// 		// Check for an entry before creating one
	// 		intPort, _, _, _, _, err := client.GetSpecificPortMappingEntry("", 42069, "udp")
	// 		if intPort != 42069 {
	// 			util.Infoln("Creating new port mapping")

	// 			// Create a new port mapping allowing all remotes to connect to us on port 42069
	// 			err = client.AddPortMapping("", 42069, "udp", 42069, localIP, true, "BIGW Voting", 900)
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 		}

	// 		util.Infoln("Port mapping is established")

	// 		// Get external IP
	// 		externalIP, err = client.GetExternalIPAddress()
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		util.Infof("Starting intermediate server at external IP: %v:42069\n", externalIP)
	// 	}
	// }

	// p2p.Setup()

	// newPeer, err := p2p.StartConnection(fmt.Sprintf("%v:%v", flagIntermediateIP, flagIntermediatePort), flagPeerIP)
	// if err != nil {
	// 	ui.Stop()
	// 	panic(err)
	// }

	// newPeer.SendMessage([]byte("Hello world!"))

	// for {
	// 	for _, p := range p2p.GetAllPeers() {
	// 		select {
	// 		case m := <-p.Messages:
	// 			util.Infof("New Packet: %v\n", string(m))
	// 		default:

	// 		}
	// 	}
	// }

	head1, p1Shares := bgw.NewVotingCircuit(2, "1.2.3.4", []string{"2.2.3.4", "3.2.3.4", "4.2.3.4"})
	head2, p2Shares := bgw.NewVotingCircuit(6, "2.2.3.4", []string{"1.2.3.4", "3.2.3.4", "4.2.3.4"})
	head3, p3Shares := bgw.NewVotingCircuit(3, "3.2.3.4", []string{"1.2.3.4", "2.2.3.4", "4.2.3.4"})
	head4, p4Shares := bgw.NewVotingCircuit(9, "4.2.3.4", []string{"1.2.3.4", "2.2.3.4", "3.2.3.4"})

	bgw.DescendCircuit(head1, append(p2Shares["1.2.3.4"], append(p3Shares["1.2.3.4"], p4Shares["1.2.3.4"]...)...))
	bgw.DescendCircuit(head2, append(p1Shares["2.2.3.4"], append(p3Shares["2.2.3.4"], p4Shares["2.2.3.4"]...)...))
	bgw.DescendCircuit(head3, append(p1Shares["3.2.3.4"], append(p2Shares["3.2.3.4"], p4Shares["3.2.3.4"]...)...))
	bgw.DescendCircuit(head4, append(p1Shares["4.2.3.4"], append(p2Shares["4.2.3.4"], p3Shares["4.2.3.4"]...)...))

	fmt.Println(head1.GetOutput())
	fmt.Println(head2.GetOutput())
	fmt.Println(head3.GetOutput())
	fmt.Println(head4.GetOutput())

	s, err := shamir.ReconstructSecret([][2]int{{1, head1.GetOutput()}, {2, head2.GetOutput()}, {3, head3.GetOutput()}, {4, head4.GetOutput()}})
	if err != nil {
		panic(err)
	}

	fmt.Println(s)
}
