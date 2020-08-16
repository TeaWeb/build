package logbuffer

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/iwind/TeaGo/Tea"
	"io"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestFileBuffer(t *testing.T) {
	buffer := NewBuffer(Tea.Root + "/logs/logbuffer")
	buffer.chunkSize = 20

	// write
	for i := 0; i < 20; i++ {
		_, err := buffer.Write([]byte("Hello " + fmt.Sprintf("%d", i)))
		if err != nil {
			t.Fatal(err)
		}
	}

	for index, f := range buffer.files {
		t.Log(index, f.Size())
	}

	// read
	for i := 0; i < 100; i++ {
		data, err := buffer.Read()
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			continue
		}
		t.Log(string(data))
	}
}

func TestBuffer_Init(t *testing.T) {
	buffer := NewBuffer(Tea.Root + "/logs/logbuffer")
	buffer.chunkSize = 20
}

func TestBuffer_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(1)

	buffer := NewBuffer(Tea.Root + "/logs/logbuffer")
	buffer.chunkSize = 1 << 20

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		// write
		for i := 0; i < 100; i++ {
			_, err := buffer.Write([]byte("Hello " + fmt.Sprintf("%d", i) + " " + strings.Repeat("HELLO", 1000)))
			if err != nil {
				t.Fatal(err)
			}
			time.Sleep(100000 * time.Nanosecond)
		}
	}()

	go func() {
		defer wg.Done()

		// read
		count := 0
		for i := 0; i < 1000; i++ {
			before := time.Now()
			data, err := buffer.Read()
			cost := time.Since(before).Seconds() * 1000
			if err != nil {
				t.Fatal(err)
			}
			if len(data) == 0 {
				time.Sleep(1000000 * time.Nanosecond)
				continue
			}
			if !bytes.HasPrefix(data, []byte("Hello")) {
				t.Fatal("invalid line")
			}
			count++
			t.Log(count, ":", cost, "ms", len(data), "bytes", string(data[:10]))
		}
	}()

	wg.Wait()
	t.Log("ok")
}

func TestBuffer_Read_Full(t *testing.T) {
	reader := bytes.NewBuffer([]byte(`1|
0123456789|
0123456789_0123456789_0123456789_0123456789_0123456789_0123456789|
0123456789_0123456789_0123456789_0123456789_0123456789_0123456789|
0123456789_0123456789_0123456789_0123456789_0123456789_0123456789|0123456789_0123456789_0123456789_0123456789_0123456789_0123456789|0123456789_0123456789_0123456789_0123456789_0123456789_0123456789|0123456789_0123456789_0123456789_0123456789_0123456789_0123456789|0123456789_0123456789_0123456789_0123456789_0123456789_0123456789|
0123456789_0123456789_0123456789_0123456789_0123456789_0123456789|
0123456789|
`))
	buf := bufio.NewReaderSize(reader, 20)
	for {
		var line []byte
		for {
			line2, err := buf.ReadSlice('\n')
			if err == nil {
				if len(line) == 0 {
					line = line2
				} else {
					line = append(line, line2...)
				}
				break
			}
			if err == bufio.ErrBufferFull {
				line = append(line, line2...)
				continue
			}
			if err != nil {
				if err == io.EOF {
					return
				}
				t.Fatal(err)
			}
		}
		t.Log(string(line[:len(line)-1]))
	}
}
