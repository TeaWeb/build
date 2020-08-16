package checkpoints

import (
	"github.com/TeaWeb/build/internal/teamemory"
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net"
	"regexp"
	"strings"
	"sync"
)

// ${cc.arg}
// TODO implement more traffic rules
type CCCheckpoint struct {
	Checkpoint

	grid *teamemory.Grid
	once sync.Once
}

func (this *CCCheckpoint) Init() {

}

func (this *CCCheckpoint) Start() {
	if this.grid != nil {
		this.grid.Destroy()
	}
	this.grid = teamemory.NewGrid(32, teamemory.NewLimitCountOpt(1000_0000))
}

func (this *CCCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = 0

	if this.grid == nil {
		this.once.Do(func() {
			this.Start()
		})
		if this.grid == nil {
			return
		}
	}

	periodString, ok := options["period"]
	if !ok {
		return
	}
	period := types.Int64(periodString)
	if period < 1 {
		return
	}

	userType, _ := options["userType"]
	userField, _ := options["userField"]
	userIndex, _ := options["userIndex"]
	userIndexInt := types.Int(userIndex)

	if param == "requests" { // requests
		var key = ""
		switch userType {
		case "ip":
			key = this.ip(req)
		case "cookie":
			if len(userField) == 0 {
				key = this.ip(req)
			} else {
				cookie, _ := req.Cookie(userField)
				if cookie != nil {
					v := cookie.Value
					if userIndexInt > 0 && len(v) > userIndexInt {
						v = v[userIndexInt:]
					}
					key = "USER@" + userType + "@" + userField + "@" + v
				}
			}
		case "get":
			if len(userField) == 0 {
				key = this.ip(req)
			} else {
				v := req.URL.Query().Get(userField)
				if userIndexInt > 0 && len(v) > userIndexInt {
					v = v[userIndexInt:]
				}
				key = "USER@" + userType + "@" + userField + "@" + v
			}
		case "post":
			if len(userField) == 0 {
				key = this.ip(req)
			} else {
				v := req.PostFormValue(userField)
				if userIndexInt > 0 && len(v) > userIndexInt {
					v = v[userIndexInt:]
				}
				key = "USER@" + userType + "@" + userField + "@" + v
			}
		case "header":
			if len(userField) == 0 {
				key = this.ip(req)
			} else {
				v := req.Header.Get(userField)
				if userIndexInt > 0 && len(v) > userIndexInt {
					v = v[userIndexInt:]
				}
				key = "USER@" + userType + "@" + userField + "@" + v
			}
		default:
			key = this.ip(req)
		}
		if len(key) == 0 {
			key = this.ip(req)
		}
		value = this.grid.IncreaseInt64([]byte(key), 1, period)
	}

	return
}

func (this *CCCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}

func (this *CCCheckpoint) ParamOptions() *ParamOptions {
	option := NewParamOptions()
	option.AddParam("请求数", "requests")
	return option
}

func (this *CCCheckpoint) Options() []OptionInterface {
	options := []OptionInterface{}

	// period
	{
		option := NewFieldOption("统计周期", "period")
		option.Value = "60"
		option.RightLabel = "秒"
		option.Size = 8
		option.MaxLength = 8
		option.Validate = func(value string) (ok bool, message string) {
			if regexp.MustCompile("^\\d+$").MatchString(value) {
				ok = true
				return
			}
			message = "周期需要是一个整数数字"
			return
		}
		options = append(options, option)
	}

	// type
	{
		option := NewOptionsOption("用户识别读取来源", "userType")
		option.Size = 10
		option.SetOptions([]maps.Map{
			{
				"name":  "IP",
				"value": "ip",
			},
			{
				"name":  "Cookie",
				"value": "cookie",
			},
			{
				"name":  "URL参数",
				"value": "get",
			},
			{
				"name":  "POST参数",
				"value": "post",
			},
			{
				"name":  "HTTP Header",
				"value": "header",
			},
		})
		options = append(options, option)
	}

	// user field
	{
		option := NewFieldOption("用户识别字段", "userField")
		option.Comment = "识别用户的唯一性字段，在用户读取来源不是IP时使用"
		options = append(options, option)
	}

	// user value index
	{
		option := NewFieldOption("字段读取位置", "userIndex")
		option.Size = 5
		option.MaxLength = 5
		option.Comment = "读取用户识别字段的位置，从0开始，比如user12345的数字ID 12345的位置就是5，在用户读取来源不是IP时使用"
		options = append(options, option)
	}

	return options
}

func (this *CCCheckpoint) Stop() {
	if this.grid != nil {
		this.grid.Destroy()
		this.grid = nil
	}
}

func (this *CCCheckpoint) ip(req *requests.Request) string {
	// X-Forwarded-For
	forwardedFor := req.Header.Get("X-Forwarded-For")
	if len(forwardedFor) > 0 {
		commaIndex := strings.Index(forwardedFor, ",")
		if commaIndex > 0 {
			return forwardedFor[:commaIndex]
		}
		return forwardedFor
	}

	// Real-IP
	{
		realIP, ok := req.Header["X-Real-IP"]
		if ok && len(realIP) > 0 {
			return realIP[0]
		}
	}

	// Real-Ip
	{
		realIP, ok := req.Header["X-Real-Ip"]
		if ok && len(realIP) > 0 {
			return realIP[0]
		}
	}

	// Remote-Addr
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		return host
	}
	return req.RemoteAddr
}
