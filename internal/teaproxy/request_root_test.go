package teaproxy

import (
	"bytes"
	"fmt"
	"github.com/dchest/siphash"
	"github.com/iwind/TeaGo/Tea"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestFormatURLForRoot(t *testing.T) {
	var uriString = "/www/test/index.html"
	var root = "/home/www/"

	if !filepath.IsAbs(root) {
		root = Tea.Root + Tea.DS + root
	}

	requestPath := uriString

	uri, err := url.ParseRequestURI(uriString)
	query := ""
	if err == nil {
		requestPath = uri.Path
		query = uri.RawQuery
	}

	// 去掉其中的奇怪的路径
	requestPath = strings.Replace(requestPath, "..\\", "", -1)

	t.Log(requestPath, query)
}

func BenchmarkFormatURLForRoot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var uriString = "/www/test/index.html"
		var root = "/home/www/"

		if !filepath.IsAbs(root) {
			root = Tea.Root + Tea.DS + root
		}

		requestPath := uriString

		/**uri, err := url.ParseRequestURI(uriString)
		query := ""
		if err == nil {
			requestPath = uri.Path
			query = uri.RawQuery
		}**/
		index := strings.Index(uriString, "?")
		query := ""
		if index > -1 {
			requestPath = uriString[:index]
			query = uriString[index+1:]
		}

		// 去掉其中的奇怪的路径
		requestPath = strings.Replace(requestPath, "..\\", "", -1)
		_ = query
		_ = requestPath
	}
}

func TestFileStat(t *testing.T) {
	gopath, _ := os.LookupEnv("GOPATH")
	if len(gopath) == 0 {
		return
	}

	before := time.Now()
	f, err := os.Open(gopath + "/src/github.com/TeaWeb/build/internal/teaproxy/request_root.go")
	t.Log("open", time.Since(before).Seconds()*1000, "ms")

	if err == nil {
		before := time.Now()
		buf := bytes.NewBuffer([]byte{})
		_, _ = io.Copy(buf, f)
		t.Log("copy", time.Since(before).Seconds()*1000, "ms")
	}

	t.Log(f, err)
	if err == nil {
		before := time.Now()
		t.Log(f.Stat())
		t.Log("stat", time.Since(before).Seconds()*1000, "ms")
	}
}

func BenchmarkFileStat(b *testing.B) {
	gopath, _ := os.LookupEnv("GOPATH")
	if len(gopath) == 0 {
		return
	}

	for i := 0; i < b.N; i++ {
		f, err := os.Open(gopath + "/src/github.com/TeaWeb/build/internal/teaproxy/request_root.go")
		if err == nil {
			_ = f.Close()
			_, _ = f.Stat()
		}
	}
}

func BenchmarkFileStat2(b *testing.B) {
	gopath, _ := os.LookupEnv("GOPATH")
	if len(gopath) == 0 {
		return
	}

	for i := 0; i < b.N; i++ {
		f, err := os.OpenFile(gopath+"/src/github.com/TeaWeb/build/internal/teaproxy/request_root.go", os.O_RDONLY, 0444)
		if err == nil {
			_ = f.Close()
		}

		_, _ = os.Stat(gopath + "/src/github.com/TeaWeb/build/internal/teaproxy/request_root.go")
	}
}

func TestFileEtag(t *testing.T) {
	etag := stringutil.Md5(fmt.Sprintf("%d,%d", 1563192836000, 1024))
	t.Log(etag)
}

func TestFileEtag_hash(t *testing.T) {
	etag := siphash.Hash(0, 0, []byte("123.txt"+strconv.FormatInt(1563192836000, 10)+strconv.FormatInt(1024, 10)))
	t.Log(fmt.Sprintf("%0x", etag))
}

func TestFileEtag_str(t *testing.T) {
	etag := siphash.Hash(0, 0, []byte("123.txt"+strconv.FormatInt(1563192836000, 10)+strconv.FormatInt(1024, 10)))
	t.Log(fmt.Sprintf("%0x", etag))
}

func BenchmarkFileEtag(b *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < b.N; i++ {
		_ = stringutil.Md5("123.txt" + fmt.Sprintf("%d,%d", 1563192836000, 1024))
	}
}

func BenchmarkFileEtag_Hash(b *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%0x", siphash.Hash(0, 0, []byte("123.txt"+strconv.FormatInt(1563192836000, 10)+strconv.FormatInt(1024, 10))))
	}
}
