package main

import (
	"flag"
)

var (
	port   int
	server string
	typ    string
)

func main() {

	flag.IntVar(&port, "p", 443, "server port")
	flag.StringVar(&server, "s", "", "server address")
	flag.StringVar(&typ, "t", "client", "programe type client/server")


	flag.Parse()

	if typ == "server" {
		serverHandler()
		return
	}

	if server == "" {
		println("with no server")
		return
	}
}