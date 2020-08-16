package api

import "time"

// api请求数限制
type APIRequestLimit struct {
	Count    uint   `yaml:"count" json:"count"`       // 请求数 TODO
	Duration string `yaml:"duration" json:"duration"` // 请求限制间隔 TODO

	duration time.Duration
}

func (this *APIRequestLimit) Validate() error {
	return nil
}
