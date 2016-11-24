package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"tracks/def"
)

func ReqRefresh(u, b string, d map[string][]uint) {
	var send def.TrackRefresh
	send.Init()
	send.User = u
	send.Branch = b
	send.Data = d
	askServer(&send)
}

func ReqQuest(c string) {
	var send def.TrackQuest
	send.Init()
	send.Cmd = c
	askServer(&send)
}

func askServer(m def.Message_itfc) {
	saddr := fmt.Sprintf("%v:%v", def.SERVER_IP, def.SERVER_PORTS)
	conn, err := net.Dial("tcp", saddr)
	if err != nil {
		fmt.Println("连接错误,addr:", saddr)
		return
	}
	bts := def.PackCmd(m)
	_, err = conn.Write(bts)
	result, err := ioutil.ReadAll(conn)
	usage, ret := def.UnPackCmd(result)
	dealRev(usage, ret)
}

func dealRev(usage uint, data []byte) {
	switch usage {
	case def.MESSAGE_TYPE_RetRefresh:
		var rev def.TrackRetRefresh
		json.Unmarshal(data, &rev)
		fmt.Println(rev)
	case def.MESSAGE_TYPE_RetQuest:
		var rev def.TrackRetQuest
		json.Unmarshal(data, &rev)
		fmt.Println(rev)
	}
}
