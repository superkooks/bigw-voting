package main

import "flag"

var (
	flagShouldUseUPNP      bool
	flagAutoConnectToPeers bool
)

func parseCommandline() {
	flag.BoolVar(&flagShouldUseUPNP, "useUPNP", true, "Should the intermediate use UPNP port forwarding")
	flag.BoolVar(&flagAutoConnectToPeers, "autoConnectToPeers", true, "Should we connect to gossipped peers")

	flag.Parse()
}
