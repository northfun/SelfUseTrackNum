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

func (u *TrackSUser) sendToMe(m Message_itfc) {
	u.conn.Write(def.PackCmd(m))
}

func (u *TrackSUser) dealConn() {
	defer conn.Close()
	bts := make([]byte, 1<<10)
	num, err := conn.Read(bts)
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
		json.Unmarshal(&rev)
		var send def.TrackRetQuest
		send.Init()
		if len(rev.Cmd) > 0 {
			send.Data[rev.Cmd] = getParams(rev.Cmd)
		} else {
			for k, v := range trackNum {
				send.Data[k] = v.paramSlc()
			}
		}
		u.sendToMe(&send)
	case def.MESSAGE_TYPE_Refresh:
		var rev def.TrackRefresh
		json.Unmarshal(&rev)
		var send def.TrackRetRefresh
		send.Init()
		send.Res = fmt.Sprintf("[%v],cmd:%v,should add Params:%v", def.COMMIT_LOG_ERR, rev.Cmd, rev.Params)
		u.sendToMe(&send)
	}
}
