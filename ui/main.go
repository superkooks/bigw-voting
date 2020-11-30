package ui

import (
	"os"

	"gitlab.com/tslocum/cview"
)

var app *cview.Application
var ir *instantRunoff
var pl *cview.List
var consoleInput *cview.InputField
var consoleOutput *cview.TextView

// Start starts up the UI
func Start() {
	app = cview.NewApplication()
	app.EnableMouse(true)

	ir = newInstantRunoff()
	ir.visible = false

	pl = newPeerList()

	subFlex := cview.NewFlex()
	subFlex.SetDirection(cview.FlexRow)
	subFlex.AddItem(ir, 0, 3, true)
	subFlex.AddItem(pl, 0, 2, false)

	rootFlex := cview.NewFlex()
	console := consoleBox()
	rootFlex.AddItem(console, 0, 3, false)
	rootFlex.AddItem(subFlex, 0, 2, true)

	app.SetRoot(rootFlex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}

	os.Exit(0)
}

// NewVote shows a new voting widget to allow voting in
func NewVote(choices []string, callback func(map[string]int)) {
	ir.SetupNewVote(choices, callback)
	ir.visible = true
}

// ClearVote hides the voting widget
func ClearVote() {
	ir.SetVisible(false)
}

// Stop closes the application
func Stop() {
	app.Stop()
}
