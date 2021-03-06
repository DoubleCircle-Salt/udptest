package main

import (
	"flag"
	"fmt"
	"math/rand"
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
		go func(src net.Addr, randstr []byte) {
			_, err = packetConn.WriteTo(randstr, src)
			if err != nil {
				println("write packet failed, err:", err.Error())
				return
			}
		} (src, buf[:n])
	}
}

const RANDSTR_LENGTH = 4

func getRandstr() []byte {

	randstr := make([]byte, RANDSTR_LENGTH)
	for i := 0; i < RANDSTR_LENGTH; i++ {
		randInt := rand.Intn(36)
		if randInt < 10 {
			randstr[i] = byte(randInt + '0')
		} else {
			randstr[i] = byte(randInt + 'a' - 10)
		}
	}

	return randstr
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

	for i := 0; i < count; i++ {
		go func() {
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
				randstr := getRandstr()
				_, err = packetConn.WriteTo(randstr, udpAddr)
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
				restr := buf[:n]

				println(endTime.Format(http.TimeFormat), ", write packet:", string(randstr), ", read packet:", string(restr), ", used time:", endTime.Sub(startTime).Milliseconds(), "ms")

				time.Sleep(time.Second)
			}
		}()
	}

	time.Sleep(365 * 24 * time.Hour)
}