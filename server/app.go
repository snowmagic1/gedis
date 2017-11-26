package server

import (
	"fmt"
	"net"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "6379"
	CONN_TYPE = "tcp"
)

type App struct {
	db *leveldb.DB
}

func (app *App) Start() error {

	db, err := leveldb.OpenFile("./testdb", nil)
	if err != nil {
		return err
	}
	app.db = db

	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go app.handleRequest(conn)
	}
}

// Handles incoming requests.
func (app *App) handleRequest(conn net.Conn) {
	client := NewClient(app, conn)
	for {
		client.ProcessRequest()
	}
}
