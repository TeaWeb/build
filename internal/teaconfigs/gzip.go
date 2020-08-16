package teaconfigs

import (
	stringutil "github.com/iwind/TeaGo/utils/string"
	"regexp"
	"strings"
)

// 默认的文件类型
var (
	DefaultGzipMimeTypes = []string{"text/html", "application/json"}
)

// gzip配置
type GzipConfig struct {
	Level     int8     `yaml:"level" json:"level"`         // 1-9
	MinLength string   `yaml:"minLength" json:"minLength"` // 比如4m, 24k
	MimeTypes []string `yaml:"mimeTypes" json:"mimeTypes"` // 比如text/html, text/*

	minLength int64
	mimeTypes []*MimeTypeRule
}

// 校验
func (this *GzipConfig) Validate() error {
	gzipMinLength, _ := stringutil.ParseFileSize(this.MinLength)
	this.minLength = int64(gzipMinLength)
	if len(this.MimeTypes) == 0 {
		this.MimeTypes = DefaultGzipMimeTypes
	}

	this.mimeTypes = []*MimeTypeRule{}
	for _, mimeType := range this.MimeTypes {
		if strings.Contains(mimeType, "*") {
			mimeType = regexp.QuoteMeta(mimeType)
			mimeType = strings.Replace(mimeType, "\\*", ".*", -1)
			reg, err := regexp.Compile("^" + mimeType + "$")
			if err != nil {
				return err
			}
			this.mimeTypes = append(this.mimeTypes, &MimeTypeRule{
				Value:  mimeType,
				Regexp: reg,
			})
		} else {
			this.mimeTypes = append(this.mimeTypes, &MimeTypeRule{
				Value:  mimeType,
				Regexp: nil,
			})
		}
	}
	return nil
}

// 可压缩最小尺寸
func (this *GzipConfig) MinBytes() int64 {
	return this.minLength
}

// 检查是否匹配Content-Type
func (this *GzipConfig) MatchContentType(contentType string) bool {
	index := strings.Index(contentType, ";")
	if index >= 0 {
		contentType = contentType[:index]
	}
	for _, mimeType := range this.mimeTypes {
		if mimeType.Regexp == nil && contentType == mimeType.Value {
			return true
		} else if mimeType.Regexp != nil && mimeType.Regexp.MatchString(contentType) {
			return true
		}
	}
	return false
}
