// +build !go1.12

package teautils

import "net/http"

// 关闭客户端连接
func CloseHTTPClient(client *http.Client) {
	// do nothing
}
