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
	AddParams    map[string][]uint
	Time         uint
}

func (l *StLogType) init(rev *def.TrackRefresh) {
	l.Time = uint(time.Now().Unix())
	l.Branch = rev.Branch
	l.User = rev.User
}

func (l *StLogType) toString() string {
	var pstr string
	for k, v := range l.AddParams {
		if len(pstr) == 0 {
			pstr = fmt.Sprintf("%v:%v", k, v)
		} else {
			pstr = fmt.Sprintf("%v,%v:%v", pstr, k, v)
		}
	}
	return fmt.Sprintf("%v,user:%v,branch:%v,addedParams:%v", l.Time, l.User, l.Branch, pstr)
}

type UsedParamType map[uint]*StLogType

func (up *UsedParamType) init() {
	(*up) = make(map[uint]*StLogType)
}

func (up *UsedParamType) addParams(p []uint, plog *StLogType) {
	for i := range p {
		(*up)[p[i]] = plog
	}
}

func (up *UsedParamType) checkParams(p []uint) bool {
	for i := range p {
		if _, find := (*up)[p[i]]; find {
			return false
		}
	}
	return true
}

func (up *UsedParamType) paramSlc() []uint {
	slc := make([]uint, len(*up))
	var i uint
	for k, _ := range *up {
		slc[i] = k
		i++
	}
	return slc
}

var trackNum map[string]UsedParamType // map[cmdNum]map[usedParam]
type TrackLogType map[uint]StLogType  // map[时间戳]

var trackLog TrackLogType // map[cmdNum][]logs

// 返回冲突信息
func refreshTrack(rev *def.TrackRefresh) (map[string][]string, map[string][]uint) {
	var log StLogType
	log.init(rev)
	cflct := make(map[string][]string)
	addok := make(map[string][]uint)
	for k, v := range rev.Data {
		if len(v) == 0 {
			continue
		}
		slc := make([]string, 0)
		okslc := make([]uint, 0)
		if used, ok := trackNum[k]; ok {
			for i := range v {
				if info, ok := used[v[i]]; ok {
					// conflict
					slc = append(slc, fmt.Sprintf("%v:used here:%v", v[i], info.toString()))
				} else {
					okslc = append(okslc, v[i])
					used.addParams([]uint{v[i]}, &log)
				}
			}
		} else {
			for i := range v {
				okslc = append(okslc, v[i])
			}
			var p UsedParamType
			p.init()
			p.addParams(v, &log)
			trackNum[k] = p
		}
		if len(slc) > 0 {
			cflct[k] = slc
		}
		if len(okslc) > 0 {
			addok[k] = okslc
		}
	}
	if len(addok) > 0 {
		log.AddParams = addok
		trackLog[log.Time] = log
		saveData() // TODO
	}
	return cflct, addok
}

func usedParam(cmd string) []uint {
	if pmap, ok := trackNum[cmd]; ok {
		var i uint
		params := make([]uint, len(pmap))
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

// TODO
func delTrack(key uint) string {
	return "404"
}

// TODO 可改成并发
func handleConnection(conn net.Conn) {
	var user TrackSUser
	user.conn = conn
	user.do()
	fmt.Println("deal from ", conn.RemoteAddr())
}

func initData() bool {
	trackNum = make(map[string]UsedParamType)
	trackLog = make(map[uint]StLogType)
	if ddbuf, err := ioutil.ReadFile(LOG_FILE); err == nil {
		if err := json.Unmarshal(ddbuf, &trackLog); err != nil {
			fmt.Printf("unmarshal file:%v err:%v\n", LOG_FILE, err)
			return false
		}
		var ok bool
		var st map[uint]*StLogType
		for _, v := range trackLog {
			for key, pslc := range v.AddParams {
				if st, ok = trackNum[key]; !ok {
					st = make(map[uint]*StLogType)
					trackNum[key] = st
				}
				for j := range pslc {
					st[pslc[j]] = &v
				}
			}
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
