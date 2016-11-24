package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"tracks/def"
)

func parseParamSlc(p string) []int {
	slc := strings.Split(p, ",")
	pslc := make([]int, len(slc))
	var err error
	for i := range slc {
		if pslc[i], err = strconv.Atoi(slc[i]); err != nil {
			fmt.Println("parse params ", p, " error:", err)
			return nil
		}
	}
	return pslc
}

func parseArgs(send *def.QuestTrackCmd) bool {
	if send == nil {
		return false
	}
	argNum := len(os.Args)
	if argNum < 2 {
		fmt.Println("asktrack <usage> <CmdType> [<ParamType>]") // kind of complex
		return false
	}
	switch os.Args[1] {
	case "quest":
		send.Serve = "Quest"
		send.Cmd = os.Args[2]
		if argNum < 3 {
			fmt.Println("need input \"cmd\"")
			return false
		}
	case "add":
		if argNum < 4 {
			fmt.Println("need input \"cmd\" & \"param\"")
			return false
		}
		send.Serve = "Add"
		send.Cmd = os.Args[2]
		send.Params = parseParamSlc(os.Args[3])
	default:
		return false
	}
	return true
}

func parseRev(rev *def.AnswerTrackCmd) {
	ok := strings.Contains(rev.Res, def.COMMIT_LOG_OK)
	if ok {
		fmt.Println("succeed!")
	} else {
		fmt.Println("failed:", rev.Res)
	}
	fmt.Println("rev:", rev)
}

func main() {
	var send def.QuestTrackCmd
	if !parseArgs(&send) {
		return
	}
	sendJson, err := json.Marshal(&send)
	if err != nil {
		fmt.Println("json错误")
		return
	}
	fmt.Println("send:", string(sendJson))
	saddr := fmt.Sprintf("%v:%v", def.SERVER_IP, def.SERVER_PORTS)
	conn, err := net.Dial("tcp", saddr)
	if err != nil {
		fmt.Println("连接错误,addr:", saddr)
		return
	}
	_, err = conn.Write(sendJson)
	result, err := ioutil.ReadAll(conn)
	var rev def.AnswerTrackCmd
	json.Unmarshal(result, &rev)
	parseRev(&rev)
}
