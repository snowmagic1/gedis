package server

func getCmdhandler(c *Client) error {
	args := c.args
	if len(args) != 1 {
		return ErrCmdParams
	}

	val, err := c.app.db.Get(args[0], nil)
	if err != nil {
		return err
	}

	if val == nil {
		c.writer.writeBulk(nil)
	} else {
		c.writer.writeBulk(val)
	}

	return nil
}

func setCmdhandler(c *Client) error {
	args := c.args
	if len(args) != 2 {
		return ErrCmdParams
	}

	err := c.app.db.Put(args[0], args[1], nil)
	if err == nil {
		c.writer.writeStatus(OK)
	}

	return nil
}

func init() {
	RegisterCmdHandler("Get", getCmdhandler)
	RegisterCmdHandler("Set", setCmdhandler)
}
