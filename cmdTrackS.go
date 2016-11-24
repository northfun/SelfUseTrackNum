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
	User, Branch string
	AddParams    map[string][]int
	Time         int64
}

func (l *StLogType) init(rev *def.TrackRefresh) {
	log.Time = time.Now().Unix()
	log.Branch = rev.Branch
	log.User = rev.User
}

func (l *StLogType) toString() string {
	return fmt.Sprintf("%v,user:%v,branch:%v,addedParams:%v", l.Time, l.User, l.Branch, l.AddParams)
}

type UsedParamType map[int]*StLogType

func (up *UsedParamType) addParams(p []int, plog *StLogType) {
	for i := range p {
		(*up)[p[i]] = plog
	}
}

func (up *UsedParamType) checkParams(p []int) bool {
	for i := range p {
		if _, find := (*up)[p[i]]; find {
			return false
		}
	}
	return true
}

func (up *UsedParamType) paramSlc() []int {
	slc := make([]int, len(*up))
	var i int
	for k, _ := range *up {
		slc[i] = k
		i++
	}
	return slc
}

var trackNum map[string]UsedParamType // map[cmdNum]map[usedParam]
type TrackLogType map[string][]StLogType

var trackLog TrackLogType // map[cmdNum][]logs

// 返回冲突信息
func refreshTrack(rev *def.TrackRefresh) (map[string][]string, map[string][]int) {
	var logs []StLogType
	var ok bool
	var log StLogType
	log.init(rev)
	cflct := make(map[string][]string)
	addok := make(map[string][]int)
	for k, v := range rev.Data {
		slc := make([]string)
		okslc := make([]int)
		if used, ok := trackNum[k]; ok {
			for i := range v {
				if info, ok := used[v[i]]; ok {
					// conflict
					slc = append(slc, info.toString())
				} else {
					okslc = append(okslc, v[i])
					used.addParams([]int(v[i]))
				}
			}
		} else {
			for i := range v {
				okslc = append(okslc, v[i])
			}
			var p UsedParamType
			p.addParams(v)
			trackNum[k] = p
		}
		if len(okslc) > 0 {
			cflct[k] = slc
		}
		if len(slc) > 0 {
			addok[k] = okslc
		}
	}
	if len(addok) > 0 {
		logs.AddParams = addok
		if logs, ok = trackLog[rev.Cmd]; !ok {
			logs = make([]StLogType, 0, LOG_MAX_NUM)
		}
		logs = append(logs, log)
		trackLog[rev.Cmd] = logs
		saveData() // TODO
	}
	return cflct, addok
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

func getParams(cmd string) UsedParamType {
	if pmap, ok := trackNum[cmd]; ok {
		return pmap
	}
	return nil
}

func handleConnection(conn net.Conn) {
	var user TrackSUser
	user.conn = conn
	user.dealRev()
	fmt.Println("deal from ", conn.RemoteAddr())
}

func initData() bool {
	trackNum = make(map[string]UsedParamType)
	trackLog = make(map[string][]StLogType)
	if ddbuf, err := ioutil.ReadFile(LOG_FILE); err == nil {
		if err := json.Unmarshal(ddbuf, &trackLog); err != nil {
			fmt.Printf("unmarshal file:%v err:%v\n", LOG_FILE, err)
			return false
		}
		for k, v := range trackLog {
			st := make(map[int]*StLogType)
			for i := range v {
				for j := range v[i].Params {
					st[v[i].Params[j]] = &v[i]
				}
			}
			trackNum[k] = st
		}
	} else {
		fmt.Printf("read file-%v err:%v\n", LOG_FILE, err)
	}
	return true
}

func saveData() {
	if bys, err := json.Marshal(&trackLog); err == nil {
		ioutil.WriteFile(LOG_FILE, bys, 0644)
	} else {
		fmt.Println("json save err")
	}
}

func main() {
	if !initData() {
		return
	}
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
