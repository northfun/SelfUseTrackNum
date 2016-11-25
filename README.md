# SelfUseTrackNum
for self use

cd $GOPATH/src

git clone https://github.com/northfun/SelfUseTrackNum.git tracks

Server:

go run ./server/*.go &

Client:

cd $workspace

go run ./tools/TranslateCmdCheck.go [init(first need)]
