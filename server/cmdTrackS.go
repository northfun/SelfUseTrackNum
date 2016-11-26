package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"tracks/def"
	"tracks/fct"
	"tracks/tools"
)

// TODO 可改成并发
func handleConnection(conn net.Conn) {
	var user TrackSUser
	user.conn = conn
	user.do()
	fmt.Println("deal from ", conn.RemoteAddr())
}

var wg_loop sync.WaitGroup

func loop(ln net.Listener) {
	defer tools.DumpStack()
	defer wg_loop.Done()
	for {
		fmt.Println("wait ...")
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("get client connection error: ", err)
		}
		handleConnection(conn)
		// go handleConnection(conn)
	}
}

type ServerData struct {
	track fct.TrackFuncManager
}

func (d *ServerData) init() bool {
	if !d.track.Init() {
		return false
	}
	// ...
	return true
}

var sdt ServerData

func main() {
	if !sdt.init() {
		return
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", def.SERVER_PORTS))
	if err != nil {
		panic(err)
	}
	for {
		wg_loop.Add(1)
		go loop(ln)
		wg_loop.Wait()
	}
}
