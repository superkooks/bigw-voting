package ui

import (
	"bigw-voting/util"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

var app *cview.Application
var ir *instantRunoff
var consoleInput *cview.InputField
var consoleOutput *cview.TextView

// Start starts up the UI
func Start() {
	app = cview.NewApplication()
	app.EnableMouse(true)

	ir = newInstantRunoff()
	ir.SetVisible(false)

	subFlex := cview.NewFlex()
	subFlex.SetDirection(cview.FlexRow)
	subFlex.AddItem(ir, 0, 3, true)
	subFlex.AddItem(demoBox("Connected Peers"), 0, 2, false)

	rootFlex := cview.NewFlex()
	console := consoleBox()
	rootFlex.AddItem(console, 0, 3, false)
	rootFlex.AddItem(subFlex, 0, 2, true)

	app.SetRoot(rootFlex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// NewVote shows a new voting widget to allow voting in
func NewVote(choices []string, callback func(map[string]int)) {
	ir.SetupNewVote(choices, callback)
	ir.SetVisible(true)
}

// ClearVote hides the voting widget
func ClearVote() {
	ir.SetVisible(false)
}

// SubmitVotes prints the votes
func SubmitVotes(votes map[string]int) {
	for k, v := range votes {
		util.Infof("%v: %v", k, v)
	}
}

func demoBox(title string) *cview.Box {
	b := cview.NewBox()
	b.SetBorder(true)
	b.SetTitle(title)
	b.SetBorderColor(tcell.ColorLime)
	b.SetTitleColor(tcell.ColorLime)
	return b
}
