package main

import "github.com/snowmagic1/gedis/server"

func main() {
	server := new(server.App)
	server.Start()
}
