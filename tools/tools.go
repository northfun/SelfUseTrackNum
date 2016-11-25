package tools

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"time"
)

func DumpStack() {
	if err := recover(); err != nil {
		var buf bytes.Buffer
		bs := make([]byte, 1<<12)
		num := runtime.Stack(bs, false)
		buf.WriteString(fmt.Sprintf("Panic: %s\n", err))
		buf.Write(bs[:num])
		dumpName := "log/dump_" + time.Now().Format("20060102-150405")
		nerr := ioutil.WriteFile(dumpName, buf.Bytes(), 0644)
		if nerr != nil {
			log.Println("write dump file error", nerr)
			log.Println(buf.Bytes())
		}
	}
}
