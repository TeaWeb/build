package teautils

import (
	"os"
	"testing"
)

func TestRlimit(t *testing.T) {
	err := SetRLimit(20480)
	gopath, _ := os.LookupEnv("GOPATH")
	if len(gopath) == 0 {
		return
	}
	if err == nil {
		for i := 0; i < 10240; i++ {
			_, err := os.Open(gopath + "/src/github.com/TeaWeb/build/internal/teautils/rlimit_test.go")
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Log("OK")
	}
}

func TestSetSuitableRLimit(t *testing.T) {
	SetSuitableRLimit()
}
