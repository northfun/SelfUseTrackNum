package client

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

func ReqDel(key string) {
	var send def.TrackReqDelTrack
	send.Init()
	send.Key = key
	askServer(&send)
}

func askServer(m def.Message_itfc) {
	saddr := fmt.Sprintf("%v:%v", def.SERVER_IP, def.SERVER_PORTS)
	conn, err := net.Dial("tcp", saddr)
	if err != nil {
		fmt.Println("连接错误,addr:", saddr)
		return
	}
	defer conn.Close()
	bts := def.PackCmd(m)
	_, err = conn.Write(bts)
	if err != nil {
		fmt.Println("write错误,err:", err)
		return
	}
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println("read错误,err:", err)
		return
	}
	usage, ret := def.UnPackCmd(result)
	dealRev(usage, ret)
}

func printConflicts(c map[string][]string) (count int) {
	if len(c) == 0 {
		return
	}
	fmt.Println("<<<---------------[WRN],some conflicts")
	for k, v := range c {
		fmt.Println(k)
		for i := range v {
			fmt.Println(v[i])
			count++
		}
	}
	fmt.Println("[WRN],some conflicts---------------->>>")
	return
}

func printAddOk(ok map[string][]uint) (count int) {
	if len(ok) == 0 {
		return
	}
	fmt.Println("========[OK],added ok params:========")
	for k, v := range ok {
		fmt.Println(k, v)
		for _ = range v {
			count++
		}
	}
	return
}

func dealRev(usage uint, data []byte) {
	switch usage {
	case def.MESSAGE_TYPE_RetRefresh:
		var rev def.TrackRetRefresh
		json.Unmarshal(data, &rev)
		cc := printConflicts(rev.Conflict)
		okc := printAddOk(rev.AddOk)
		fmt.Println("[total],conflicts:", cc, ",ok:", okc, ",key:", rev.Key)
	case def.MESSAGE_TYPE_RetQuest:
		var rev def.TrackRetQuest
		json.Unmarshal(data, &rev)
		fmt.Println("-------------list-------------")
		fmt.Println(rev.Data)
	case def.MESSAGE_TYPE_RetDelTrack:
		var rev def.TrackRetDelTrack
		json.Unmarshal(data, &rev)
		fmt.Println("del return:")
		fmt.Println(rev.Res)
	}
}
