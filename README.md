# SelfUseTrackNum
for self use

cd $GOPATH/src

git clone http://192.168.1.201:3000/fbs/CmdTrackTool.git tracks

Server:

go run ./server/*.go &
Client:

cd $workspace
go run ./tools/TranslateCmdCheck.go [init(first need)]
