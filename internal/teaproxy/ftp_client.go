package teaproxy

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"mime"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"
)

// FTP客户端
type FTPClient struct {
	backend *teaconfigs.BackendConfig
	pool    *FTPConnectionPool
}

// 执行请求
func (this *FTPClient) Do(req *http.Request) (*http.Response, error) {
	return this.doRetries(req, 0)
}

// 执行请求
func (this *FTPClient) doRetries(req *http.Request, retries int) (*http.Response, error) {
	conn, err := this.pool.Get()
	if err != nil {
		// retry
		if err == ErrFTPTooManyConnections && retries == 0 {
			time.Sleep(1 * time.Second)
			return this.doRetries(req, retries+1)
		}
		return nil, err
	}

	// file size
	path := strings.TrimLeft(req.URL.Path, "/")
	size, err := conn.FileSize(path)
	if err != nil {
		isDisconnected := false
		textErr, ok := err.(*textproto.Error)
		if ok {
			if textErr.Code != ftp.StatusFileUnavailable &&
				textErr.Code != ftp.StatusBadFileName {
				isDisconnected = true
			} else {
				// not found
				this.pool.Put(conn)
				message := "404 page not found"
				return &http.Response{
					Status:        http.StatusText(http.StatusNotFound),
					StatusCode:    http.StatusNotFound,
					Body:          ioutil.NopCloser(bytes.NewReader([]byte(message))),
					ContentLength: int64(len(message)),
					Close:         false,
				}, nil
			}
		} else {
			isDisconnected = true
		}

		// retry
		if isDisconnected {
			this.pool.Decrease()
			_ = this.pool.CloseAll()

			if retries == 0 {
				return this.doRetries(req, retries+1)
			}
		} else {
			this.pool.Put(conn)
		}
		return nil, err
	}

	// read file
	response, err := conn.Retr(path)
	if err != nil {
		if response != nil {
			_ = response.Close()
		}
		return nil, err
	}

	ext := filepath.Ext(path)
	headers := map[string][]string{}
	if len(ext) > 0 {
		mimeType := mime.TypeByExtension(ext)
		if len(mimeType) > 0 {
			headers["Content-Type"] = []string{mimeType}
		}
	}

	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Header:     headers,
		Body: &FTPResponseBody{
			ReadCloser: response,
			callback: func() {
				this.pool.Put(conn)
			},
		},
		ContentLength: size,
		Close:         false,
	}
	return resp, nil
}

// 关闭客户端
func (this *FTPClient) Close() error {
	if this.pool == nil {
		return nil
	}
	return this.pool.CloseAll()
}
