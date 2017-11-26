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

func init() {
	RegisterCmdHandler("Get", getCmdhandler)
}
