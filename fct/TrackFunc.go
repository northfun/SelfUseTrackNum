package fct

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (l *StLogType) key() string {
	return fmt.Sprintf("%v_%v_%v", l.Time, l.Branch, l.User)
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
	return fmt.Sprintf("key:%v,user:%v,branch:%v,addedParams:%v", l.key(), l.User, l.Branch, pstr)
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

func (up *UsedParamType) ParamSlc() []uint {
	slc := make([]uint, len(*up))
	var i uint
	for k, _ := range *up {
		slc[i] = k
		i++
	}
	return slc
}

type TrackLogType map[string]StLogType // map[key]

type TrackFuncManager struct {
	trackNum map[string]UsedParamType // map[cmdNum]map[usedParam]
	trackLog TrackLogType             // map[cmdNum][]logs
}

// 返回冲突信息
func (m *TrackFuncManager) RefreshTrack(rev *def.TrackRefresh) (map[string][]string, map[string][]uint, string) {
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
		if used, ok := m.trackNum[k]; ok {
			for i := range v {
				if info, ok := used[v[i]]; ok {
					// conflict
					slc = append(slc, fmt.Sprintf("%v:used here:%v", v[i], info.key()))
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
			m.trackNum[k] = p
		}
		if len(slc) > 0 {
			cflct[k] = slc
		}
		if len(okslc) > 0 {
			addok[k] = okslc
		}
	}
	var key string
	if len(addok) > 0 {
		log.AddParams = addok
		key = log.key()
		m.trackLog[key] = log
		m.saveData() // TODO
	}
	return cflct, addok, key
}

func (m *TrackFuncManager) UsedParam(cmd string) []uint {
	if pmap, ok := m.trackNum[cmd]; ok {
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

func (m *TrackFuncManager) GetParams(cmd string) UsedParamType {
	if pmap, ok := m.trackNum[cmd]; ok {
		return pmap
	}
	return nil
}

// TODO
func (m *TrackFuncManager) DelTrack(key string) string {
	if track, ok := m.trackLog[key]; ok {
		for k, v := range track.AddParams {
			if t, ok := m.trackNum[k]; ok {
				for i := range v {
					delete(t, v[i])
				}
			}
		}
		delete(m.trackLog, key)
		m.saveData()
		return track.toString()
	}
	return "not found"
}

func (m *TrackFuncManager) GetAllSlc() map[string][]uint {
	mp := make(map[string][]uint)
	for k, v := range m.trackNum {
		mp[k] = v.ParamSlc()
	}
	return mp
}

func (m *TrackFuncManager) Init() bool {
	return m.initData()
}

func (m *TrackFuncManager) initData() bool {
	m.trackNum = make(map[string]UsedParamType)
	m.trackLog = make(map[string]StLogType)
	if ddbuf, err := ioutil.ReadFile(LOG_FILE); err == nil {
		if err := json.Unmarshal(ddbuf, &m.trackLog); err != nil {
			fmt.Printf("unmarshal file:%v err:%v\n", LOG_FILE, err)
			return false
		}
		var ok bool
		var st map[uint]*StLogType
		for _, v := range m.trackLog {
			for key, pslc := range v.AddParams {
				if st, ok = m.trackNum[key]; !ok {
					st = make(map[uint]*StLogType)
					m.trackNum[key] = st
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

func (m *TrackFuncManager) saveData() {
	if bys, err := json.Marshal(&m.trackLog); err == nil {
		ioutil.WriteFile(LOG_FILE, bys, 0644)
	} else {
		fmt.Println("json save err:", err)
	}

	//	if bys, err := json.Marshal(&m.trackLog); err == nil {
	//		ioutil.WriteFile(LOG_FILE, bys, 0644)
	//	} else {
	//		fmt.Println("json save err")
	//	}
}
