package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

var plItems []*cview.ListItem

func newPeerList() *cview.List {
	list := cview.NewList()
	list.SetBorder(true)
	list.SetTitle("Peers")
	list.SetBorderColor(tcell.ColorLime)
	list.SetTitleColor(tcell.ColorLime)
	list.SetMainTextColor(tcell.ColorRed)

	list.SetPadding(0, 0, 2, 0)
	list.SetSelectedFocusOnly(true)
	list.SetSelectedAlwaysVisible(false)

	return list
}

// AddPeerToList adds a peer to the list
func AddPeerToList(ip string, status string) {
	item := cview.NewListItem(ip)
	item.SetSecondaryText(status)
	pl.AddItem(item)

	plItems = append(plItems, item)
}

// SetStatusOfPeer sets the secondary text of the list item to status
func SetStatusOfPeer(ip string, status string) {
	for _, v := range plItems {
		if strings.Split(v.GetMainText(), " ")[0] == ip {
			v.SetSecondaryText(status)
		}
	}
}

// SetNickOfPeer sets the nickname of the peer
func SetNickOfPeer(ip string, nick string) {
	for _, v := range plItems {
		if strings.Split(v.GetMainText(), " ")[0] == ip {
			v.SetMainText(ip + " (" + nick + ")")
		}
	}
}
