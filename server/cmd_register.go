package server

import "strings"

type CmdHandler func(c *Client) error

var registeredCmds = make(map[string]CmdHandler)

func RegisterCmdHandler(name string, handler CmdHandler) {
	lowerName := strings.ToLower(name)
	if registeredCmds[lowerName] != nil {
		panic("already registered")
	}

	registeredCmds[lowerName] = handler
}
