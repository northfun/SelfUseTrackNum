package main

func main() {
	mp := make(map[string][]uint)
	mp["SceneCmd"] = []uint{23, 24, 25}
	mp["UserCmd"] = []uint{1, 2}
	ReqRefresh("fbs", "ndr", mp)
	ReqQuest("User")
	mp["SceneCmd"] = []uint{23, 24, 28}
	mp["UserCmd"] = []uint{1, 2, 5, 6}
	mp["ActCmd"] = []uint{5, 6}
	ReqRefresh("ff", "tw", mp)
	ReqQuest("User")
}
