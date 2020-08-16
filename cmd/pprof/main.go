package main

import (
	"github.com/TeaWeb/build/internal/teaweb"
	"github.com/iwind/TeaGo/logs"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		err := http.ListenAndServe("127.0.0.1:6060", nil)
		if err != nil {
			logs.Println("[error]" + err.Error())
		}
	}()
	teaweb.Start()
}
