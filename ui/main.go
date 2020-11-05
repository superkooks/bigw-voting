package ui

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

var app *cview.Application
var ir *instantRunoff

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
		Infof("%v: %v", k, v)
	}
}

func consoleBox() *cview.TextView {
	t := cview.NewTextView()
	t.SetBorder(true)
	t.SetTitle("Console")
	t.SetBorderColor(tcell.ColorLime)
	t.SetTitleColor(tcell.ColorLime)
	t.SetTextColor(tcell.ColorLime)
	t.SetDynamicColors(true)
	t.SetChangedFunc(func() { app.QueueUpdateDraw(func() {}) })

	log.SetOutput(t)
	log.SetFlags(log.Ltime)

	return t
}

// Errorf prints a formatted error to the ui console
func Errorf(fmtString string, fmtArgs ...interface{}) {
	log.Printf("[red]error:[lime] %v", fmt.Sprintf(fmtString, fmtArgs...))
}

// Errorln prints text to the ui console
func Errorln(msg ...interface{}) {
	v := append([]interface{}{"[red]error:[lime]"}, msg...)
	log.Println(v...)
}

// Warnf prints a formatted warning to the ui console
func Warnf(fmtString string, fmtArgs ...interface{}) {
	log.Printf("[orange]warn:[lime] %v", fmt.Sprintf(fmtString, fmtArgs...))
}

// Warnln prints text to the ui console
func Warnln(msg ...interface{}) {
	v := append([]interface{}{"[orange]warn:[lime]"}, msg...)
	log.Println(v...)
}

// Infof prints a formatted Infoing to the ui console
func Infof(fmtString string, fmtArgs ...interface{}) {
	log.Printf("[lime] %v", fmt.Sprintf(fmtString, fmtArgs...))
}

// Infoln prints text to the ui console
func Infoln(msg ...interface{}) {
	v := append([]interface{}{"[lime]"}, msg...)
	log.Println(v...)
}

func demoBox(title string) *cview.Box {
	b := cview.NewBox()
	b.SetBorder(true)
	b.SetTitle(title)
	b.SetBorderColor(tcell.ColorLime)
	b.SetTitleColor(tcell.ColorLime)
	return b
}
