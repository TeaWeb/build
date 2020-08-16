package teaproxy

import "net/http"

// 空的响应Writer
type EmptyResponseWriter struct {
	header http.Header
}

func NewEmptyResponseWriter() *EmptyResponseWriter {
	return &EmptyResponseWriter{
		header: http.Header{},
	}
}

func (this *EmptyResponseWriter) Header() http.Header {
	return this.header
}

func (this *EmptyResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (this *EmptyResponseWriter) WriteHeader(statusCode int) {

}
