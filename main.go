package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	port   int
	server string
	typ    string
)

func serverHandler() {
	packetConn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		println("listen packet failed, err:", err.Error())
		return
	}

	for {
		buf := make([]byte, 4096)
		n, src, err := packetConn.ReadFrom(buf)
		if err != nil {
			println("read packet failed, err:", err.Error())
			return
		}
		println("read packet:", string(buf[:n]))

		_, err = packetConn.WriteTo([]byte("response"), src)
		if err != nil {
			println("write packet failed, err:", err.Error())
			return
		}
	}
}

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