package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"
)

var (
	port   int
	server string
	typ    string
	count  int
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
		println("read packet:", string(buf[:n]), ", src addr:", src.String())

		_, err = packetConn.WriteTo([]byte("response"), src)
		if err != nil {
			println("write packet failed, err:", err.Error())
			return
		}
	}
}

func main() {

	flag.IntVar(&port, "p", 8443, "server port")
	flag.StringVar(&server, "s", "", "server address")
	flag.StringVar(&typ, "t", "client", "programe type client/server")
	flag.IntVar(&count, "c", 1, "count")


	flag.Parse()

	if typ == "server" {
		serverHandler()
		return
	}

	if server == "" {
		println("with no server")
		return
	}

	packetConn, err := net.ListenPacket("udp", "")
	if err != nil {
		println("listen packet failed, err:", err.Error())
		return
	}

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		println("resolve udp addr failed, err:", err.Error())
		return
	}

	for {
		startTime := time.Now()
		_, err = packetConn.WriteTo([]byte("request"), udpAddr)
		if err != nil {
			println("write packet failed, err:", err.Error())
			return
		}

		buf := make([]byte, 4096)
		packetConn.SetReadDeadline(time.Now().Add(3 * time.Second))
		n, _, err := packetConn.ReadFrom(buf)
		if err != nil {
			println(time.Now().Format(http.TimeFormat), ", read packet failed, err:", err.Error())
			continue
		}
		endTime := time.Now()
		println(endTime.Format(http.TimeFormat), ", read packet:", string(buf[:n]), ", used time:", endTime.Sub(startTime).Milliseconds(), "ms")

		time.Sleep(time.Second)
	}

}