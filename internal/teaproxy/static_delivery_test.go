package teaproxy

import (
	"os"
	"testing"
	"time"
)

func TestStaticDelivery_Read(t *testing.T) {
	delivery := NewStaticDelivery()

	root, _ := os.LookupEnv("GOPATH")
	stat, err := os.Stat(root + "/src/github.com/TeaWeb/build/internal/teaproxy/static_delivery.go")
	if err != nil {
		t.Fatal(err)
	}

	{
		before := time.Now()
		reader, shouldClose, err := delivery.Read(root+"/src/github.com/TeaWeb/build/internal/teaproxy/static_delivery.go", stat)
		if err != nil {
			t.Fatal(err)
		}
		if reader == nil {
			t.Fatal("reader is nil")
		}

		t.Log("shouldClose:", shouldClose)
		t.Log(1 / time.Since(before).Seconds())
	}

	for i := 0; i < 10; i ++ {
		before := time.Now()
		reader, _, _ := delivery.Read(root+"/src/github.com/TeaWeb/build/internal/teaproxy/static_delivery.go", stat)
		if reader == nil {
			t.Fatal("reader is nil")
		}
		t.Log(1 / time.Since(before).Seconds())
	}
}
