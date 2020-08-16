package accesslogs

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teautils/logbuffer"
	"github.com/TeaWeb/build/internal/lego/log"
	"github.com/iwind/TeaGo/logs"
	"github.com/mailru/easyjson"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestFileBuffer_Write(t *testing.T) {
	runtime.GOMAXPROCS(1)

	buf := logbuffer.NewBuffer("hello.log")
	writeAccessLogToBuffer(buf, "hello")
	writeAccessLogToBuffer(buf, "world")

	max := 10000

	go func() {
		i := 0
		for {
			i++
			if i > max {
				break
			}
			if i%1000 == 0 {
				time.Sleep(1 * time.Second)
			}
			before := time.Now()
			writeAccessLogToBuffer(buf, "Fine "+fmt.Sprintf("%d", i))
			logs.Println(time.Since(before).Seconds()*1000, "ms")
		}
	}()

	j := 0
	for {
		data, err := buf.Read()
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		j++
		if j >= max {
			break
		}
		log.Println("line:", j, len(data), string(data))
	}

	_ = os.Remove("hello.log.0.log")
}

func writeAccessLogToBuffer(buf *logbuffer.Buffer, path string) {
	accessLog := &AccessLog{
		RequestPath: path,
	}

	data, err := easyjson.Marshal(accessLog)

	if err != nil {
		logs.Error(err)
		return
	}
	_, err = buf.Write(data)
	if err != nil {
		logs.Error(err)
	}
}
