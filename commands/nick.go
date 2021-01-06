package commands

import (
	"bigw-voting/p2p"
	"bigw-voting/util"
	"strings"
)

// CommandNick changes the nickname shown by other peers
func CommandNick(args []string) {
	err := p2p.BroadcastMessage([]byte("Nick "+strings.Join(args, " ")), 0)
	if err != nil {
		util.Errorln(err)
		return
	}
}
