package commands

import (
	"bigw-voting/p2p"
	"bigw-voting/util"
)

// CommandNick changes the nickname shown by other peers
func CommandNick(args []string) {
	err := p2p.BroadcastMessage([]byte("Nick "+args[0]), 0)
	if err != nil {
		util.Errorln(err)
		return
	}
}
