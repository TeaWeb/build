package teaproxy

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/logs"
	"net"
	"net/http"
)

// 响应Writer
type ResponseWriter struct {
	writer http.ResponseWriter

	gzipConfig *teaconfigs.GzipConfig
	gzipWriter *gzip.Writer

	statusCode    int
	sentBodyBytes int64

	bodyCopying    bool
	body           []byte
	gzipBodyBuffer *bytes.Buffer // 当使用gzip压缩时使用
	gzipBodyWriter *gzip.Writer  // 当使用gzip压缩时使用
}

// 包装对象
func NewResponseWriter(httpResponseWriter http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		writer: httpResponseWriter,
	}
}

// 重置
func (this *ResponseWriter) Reset(httpResponseWriter http.ResponseWriter) {
	this.writer = httpResponseWriter

	this.gzipConfig = nil
	this.gzipWriter = nil

	this.statusCode = 0
	this.sentBodyBytes = 0

	this.bodyCopying = false
	this.body = nil
	this.gzipBodyBuffer = nil
	this.gzipBodyWriter = nil
}

// 设置Gzip
func (this *ResponseWriter) Gzip(config *teaconfigs.GzipConfig) {
	this.gzipConfig = config
}

// 准备输出
func (this *ResponseWriter) Prepare(size int64) {
	if this.gzipConfig == nil || this.gzipConfig.Level <= 0 {
		return
	}

	// 尺寸和类型
	if size < this.gzipConfig.MinBytes() {
		return
	}

	contentType := this.Header().Get("Content-Type")
	if !this.gzipConfig.MatchContentType(contentType) {
		return
	}

	// 如果已经有编码则不处理
	if len(this.writer.Header().Get("Content-Encoding")) > 0 {
		return
	}

	// gzip writer
	var err error = nil
	this.gzipWriter, err = gzip.NewWriterLevel(this.writer, int(this.gzipConfig.Level))
	if err != nil {
		logs.Error(err)
		return
	}

	// body copy
	if this.bodyCopying {
		this.gzipBodyBuffer = bytes.NewBuffer([]byte{})
		this.gzipBodyWriter, err = gzip.NewWriterLevel(this.gzipBodyBuffer, int(this.gzipConfig.Level))
		if err != nil {
			logs.Error(err)
		}
	}

	header := this.writer.Header()
	header.Set("Content-Encoding", "gzip")
	header.Set("Transfer-Encoding", "chunked")
	header.Set("Vary", "Accept-Encoding")
	header.Del("Content-Length")
}

// 包装前的原始的Writer
func (this *ResponseWriter) Raw() http.ResponseWriter {
	return this.writer
}

// 获取Header
func (this *ResponseWriter) Header() http.Header {
	if this.writer == nil {
		return http.Header{}
	}
	return this.writer.Header()
}

// 添加一组Header
func (this *ResponseWriter) AddHeaders(header http.Header) {
	if this.writer == nil {
		return
	}
	for key, value := range header {
		for _, v := range value {
			this.writer.Header().Add(key, v)
		}
	}
}

// 写入数据
func (this *ResponseWriter) Write(data []byte) (n int, err error) {
	if this.writer != nil {
		if this.gzipWriter != nil {
			n, err = this.gzipWriter.Write(data)
		} else {
			n, err = this.writer.Write(data)
		}
		if n > 0 {
			this.sentBodyBytes += int64(n)
		}
	} else {
		if n == 0 {
			n = len(data) // 防止出现short write错误
		}
	}
	if this.bodyCopying {
		if this.gzipBodyWriter != nil {
			_, err := this.gzipBodyWriter.Write(data)
			if err != nil {
				logs.Error(err)
			}
		} else {
			this.body = append(this.body, data...)
		}
	}
	return
}

// 写入字符串
func (this *ResponseWriter) WriteString(s string) (n int, err error) {
	return this.Write([]byte(s))
}

// 读取发送的字节数
func (this *ResponseWriter) SentBodyBytes() int64 {
	return this.sentBodyBytes
}

// 写入状态码
func (this *ResponseWriter) WriteHeader(statusCode int) {
	if this.writer != nil {
		this.writer.WriteHeader(statusCode)
	}
	this.statusCode = statusCode
}

// 读取状态码
func (this *ResponseWriter) StatusCode() int {
	if this.statusCode == 0 {
		return http.StatusOK
	}
	return this.statusCode
}

// 设置拷贝Body数据
func (this *ResponseWriter) SetBodyCopying(b bool) {
	this.bodyCopying = b
}

// 判断是否在拷贝Body数据
func (this *ResponseWriter) BodyIsCopying() bool {
	return this.bodyCopying
}

// 读取拷贝的Body数据
func (this *ResponseWriter) Body() []byte {
	return this.body
}

// 读取Header二进制数据
func (this *ResponseWriter) HeaderData() []byte {
	if this.writer == nil {
		return nil
	}

	resp := &http.Response{}
	resp.Header = this.Header()
	if this.statusCode == 0 {
		this.statusCode = http.StatusOK
	}
	resp.StatusCode = this.statusCode
	resp.ProtoMajor = 1
	resp.ProtoMinor = 1

	resp.ContentLength = 1 // Trick：这样可以屏蔽Content-Length

	writer := bytes.NewBuffer([]byte{})
	_ = resp.Write(writer)
	return writer.Bytes()
}

// 关闭
func (this *ResponseWriter) Close() {
	if this.gzipWriter != nil {
		if this.bodyCopying && this.gzipBodyWriter != nil {
			_ = this.gzipBodyWriter.Close()
			this.body = this.gzipBodyBuffer.Bytes()
		}
		_ = this.gzipWriter.Close()
		this.gzipWriter = nil
	}
}

// Hijack
func (this *ResponseWriter) Hijack() (conn net.Conn, buf *bufio.ReadWriter, err error) {
	hijack, ok := this.writer.(http.Hijacker)
	if ok {
		return hijack.Hijack()
	}
	return
}

// Flush
func (this *ResponseWriter) Flush() {
	flusher, ok := this.writer.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
