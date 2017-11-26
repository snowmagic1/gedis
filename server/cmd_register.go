package server

type CmdHandler func(c *Client) error

var registeredCmds = make(map[string]CmdHandler)

func RegisterCmdHandler(name string, handler CmdHandler) {
	if registeredCmds[name] != nil {
		panic("already registered")
	}

	registeredCmds[name] = handler
}