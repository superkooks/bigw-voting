package ui

import (
	"bigw-voting/commands"
	"fmt"
	"io"
	"log"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

func consoleBox() *cview.Flex {
	flex := cview.NewFlex()
	flex.SetBorder(true)
	flex.SetDirection(cview.FlexRow)
	flex.SetBorderColor(tcell.ColorLime)
	flex.SetTitleColor(tcell.ColorLime)
	flex.SetTitle("Console")

	t := cview.NewTextView()
	t.SetTextColor(tcell.ColorLime)
	t.SetDynamicColors(true)
	t.SetChangedFunc(func() { app.QueueUpdateDraw(func() {}) })

	consoleOutput = t

	flex.AddItem(t, 0, 20, false)

	i := cview.NewInputField()
	i.SetPadding(0, 0, 2, 0)
	i.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		cview.Print(screen, []byte("> "), x, y, width, cview.AlignLeft, tcell.ColorLime)
		return x + 2, y, width - 2, height
	})

	i.SetFieldBackgroundColor(tcell.ColorBlack.TrueColor())
	i.SetFieldBackgroundColorFocused(tcell.ColorBlack.TrueColor())
	i.SetFieldTextColor(tcell.ColorLime)
	i.SetFieldTextColorFocused(tcell.ColorLime)
	i.SetDoneFunc(runCommand)

	consoleInput = i

	flex.AddItem(i, 1, 3, true)

	log.SetOutput(t)
	log.SetFlags(log.Ltime)

	return flex
}

func runCommand(event tcell.Key) {
	if event == tcell.KeyEnter {
		t := consoleInput.GetText()
		consoleInput.SetText("")
		GetConsoleWriter().Write([]byte("\n"))
		GetConsoleWriter().Write([]byte(fmt.Sprintf("> %v\n", t)))

		commands.Parse(t)
	}
}

// GetConsoleWriter returns the io.Writer to write to the console
func GetConsoleWriter() io.Writer {
	return consoleOutput
}
