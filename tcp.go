package main

import (
	"log"
	"net"
)

var TcpPool = make(map[string]map[int]net.Conn, 10)

func InitTcp(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleTcp(conn)

	}
}

func handleTcp(conn net.Conn) {
	defer conn.Close()
	var headerBuf = make([]byte, 8)
	var buf = make([]byte, 1024)
	for {
		n, err := conn.Read(headerBuf)
		if err != nil {
			log.Println("read header error:", err)
			return
		}
		if n != 8 {
			return
		}
		var size int
		if size > len(buf) {
			buf = make([]byte, size)
		}
		_, err = conn.Read(buf[:size])
		if err != nil {
			log.Println("read error:", err)
			return
		}
		for _, targets := range TcpPool {
			for _, target := range targets {
				_, err := target.Write(buf[:n])
				if err != nil {
					log.Println("write error:", err)
					continue
				}
				break
			}
			break
		}
	}
}
