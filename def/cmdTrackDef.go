package def

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

const SERVER_PORTS = 10777
const SERVER_IP = "192.168.1.236"

const COMMIT_LOG_ERR = "err"
const COMMIT_LOG_OK = "ok"

func BinRead(buf *bytes.Buffer, data interface{}) {
	binary.Read(buf, binary.LittleEndian, data)
}

func BinWrite(buf *bytes.Buffer, data interface{}) {
	binary.Write(buf, binary.LittleEndian, data)
}

func UnPackCmd(rev []byte) (uint, []byte) {
	var usage int16
	buf := bytes.NewBuffer(rev)
	BinRead(buf, &usage)
	return uint(usage), buf.Bytes()
}

func PackCmd(m Message_itfc) []byte {
	var buf bytes.Buffer
	BinWrite(&buf, uint16(m.Usage()))
	var bts []byte
	var err error
	if bts, err = json.Marshal(m); err == nil {
		BinWrite(&buf, bts)
	} else {
		fmt.Println("pack cmd json err:%v", err)
	}
	return buf.Bytes()
}

type Message_itfc interface {
	Usage() uint
}

type MessageBase struct {
	usage uint
}

func (b *MessageBase) Usage() uint {
	return b.usage
}

// C->S 请求查询
const MESSAGE_TYPE_Quest = 1

type TrackQuest struct {
	MessageBase
	Cmd string // "CmdType"
}

func (t *TrackQuest) Init() {
	t.usage = MESSAGE_TYPE_Quest
}

// C->S 请求更新
const MESSAGE_TYPE_Refresh = 2

type TrackRefresh struct {
	MessageBase
	Data         map[string][]uint
	User, Branch string
}

func (t *TrackRefresh) Init() {
	t.usage = MESSAGE_TYPE_Refresh
}

// S->C 响应更新
const MESSAGE_TYPE_RetRefresh = 3

type TrackRetRefresh struct {
	MessageBase
	Conflict map[string][]string // map[cmdName]conflict infos
	AddOk    map[string][]uint
	Key      string
}

func (t *TrackRetRefresh) Init() {
	t.usage = MESSAGE_TYPE_RetRefresh
}

// S->C 响应查询
const MESSAGE_TYPE_RetQuest = 4

type TrackRetQuest struct {
	MessageBase
	Data map[string][]uint
}

func (t *TrackRetQuest) Init() {
	t.usage = MESSAGE_TYPE_RetQuest
	t.Data = make(map[string][]uint)
}

// C->S 删除记录
const MESSAGE_TYPE_ReqDelTrack = 5

type TrackReqDelTrack struct {
	MessageBase
	Key string
}

func (t *TrackReqDelTrack) Init() {
	t.usage = MESSAGE_TYPE_ReqDelTrack
}

// S->C 删除记录
const MESSAGE_TYPE_RetDelTrack = 6

type TrackRetDelTrack struct {
	MessageBase
	Res string
}

func (t *TrackRetDelTrack) Init() {
	t.usage = MESSAGE_TYPE_RetDelTrack
}
