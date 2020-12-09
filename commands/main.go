package commands

import (
	"strings"
)

var commands = make(map[string]func([]string))

// RegisterAll registers all the commands present using RegisterCommand
func RegisterAll() {
	RegisterCommand("connect", CommandConnect)
	RegisterCommand("nick", CommandNick)
}

// RegisterCommand takes a name and a callback
func RegisterCommand(name string, cmd func([]string)) {
	commands[name] = cmd
}

// Parse parses a command and calls any callbacks
func Parse(cmd string) {
	split := strings.Split(cmd, " ")
	for k, v := range commands {
		if split[0] == k {
			go v(split[1:])
			return
		}
	}
}
