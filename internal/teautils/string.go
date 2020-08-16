package teautils

import (
	"strings"
	"unsafe"
)

// convert bytes to string
func UnsafeBytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// convert string to bytes
func UnsafeStringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// format address
func FormatAddress(addr string) string {
	if strings.HasSuffix(addr, "unix:") {
		return addr
	}
	addr = strings.Replace(addr, " ", "", -1)
	addr = strings.Replace(addr, "\t", "", -1)
	addr = strings.Replace(addr, "：", ":", -1)
	addr = strings.TrimSpace(addr)
	return addr
}

// format address list
func FormatAddressList(addrList []string) []string {
	result := []string{}
	for _, addr := range addrList {
		result = append(result, FormatAddress(addr))
	}
	return result
}
