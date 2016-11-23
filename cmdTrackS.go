package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"
	"tracks/def"
)

const LOG_MAX_NUM = 1 << 8
const LOG_FILE = "track.log"

type StLogType struct {
	def.QuestTrackCmd
	time int64
}

func (l *StLogType) toString() string {
	return fmt.Sprintf("%v,user:%v,branch:%v,cmd:%v,param:%v", l.time, l.User, l.Branch, l.Cmd, l.Param)
}

var trackNum map[string]map[int]bool // map[cmdNum]map[usedParam]
type TrackLogType map[string][]*StLogType

var trackLog TrackLogType // map[cmdNum][]logs

func addLog(rev *def.QuestTrackCmd, send *def.AnswerTrackCmd) {
	var logs []string
	var ok bool
	var log StLogType
	log.time = time.Now().Unix()
	log.QuestTrackCmd = *rev
	if logs, ok = trackLog[rev.Cmd]; !ok {
		logs = make([]*StLogType, 0, LOG_MAX_NUM)
	}
	logs = append(logs, &log)
	trackLog[rev.Cmd] = logs
	saveData() // TODO
	if send != nil {
		send.Res = log.toString()
	}
	//appendToFile(log)
}

func usedParam(cmd string) []int {
	if pmap, ok := trackNum[cmd]; ok {
		var i int
		params := make([]int, len(pmap))
		for p, _ := range pmap {
			params[i] = p
			i++
		}
		return params
	}
	return nil
}

func checkRevParams(rev *def.QuestTrackCmd) bool {
	if pmap, ok := trackNum[rev.Cmd]; ok {
		for i := range rev.Params {
			if _, find := pmap[rev.Params[i]]; find {
				return false
			}
		}
	}
	return true
}

func dealRev(rev *def.QuestTrackCmd, send *def.AnswerTrackCmd) {
	if rev == nil || send == nil {
		return
	}
	send.Cmd = rev.Cmd
	var ok bool
	switch rev.Serve {
	case "Quest":
		send.UsedParam = usedParam(rev.Cmd)
	case "Add":
		if curParam, ok := trackNum[rev.Cmd]; ok {
			if curParam == LOG_MAX_NUM {
				send.Res = fmt.Sprintf("[%v],starve:%v", def.COMMIT_LOG_ERR, send.Cmd)
				return
			}
			if curParam+1 == rev.Param {
				trackNum[rev.Cmd] = rev.Param
				addLog(rev, send)
			} else {
				send.Res = fmt.Sprintf("[%v],cmd:%v,nextParam:%v", def.COMMIT_LOG_ERR, rev.Cmd, curParam+1)
			}
		} else if rev.Param == 1 {
			trackNum[rev.Cmd] = 1
			addLog(rev, send)
		} else {
			send.Res = fmt.Sprintf("[%v],cmd:%v,nextParam:%v", def.COMMIT_LOG_ERR, rev.Cmd, 1)
		}
	default:
		send.Res = fmt.Sprintf("[%v],undefined serve:%v", rev.Serve)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	bts := make([]byte, 1<<10)
	num, err := conn.Read(bts)
	if err != nil {
		return
	}
	var rev def.QuestTrackCmd
	json.Unmarshal(bts[:num], &rev)
	var send def.AnswerTrackCmd

	dealRev(&rev, &send)
	if sendJson, err := json.Marshal(&send); err == nil {
		fmt.Fprintf(conn, string(sendJson))
	} else {
		fmt.Println("send json err:", err)
	}
	fmt.Println("server done :", send)
}

func initData() {
	trackNum = make(map[string]int)
	trackLog = make(map[string][]string)
	if ddbuf, err := ioutil.ReadFile(LOG_FILE); err == nil {
		if err := json.Unmarshal(ddbuf, trackLog); err == nil {
			fmt.Println("initData unmarshal err")
		} else {
			fmt.Printf("unmarshal file-%v err:%v\n", LOG_FILE, err)
		}
	} else {
		fmt.Printf("read file-%v err:%v\n", LOG_FILE, err)
	}

}

func saveData() {
	if bys, err := json.Marshal(trackLog); err == nil {
		ioutil.WriteFile(LOG_FILE, bys, 0644)
	} else {
		fmt.Println("json save err")
	}
}

func main() {
	initData()
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", def.SERVER_PORTS))
	if err != nil {
		panic(err)
	}
	defer saveData()
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
