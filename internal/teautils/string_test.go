package teautils

import (
	"strings"
	"testing"
)

func TestBytesToString(t *testing.T) {
	t.Log(UnsafeBytesToString([]byte("Hello,World")))
}

func TestStringToBytes(t *testing.T) {
	t.Log(string(UnsafeStringToBytes("Hello,World")))
}

func BenchmarkBytesToString(b *testing.B) {
	data := []byte("Hello,World")
	for i := 0; i < b.N; i++ {
		_ = UnsafeBytesToString(data)
	}
}

func BenchmarkBytesToString2(b *testing.B) {
	data := []byte("Hello,World")
	for i := 0; i < b.N; i++ {
		_ = string(data)
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	s := strings.Repeat("Hello,World", 1024)
	for i := 0; i < b.N; i++ {
		_ = UnsafeStringToBytes(s)
	}
}

func BenchmarkStringToBytes2(b *testing.B) {
	s := strings.Repeat("Hello,World", 1024)
	for i := 0; i < b.N; i++ {
		_ = []byte(s)
	}
}

func TestFormatAddress(t *testing.T) {
	t.Log(FormatAddress("127.0.0.1:1234"))
	t.Log(FormatAddress("127.0.0.1 : 1234"))
	t.Log(FormatAddress("127.0.0.1：1234"))
}

func TestFormatAddressList(t *testing.T) {
	t.Log(FormatAddressList([]string{
		"127.0.0.1:1234",
		"127.0.0.1 : 1234",
		"127.0.0.1：1234",
	}))
}
