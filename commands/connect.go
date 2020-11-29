package commands

import (
	"bigw-voting/p2p"
	"bigw-voting/util"
	"fmt"
	"log"

	flag "github.com/spf13/pflag"
)

// CommandConnect connects to a new peer
func CommandConnect(args []string) {
	flagSet := flag.NewFlagSet("connect", flag.ContinueOnError)
	intIP := flagSet.StringP("intermediateIP", "i", "", "The IPv4 of the intermediate to conect through")
	intPort := flagSet.IntP("intermediatePort", "o", 42069, "The UDP port of the intermediate to connect through, the default is 42069")
	peerIP := flagSet.StringP("peerIP", "p", "", "The IPv4 address of the peer to begin connecting with")

	flagSet.SetOutput(log.Writer())

	err := flagSet.Parse(args)
	if err == flag.ErrHelp {
		return
	} else if err != nil {
		util.Errorf("could not parse command: %v\n", err)
		return
	}

	p, err := p2p.StartConnection(fmt.Sprintf("%v:%v", *intIP, *intPort), *peerIP)
	if err != nil {
		util.Errorln(err)
		return
	}

	err = p.SendMessage([]byte("What is up?"))
	if err != nil {
		util.Errorln(err)
		return
	}
}
