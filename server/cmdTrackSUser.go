package main

import (
	"encoding/json"
	"fmt"
	"net"
	"tracks/def"
)

type TrackSUser struct {
	conn net.Conn
}

func (u *TrackSUser) sendToMe(m def.Message_itfc) {
	u.conn.Write(def.PackCmd(m))
}

func (u *TrackSUser) do() {
	defer u.conn.Close()
	bts := make([]byte, 1<<10)
	num, err := u.conn.Read(bts)
	if err != nil {
		fmt.Println("read err 1", err)
		return
	}
	fmt.Println("deal:", num)
	usage, data := def.UnPackCmd(bts)
	u.dealRev(usage, data)
}

func (u *TrackSUser) dealRev(usage uint, data []byte) {
	switch usage {
	case def.MESSAGE_TYPE_Quest:
		var rev def.TrackQuest
		json.Unmarshal(data, &rev)
		var send def.TrackRetQuest
		send.Init()
		if len(rev.Cmd) > 0 {
			if ps := getParams(rev.Cmd); ps != nil {
				send.Data[rev.Cmd] = ps.paramSlc()
			}
		} else {
			for k, v := range trackNum {
				send.Data[k] = v.paramSlc()
			}
		}
		u.sendToMe(&send)
	case def.MESSAGE_TYPE_Refresh:
		var rev def.TrackRefresh
		json.Unmarshal(data, &rev)
		var send def.TrackRetRefresh
		send.Init()
		u.sendToMe(&send)
	}
}
