package checkpoints

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"github.com/iwind/TeaGo/lists"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// ${requestUpload.arg}
type RequestUploadCheckpoint struct {
	Checkpoint
}

func (this *RequestUploadCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = ""
	if param == "minSize" || param == "maxSize" {
		value = 0
	}

	if req.Method != http.MethodPost {
		return
	}

	if req.Body == nil {
		return
	}

	if req.MultipartForm == nil {
		if len(req.BodyData) == 0 {
			data, err := req.ReadBody(32 * 1024 * 1024)
			if err != nil {
				sysErr = err
				return
			}

			req.BodyData = data
			defer req.RestoreBody(data)
		}
		oldBody := req.Body
		req.Body = ioutil.NopCloser(bytes.NewBuffer(req.BodyData))

		err := req.ParseMultipartForm(32 * 1024 * 1024)

		// 还原
		req.Body = oldBody

		if err != nil {
			userErr = err
			return
		}

		if req.MultipartForm == nil {
			return
		}
	}

	if param == "field" { // field
		fields := []string{}
		for field := range req.MultipartForm.File {
			fields = append(fields, field)
		}
		value = strings.Join(fields, ",")
	} else if param == "minSize" { // minSize
		minSize := int64(0)
		for _, files := range req.MultipartForm.File {
			for _, file := range files {
				if minSize == 0 || minSize > file.Size {
					minSize = file.Size
				}
			}
		}
		value = minSize
	} else if param == "maxSize" { // maxSize
		maxSize := int64(0)
		for _, files := range req.MultipartForm.File {
			for _, file := range files {
				if maxSize < file.Size {
					maxSize = file.Size
				}
			}
		}
		value = maxSize
	} else if param == "name" { // name
		names := []string{}
		for _, files := range req.MultipartForm.File {
			for _, file := range files {
				if !lists.ContainsString(names, file.Filename) {
					names = append(names, file.Filename)
				}
			}
		}
		value = strings.Join(names, ",")
	} else if param == "ext" { // ext
		extensions := []string{}
		for _, files := range req.MultipartForm.File {
			for _, file := range files {
				if len(file.Filename) > 0 {
					exit := strings.ToLower(filepath.Ext(file.Filename))
					if !lists.ContainsString(extensions, exit) {
						extensions = append(extensions, exit)
					}
				}
			}
		}
		value = strings.Join(extensions, ",")
	}

	return
}

func (this *RequestUploadCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}

func (this *RequestUploadCheckpoint) ParamOptions() *ParamOptions {
	option := NewParamOptions()
	option.AddParam("最小文件尺寸", "minSize")
	option.AddParam("最大文件尺寸", "maxSize")
	option.AddParam("扩展名(如.txt)", "ext")
	option.AddParam("原始文件名", "name")
	option.AddParam("表单字段名", "field")
	return option
}
