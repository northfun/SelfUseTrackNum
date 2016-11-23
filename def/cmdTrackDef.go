package def

const SERVER_PORTS = 10777
const SERVER_IP = "127.0.0.1"

const COMMIT_LOG_ERR = "err"
const COMMIT_LOG_OK = "ok"

type QuestTrackCmd struct {
	Serve        string // "Quest" "Add"
	Cmd          string // "CmdType"
	Params       []int  // "ParamType"
	User, Branch string
}

type AnswerTrackCmd struct {
	Cmd       string
	UsedParam []int
	Res       string
}
