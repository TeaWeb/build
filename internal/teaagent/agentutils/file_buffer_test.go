package agentutils

import (
	"fmt"
	"github.com/iwind/TeaGo/logs"
	"strings"
	"testing"
	"time"
)

func TestNewFileBuffer(t *testing.T) {
	logs.Println("starting ...")
	buf := NewFileBuffer("a")
	buf.Debug()
	longText := []byte(strings.Repeat("abc", 10000) + "\n")
	go func() {
		i := 0
		for {
			i++
			go buf.Write(longText)
			time.Sleep(2 * time.Second)
		}
	}()
	go func() {
		time.Sleep(10 * time.Second)
		buf.Close()
	}()
	for {
		fmt.Println("reading ...")
		if !buf.Next() {
			break
		}
		buf.Read(func(data []byte) {
			fmt.Println("read:", len(data))
		})
		time.Sleep(1 * time.Second)
	}
}
