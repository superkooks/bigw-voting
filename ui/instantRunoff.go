package ui

import (
	"fmt"
	"strconv"
	"sync"
	"unicode"

	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

type instantRunoff struct {
	*cview.Box

	allCandidates []string
	cursor        int
	votes         map[int]int
	submitFunc    func(map[string]int)
	currentStatus string

	l sync.RWMutex
}

// newInstantRunoff creates a new instantRunoff widget
func newInstantRunoff() *instantRunoff {
	i := &instantRunoff{
		Box:   cview.NewBox(),
		votes: make(map[int]int),
	}

	i.SetBorder(true)
	i.SetTitle("Voting")
	i.SetTitleColor(tcell.ColorLime)
	i.SetBorderColor(tcell.ColorLime)

	return i
}

func (i *instantRunoff) SetupNewVote(choices []string, submit func(map[string]int)) {
	i.l.Lock()
	defer i.l.Unlock()

	i.allCandidates = choices
	i.submitFunc = submit
}

// Draw renders the widget to text
func (i *instantRunoff) Draw(screen tcell.Screen) {
	i.l.Lock()
	defer i.l.Unlock()

	i.Box.Draw(screen)
	x, y, width, height := i.GetInnerRect()

	// Print the candidates with cursor and votes
	for index, candidate := range i.allCandidates {
		if index >= height {
			break
		}

		cursor := " "
		if i.cursor == index {
			cursor = ">"
		}

		transferrableVote := []byte(" ")
		if transfer, ok := i.votes[index]; ok {
			transferrableVote = strconv.AppendInt([]byte{}, int64(transfer), 10)
		}

		line := fmt.Sprintf("%v [%v] %v", cursor, string(transferrableVote), candidate)
		cview.Print(screen, []byte(line), x, y+index, width, cview.AlignLeft, tcell.ColorLime)
	}

	// Print submit button
	buttonText := "<Submit>"
	style := tcell.StyleDefault.Foreground(tcell.ColorLime)
	if i.cursor == len(i.allCandidates) {
		style = style.Background(tcell.ColorLime).Foreground(tcell.ColorBlack)
	}

	cview.PrintStyle(screen, []byte(buttonText), x, y+len(i.allCandidates)+1, width, cview.AlignCenter, style)

	// Print the status message
	cview.Print(screen, []byte(i.currentStatus), x, y+len(i.allCandidates)+2, width, cview.AlignCenter, tcell.ColorRed)
}

// InputHandler allows user to input things (duh)
func (i *instantRunoff) InputHandler() func(event *tcell.EventKey, setFocus func(p cview.Primitive)) {
	return i.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p cview.Primitive)) {
		switch event.Key() {
		case tcell.KeyUp:
			i.cursor--

			// Check whether we have gone before the first item
			if i.cursor < 0 {
				i.cursor = 0
			}

		case tcell.KeyDown:
			i.cursor++

			// Check whether we have gone past the button
			if i.cursor > len(i.allCandidates) {
				i.cursor = len(i.allCandidates)
			}

		case tcell.KeyEnter:
			if i.cursor == len(i.allCandidates) {
				i.currentStatus = ""

				voteMap := make(map[string]int)
				if len(i.votes) < len(i.allCandidates) {
					i.currentStatus = "Please number every candidate"
					return
				}

				usedTransfers := make([]int, len(i.votes)+1)
				for _, v := range i.votes {
					if usedTransfers[v] == v {
						i.currentStatus = "Votes are not transferrable"
						return
					}

					usedTransfers[v] = v
				}

				for k, v := range i.allCandidates {
					voteCount, vote := i.votes[k]
					if !vote {
						voteCount = 0
					}

					voteMap[v] = voteCount
				}

				i.submitFunc(voteMap)
			}
		}

		// Only allow number inputs for votes
		if unicode.IsNumber(event.Rune()) {
			num, _ := strconv.Atoi(string(event.Rune()))

			// Only allow inputs in [1, number of candidates]
			if num <= len(i.allCandidates) && num != 0 {
				i.votes[i.cursor] = num
			}
		}
	})
}

// MouseHandler allows users to click on options like peasants
func (i *instantRunoff) MouseHandler() func(action cview.MouseAction, event *tcell.EventMouse, setFocus func(p cview.Primitive)) (consumed bool, capture cview.Primitive) {
	return i.WrapMouseHandler(func(action cview.MouseAction, event *tcell.EventMouse, setFocus func(p cview.Primitive)) (consumed bool, capture cview.Primitive) {
		if action == cview.MouseLeftClick && i.InRect(event.Position()) {
			setFocus(i)

			// Allow clicking on candidates
			_, innerY, _, _ := i.GetInnerRect()
			_, y := event.Position()
			for index := range i.allCandidates {
				if y-innerY == index {
					i.cursor = index
					return true, capture
				}
			}

			// Allow clicking on submit button
			if y-innerY == len(i.allCandidates)+1 {
				i.cursor = len(i.allCandidates)

				voteMap := make(map[string]int)
				usedTransfers := make([]int, len(i.votes)+1)
				if len(voteMap) < len(i.allCandidates) {
					i.currentStatus = "Please number every candidate"
					return true, capture
				}

				for _, v := range i.votes {
					if usedTransfers[v] == v {
						i.currentStatus = "Votes are not transferrable"
						return true, capture
					}

					usedTransfers[v] = v
				}

				for k, v := range i.allCandidates {
					voteCount, vote := i.votes[k]
					if !vote {
						voteCount = 0
					}

					voteMap[v] = voteCount
				}

				i.submitFunc(voteMap)

				return true, capture
			}

			return true, capture
		}

		return false, capture
	})
}
