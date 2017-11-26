package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/snowmagic1/gedis/utils"
)

type Client struct {
	app    *App
	conn   net.Conn
	reader *utils.RespReader
	writer *responseWriter
	args   [][]byte
}

func NewClient(app *App, conn net.Conn) *Client {
	c := &Client{
		app:  app,
		conn: conn,
	}

	br := bufio.NewReaderSize(conn, 4*1024)
	c.reader = utils.NewRespReader(br)

	w := new(responseWriter)
	w.b = bufio.NewWriterSize(conn, 4*1024)
	c.writer = w

	return c
}

func (c *Client) ProcessRequest() error {
	req, err := c.reader.ParseRequest()
	if err != nil {
		fmt.Println("failed to read, ", err)
		return err
	}

	cmd := strings.ToLower(string(req[0]))
	c.args = req[1:]

	log.Printf("[%v]\n", cmd)
	for _, arg := range c.args {
		log.Printf("  %v\n", string(arg))
	}

	if handler := registeredCmds[cmd]; handler == nil {
		err = ErrNotFound
	} else {
		err = handler(c)
	}

	if err != nil {
		c.writer.writeError(err)
	}
	c.writer.flush()

	return nil
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

func (rw *responseWriter) writeBulk(b []byte) {
	rw.b.WriteByte('$')
	if b == nil {
		rw.b.Write(NullBulk)
	} else {
		rw.b.Write([]byte(strconv.Itoa(len(b))))
		rw.b.Write(Delims)
		rw.b.Write(b)
	}

	rw.b.Write(Delims)
}

func (rw *responseWriter) writeStatus(status string) {
	rw.b.WriteByte('+')
	rw.b.Write([]byte(status))
	rw.b.Write(Delims)
}

func (rw *responseWriter) flush() {
	rw.b.Flush()
}
