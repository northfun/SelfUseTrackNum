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
	usage, data := def.UnPackCmd(bts[:num])
	u.dealRev(usage, data)
}

func (u *TrackSUser) dealRev(usage uint, data []byte) {
	switch usage {
	case def.MESSAGE_TYPE_Quest:
		var rev def.TrackQuest
		if err := json.Unmarshal(data, &rev); err != nil {
			fmt.Println("err:", err)
		}
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
		send.Conflict, send.AddOk = refreshTrack(&rev)
		u.sendToMe(&send)
	case def.MESSAGE_TYPE_ReqDelTrack:
		var rev def.TrackReqDelTrack
		json.Unmarshal(data, &rev)
		var send def.TrackRetDelTrack
		send.Init()
		send.Res = delTrack(rev.Key)
		u.sendToMe(&send)
	default:
		fmt.Println("err usage")
	}
}
