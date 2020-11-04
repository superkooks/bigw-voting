package main

import (
	flag "github.com/spf13/pflag"
)

var (
	flagIntermediateIP   string
	flagIntermediatePort int
	flagPeerIP           string

	flagNoUPNP                bool
	flagConfirmNewConnections bool
)

func parseCommandline() {
	flag.StringVarP(&flagIntermediateIP, "intermediateIP", "i", "", "The IPv4 address of the intermediate to connect through")
	flag.IntVarP(&flagIntermediatePort, "intermediatePort", "o", 42069, "The UDP port of the intermediate to connect through")
	flag.StringVarP(&flagPeerIP, "peerIP", "p", "", "The IPv4 address of the peer to begin connecting with")

	flag.BoolVar(&flagNoUPNP, "noUPNP", false, "Should the local intermediate server use UPNP port forwarding")
	flag.BoolVar(&flagConfirmNewConnections, "confirmNewConnections", false, "Require confirmation before connecting to gossipped peers")

	flag.Parse()

	if flagIntermediateIP == "" {
		panic("intermediate IP address is a required flag")
	}

	if flagPeerIP == "" {
		panic("peer IP address is a required flag")
	}
}
