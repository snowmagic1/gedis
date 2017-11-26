package server

import (
	"bufio"
	"fmt"
	"net"

	"github.com/snowmagic1/gedis/utils"
)

type Client struct {
	conn   net.Conn
	reader *utils.RespReader
	writer *responseWriter
	args   [][]byte
}

func NewClient(conn net.Conn) *Client {
	c := &Client{
		conn: conn,
	}

	br := bufio.NewReaderSize(conn, 4*1024)
	c.reader = utils.NewRespReader(br)

	w := new(responseWriter)
	w.b = bufio.NewWriterSize(conn, 4*1024)
	c.writer = w

	return c
}

func (c *Client) Run() {
	req, err := c.reader.ParseRequest()
	if err != nil {
		fmt.Println("failed to read, ", err)
		return
	}

	cmd := string(req[0])
	c.args = req[1:]

	if handler := registeredCmds[cmd]; handler == nil {
		err = ErrNotFound
	} else {
		err = handler(c)
	}

	if err != nil {
		c.writer.writeError(err)
	}
	c.writer.flush()
}

type responseWriter struct {
	b *bufio.Writer
}

func (rw *responseWriter) writeError(err error) {
	rw.b.Write([]byte("-"))
	if err != nil {
		rw.b.Write([]byte(err.Error()))
	}
	rw.b.Write(Delims)
}

func (rw *responseWriter) flush() {
	rw.b.Flush()
}
