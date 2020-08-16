package teaproxy

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/jlaffaye/ftp"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestFTPConnectionPool_Get(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}
	pool := &FTPConnectionPool{
		addr:    "192.168.2.30:21",
		timeout: 5 * time.Second,
		dir:     "",
		c:       make(chan *ftp.ServerConn, 5),
	}
	wg := sync.WaitGroup{}
	count := 100
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Int()%10) * time.Second)
			client, err := pool.Get()
			if err != nil {
				t.Fatal(err)
			}
			pool.Put(client)
		}()
	}
	wg.Wait()
	t.Log(len(pool.c))
}

func TestFTPConnectionPool_Get2(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}
	pool := &FTPConnectionPool{
		addr:    "192.168.2.31:21",
		timeout: 5 * time.Second,
		dir:     "",
		c:       make(chan *ftp.ServerConn, 5),
	}
	wg := sync.WaitGroup{}
	count := 100
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Int()%10) * time.Second)
			client, err := pool.Get()
			if err != nil {
				t.Fatal(err)
			}
			pool.Put(client)
		}()
	}
	wg.Wait()
	t.Log(len(pool.c))
}

func TestChanFull_Write(t *testing.T) {
	c := make(chan int, 5)
	for i := 0; i < 5; i++ {
		c <- i
	}
	t.Log(len(c))
	select {
	case c <- 6:
		t.Log("write 6")
	default:
		t.Log("write failed")
	}
}

func TestChanFull_Read(t *testing.T) {
	c := make(chan int, 5)
	for i := 0; i < 5; i++ {
		c <- i
	}

FOR:
	for {
		select {
		case x := <-c:
			t.Log("read", x)
		default:
			t.Log("read failed")
			break FOR
		}
	}
}
