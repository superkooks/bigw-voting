package commands

import "strings"

var commands map[string]func([]string)

// RegisterAll registers all the commands present using RegisterCommand
func RegisterAll() {

}

// RegisterCommand takes a name and a callback
func RegisterCommand(name string, cmd func([]string)) {

}

// Parse parses a command and calls any callbacks
func Parse(cmd string) {
	split := strings.Split(cmd, " ")
	for k, v := range commands {
		if split[0] == k {
			v(split[1:])
			return
		}
	}
}
