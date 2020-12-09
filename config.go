package main

import (
	"bigw-voting/util"

	flag "github.com/spf13/pflag"
)

var (
	flagIntermediateIP   string
	flagIntermediatePort int
	flagPeerIP           string
	flagVotepackFilename string

	flagLogFile bool
	flagNoUPNP  bool

	flagTrusteeExport bool
)

func parseCommandline() {
	flag.StringVarP(&flagIntermediateIP, "intermediateIP", "i", "", "The IPv4 address of the intermediate to connect through")
	flag.IntVarP(&flagIntermediatePort, "intermediatePort", "o", 42069, "The UDP port of the intermediate to connect through, the default is 42069")
	flag.StringVarP(&flagPeerIP, "peerIP", "p", "", "The IPv4 address of the peer to begin connecting with")
	flag.StringVarP(&flagVotepackFilename, "votepack", "v", "", "The filename of the votepack to use")

	flag.BoolVar(&flagLogFile, "log", false, "Should messages be written to debug.log")
	flag.BoolVar(&flagNoUPNP, "noUPNP", false, "Should the local intermediate server use UPNP port forwarding")

	flag.BoolVar(&flagTrusteeExport, "exportTrusteeVote", false, "Export a new trustee vote")

	flag.Parse()

	// Check for config errors
	if flagVotepackFilename == "" {
		panic("a votepack must be specified")
	}

	if !flagTrusteeExport {
		if flagIntermediateIP == "" {
			panic("intermediate IP address is a required flag")
		}

		if flagPeerIP == "" {
			panic("peer IP address is a required flag")
		}
	}

	// Check whether we should start the debug log
	if flagLogFile {
		util.SetDualLogging()
	}
}
