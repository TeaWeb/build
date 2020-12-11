package accesslogs

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/TeaWeb/build/internal/teageo"
	"github.com/TeaWeb/build/internal/teamemory"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/uaparser"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/mailru/easyjson"
	"github.com/pquerna/ffjson/ffjson"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var userAgentGrid = teamemory.NewGrid(32, teamemory.NewLimitCountOpt(100_0000))
var charsetReg = regexp.MustCompile("(?i)charset\\s*=\\s*([\\w-]+)")
var headerReg = regexp.MustCompile("([A-Z])")

type AccessLog struct {
	Id shared.ObjectId `var:"id" bson:"_id" json:"id"` // 数据库存储的ID

	ServerId   string `var:"serverId" bson:"serverId" json:"serverId"`       // 服务ID
	BackendId  string `var:"backendId" bson:"backendId" json:"backendId"`    // 后端服务ID
	LocationId string `var:"locationId" bson:"locationId" json:"locationId"` // 路径配置ID
	FastcgiId  string `var:"fastcgiId" bson:"fastcgiId" json:"fastcgiId"`    // Fastcgi配置ID
	RewriteId  string `var:"rewriteId" bson:"rewriteId" json:"rewriteId"`    // 重写规则ID

	TeaVersion      string  `var:"teaVersion" bson:"teaVersion" json:"teaVersion"`                // TeaWeb版本
	RemoteAddr      string  `var:"remoteAddr" bson:"remoteAddr" json:"remoteAddr"`                // 终端地址
	RawRemoteAddr   string  `var:"rawRemoteAddr" bson:"rawRemoteAddr" json:"rawRemoteAddr"`       // 直接连接服务器的终端地址
	RemotePort      int     `var:"remotePort" bson:"remotePort" json:"remotePort"`                // 终端端口
	RemoteUser      string  `var:"remoteUser" bson:"remoteUser" json:"remoteUser"`                // 终端用户，基于BasicAuth认证
	RequestURI      string  `var:"requestURI" bson:"requestURI" json:"requestURI"`                // 请求URI
	RequestPath     string  `var:"requestPath" bson:"requestPath" json:"requestPath"`             // 请求URI中的路径
	APIPath         string  `var:"apiPath" bson:"apiPath" json:"apiPath"`                         // API路径
	APIStatus       string  `var:"apiStatus" bson:"apiStatus" json:"apiStatus"`                   // API状态码
	RequestLength   int64   `var:"requestLength" bson:"requestLength" json:"requestLength"`       // 请求内容长度
	RequestTime     float64 `var:"requestTime" bson:"requestTime" json:"requestTime"`             // 从请求到所有响应数据发送到请求端所花时间，单位为带有小数点的秒，精确到纳秒，比如：0.000260081
	RequestMethod   string  `var:"requestMethod" bson:"requestMethod" json:"requestMethod"`       // 请求方法
	RequestFilename string  `var:"requestFilename" bson:"requestFilename" json:"requestFilename"` // 请求的文件名，包含完整的路径
	Scheme          string  `var:"scheme" bson:"scheme" json:"scheme"`                            // 请求协议，http或者https
	Proto           string  `var:"proto" bson:"proto" json:"proto"`                               // 请求协议，比如HTTP/1.0, HTTP/1.1

	BytesSent     int64               `var:"bytesSent" bson:"bytesSent" json:"bytesSent"`             // 响应的字节数，目前同 bodyBytesSent
	BodyBytesSent int64               `var:"bodyBytesSent" bson:"bodyBytesSent" json:"bodyBytesSent"` // 响应的字节数
	Status        int                 `var:"status" bson:"status" json:"status"`                      // 响应的状态码
	StatusMessage string              `var:"statusMessage" bson:"statusMessage" json:"statusMessage"` // 响应的信息
	SentHeader    map[string][]string `var:"sentHeader" bson:"sentHeader" json:"sentHeader"`          // 响应的头信息

	TimeISO8601    string              `var:"timeISO8601" bson:"timeISO8601" json:"timeISO8601"`          // ISO 8601格式的本地时间，比如 2018-07-16T23:52:24.839+08:00
	TimeLocal      string              `var:"timeLocal" bson:"timeLocal" json:"timeLocal"`                // 本地时间，比如 17/Jul/2018:09:52:24 +0800
	Msec           float64             `var:"msec" bson:"msec" json:"msec"`                               // 带有毫秒的时间，比如 1531756823.054
	Timestamp      int64               `var:"timestamp" bson:"timestamp" json:"timestamp"`                // unix时间戳，单位为秒
	Host           string              `var:"host" bson:"host" json:"host"`                               // 主机名
	Referer        string              `var:"referer" bson:"referer" json:"referer"`                      // 请求来源URL
	UserAgent      string              `var:"userAgent" bson:"userAgent" json:"userAgent"`                // 客户端信息
	Request        string              `var:"request" bson:"request" json:"request"`                      // 请求的简要说明，格式类似于 GET /hello/world HTTP/1.1
	ContentType    string              `var:"contentType" bson:"contentType" json:"contentType"`          // 请求头部的Content-Type
	Cookie         map[string]string   `bson:"cookie" json:"cookie"`                                      // Cookie cookie.name, cookie.sid
	Arg            map[string][]string `bson:"arg" json:"arg"`                                            // arg_name, arg_id
	Args           string              `var:"args" bson:"args" json:"args"`                               // name=liu&age=20
	QueryString    string              `var:"queryString" bson:"queryString" json:"queryString"`          // 同 Args
	Header         map[string][]string `bson:"header" json:"header"`                                      // 请求的头部信息，支持header_*和http_*，header_content_type, header_expires, http_content_type, http_user_agent
	ServerName     string              `var:"serverName" bson:"serverName" json:"serverName"`             // 接收请求的服务器名
	ServerPort     int                 `var:"serverPort" bson:"serverPort" json:"serverPort"`             // 服务器端口
	ServerProtocol string              `var:"serverProtocol" bson:"serverProtocol" json:"serverProtocol"` // 服务器协议，类似于HTTP/1.0”
	Hostname       string              `var:"hostname" bson:"hostname" json:"hostname"`

	// 代理相关
	BackendAddress string `var:"backendAddress" bson:"backendAddress" json:"backendAddress"` // 代理的后端的地址
	BackendCode    string `var:"backendCode" bson:"backendCode" json:"backendCode"`          // 后端代号
	BackendScheme  string `var:"backendScheme" bson:"backendScheme" json:"backendScheme"`    // 后端协议
	FastcgiAddress string `var:"fastcgiAddress" bson:"fastcgiAddress" json:"fastcgiAddress"` // Fastcgi后端地址

	// 调试用
	RequestData        []byte `var:"" bson:"requestData" json:"requestData"`               // 请求数据
	ResponseHeaderData []byte `var:"" bson:"responseHeaderData" json:"responseHeaderData"` // 响应Header数据
	ResponseBodyData   []byte `var:"" bson:"responseBodyData" json:"responseBodyData"`     // 响应Body数据

	// 错误信息
	Errors    []string `var:"errors" bson:"errors" json:"errors"`          // 错误信息
	HasErrors bool     `var:"hasErrors" bson:"hasErrors" json:"hasErrors"` // 是否包含有错误信息

	// 扩展
	Extend *AccessLogExtend  `bson:"extend" json:"extend"`
	Attrs  map[string]string `bson:"attrs" json:"attrs"`

	// 日志
	StorageOnly      bool     `bson:"storageOnly" json:"storageOnly"`           // 是否只存储到日志策略中
	StoragePolicyIds []string `bson:"storagePolicyIds" json:"storagePolicyIds"` // 日志策略Ids

	shouldStat    bool  // 是否应该统计
	shouldWrite   bool  // 是否写入
	writingFields []int // 写入的字段
}

type AccessLogExtend struct {
	File   AccessLogFile   `bson:"file" json:"file"`
	Client AccessLogClient `bson:"client" json:"client"`
	Geo    AccessLogGeo    `bson:"geo" json:"geo"`
}

type AccessLogFile struct {
	MimeType  string `bson:"mimeType" json:"mimeType"`   // 类似于 image/jpeg
	Extension string `bson:"extension" json:"extension"` // 扩展名，不带点（.）
	Charset   string `bson:"charset" json:"charset"`     // 字符集，统一大写
}

type AccessLogClient struct {
	OS      AccessLogClientOS      `bson:"os" json:"os"`
	Device  AccessLogClientDevice  `bson:"device" json:"device"`
	Browser AccessLogClientBrowser `bson:"browser" json:"browser"`
}

type AccessLogClientOS struct {
	Family     string `bson:"family" json:"family"`
	Major      string `bson:"major" json:"major"`
	Minor      string `bson:"minor" json:"minor"`
	Patch      string `bson:"patch" json:"patch"`
	PatchMinor string `bson:"patchMinor" json:"patchMinor"`
}

type AccessLogClientDevice struct {
	Family string `bson:"family" json:"family"`
	Brand  string `bson:"brand" json:"brand"`
	Model  string `bson:"model" json:"model"`
}

type AccessLogClientBrowser struct {
	Family string `bson:"family" json:"family"`
	Major  string `bson:"major" json:"major"`
	Minor  string `bson:"minor" json:"minor"`
	Patch  string `bson:"patch" json:"patch"`
}

type AccessLogGeo struct {
	Region   string               `bson:"region" json:"region"`
	State    string               `bson:"state" json:"state"`
	City     string               `bson:"city" json:"city"`
	Location AccessLogGeoLocation `bson:"location" json:"location"`
}

type AccessLogGeoLocation struct {
	Latitude       float64 `bson:"latitude" json:"latitude"`
	Longitude      float64 `bson:"longitude" json:"longitude"`
	TimeZone       string  `bson:"timeZone" json:"timeZone"`
	AccuracyRadius uint16  `bson:"accuracyRadius" json:"accuracyRadius"`
	MetroCode      uint    `bson:"metroCode" json:"metroCode"`
}

// 获取新对象
func NewAccessLog() *AccessLog {
	return &AccessLog{}
}

// 获取访问日志的请求时间
func (this *AccessLog) Time() time.Time {
	return time.Unix(this.Timestamp, 0)
}

func (this *AccessLog) SentContentType() string {
	contentType, ok := this.SentHeader["Content-Type"]
	if ok && len(contentType) > 0 {
		return contentType[0]
	}
	return ""
}

func (this *AccessLog) Format(format string) string {
	refValue := reflect.ValueOf(*this)

	// 处理变量${varName}
	format = teautils.ParseVariables(format, func(varName string) (value string) {
		fieldName, found := accessLogVars[varName]
		if found {
			field := refValue.FieldByName(fieldName)
			if field.IsValid() {
				kind := field.Kind()
				switch kind {
				case reflect.String:
					return field.String()
				case reflect.Float32, reflect.Float64:
					return fmt.Sprintf("%f", field.Interface())
				default:
					return types.String(field.Interface())
				}
			}

			return ""
		}

		// json log
		if varName == "log" {
			data, _ := easyjson.Marshal(this)
			return string(data)
		}

		// arg
		if strings.HasPrefix(varName, "arg.") {
			varName = varName[4:]
			values, found := this.Arg[varName]
			if found {
				countValues := len(values)
				if countValues == 1 {
					return values[0]
				} else if countValues > 1 {
					return "[" + strings.Join(values, ",") + "]"
				}
			}
			return ""
		}

		// cookie
		if strings.HasPrefix(varName, "cookie.") {
			varName = varName[7:]
			value, found := this.Cookie[varName]
			if found {
				return value
			}
			return ""
		}

		// http
		if strings.HasPrefix(varName, "http.") {
			varName = varName[5:]
			values, found := this.Header[varName]
			if found {
				if len(values) > 0 {
					return values[0]
				}
			} else {
				varName = strings.TrimPrefix(headerReg.ReplaceAllString(varName, "-${1}"), "-")
				values, found := this.Header[varName]
				if found && len(values) > 0 {
					return values[0]
				}
			}

			return ""
		}

		// header
		if strings.HasPrefix(varName, "header.") {
			varName = varName[7:]
			values, found := this.Header[varName]
			if found {
				if len(values) > 0 {
					return values[0]
				}
			} else {
				varName = strings.TrimPrefix(headerReg.ReplaceAllString(varName, "-${1}"), "-")
				values, found := this.Header[varName]
				if found && len(values) > 0 {
					return values[0]
				}
			}

			return ""
		}

		// extend
		if strings.HasPrefix(varName, "extend.") {
			value := teautils.Get(this.Extend, strings.Split(varName[7:], "."))
			jsonValue, err := ffjson.Marshal(value)
			if err != nil {
				logs.Error(err)
			} else {
				return string(jsonValue)
			}
		}

		// backend
		if strings.HasPrefix(varName, "backend.") {
			subVarName := varName[8:]
			switch subVarName {
			case "id":
				return this.BackendId
			case "address":
				return this.BackendAddress
			case "host":
				index := strings.Index(this.BackendAddress, ":")
				if index > -1 {
					return this.BackendAddress[:index]
				} else {
					return ""
				}
			case "scheme":
				return this.BackendScheme
			case "code":
				return this.BackendCode
			}
		}

		return "${" + varName + "}"
	})

	return format
}

// 获取Header内容
func (this *AccessLog) GetHeader(name string) string {
	v, found := this.Header[name]
	if !found {
		return ""
	}
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// 分析mime，扩展名、代理设置等
func (this *AccessLog) Parse() {
	if (len(this.writingFields) == 0 || lists.ContainsInt(this.writingFields, AccessLogFieldExtend)) || this.shouldStat {
		this.parseMime()
		this.parseExtension()
		this.parseUserAgent()
		this.parseGeoIP()
	}
}

// 是否支持统计
func (this *AccessLog) ShouldStat() bool {
	return this.shouldStat
}

// 设置是否支持统计
func (this *AccessLog) SetShouldStat(b bool) {
	this.shouldStat = b
}

// 是否支持写入
func (this *AccessLog) ShouldWrite() bool {
	return this.shouldWrite
}

// 设置是否写入
func (this *AccessLog) SetShouldWrite(b bool) {
	this.shouldWrite = b
}

// 设置写入的字段
func (this *AccessLog) SetWritingFields(fields []int) {
	this.writingFields = fields
}

// 清除不必要的的字段
func (this *AccessLog) CleanFields() {
	l := len(this.writingFields)
	if l == 0 {
		return
	}
	for _, code := range AccessLogFieldsCodes {
		if lists.ContainsInt(this.writingFields, code) {
			continue
		}
		switch code {
		case AccessLogFieldHeader:
			this.Header = nil
		case AccessLogFieldSentHeader:
			this.SentHeader = nil
		case AccessLogFieldArg:
			this.Arg = nil
		case AccessLogFieldCookie:
			this.Cookie = nil
		case AccessLogFieldExtend:
			this.Extend = nil
		case AccessLogFieldReferer:
			this.Referer = ""
		case AccessLogFieldUserAgent:
			this.UserAgent = ""
		}
	}
}

// 设置数据库列值
func (this *AccessLog) SetDBColumns(v maps.Map) {
	idString := v.GetString("_id")
	if len(idString) > 0 {
		id, err := shared.ObjectIdFromHex(idString)
		if err != nil {
			logs.Error(err)
		} else {
			this.Id = id
		}
	}
	this.ServerId = v.GetString("serverId")
	this.BackendId = v.GetString("backendId")
	this.LocationId = v.GetString("locationId")
	this.FastcgiId = v.GetString("fastcgiId")
	this.RewriteId = v.GetString("rewriteId")
	this.TeaVersion = v.GetString("teaVersion")
	this.RemoteAddr = v.GetString("remoteAddr")
	this.RemotePort = v.GetInt("remotePort")
	this.RemoteUser = v.GetString("remoteUser")
	this.RequestURI = v.GetString("requestURI")
	this.RequestPath = v.GetString("requestPath")
	this.RequestLength = v.GetInt64("requestLength")
	this.RequestTime = v.GetFloat64("requestTime")
	this.RequestMethod = v.GetString("requestMethod")
	this.RequestFilename = v.GetString("requestFilename")
	this.Scheme = v.GetString("scheme")
	this.Proto = v.GetString("proto")
	this.BytesSent = v.GetInt64("bytesSent")
	this.BodyBytesSent = v.GetInt64("bodyBytesSent")
	this.Status = v.GetInt("status")
	this.StatusMessage = v.GetString("statusMessage")
	this.jsonDecode(v.Get("sentHeader"), &this.SentHeader)
	this.TimeISO8601 = v.GetString("timeISO8601")
	this.TimeLocal = v.GetString("timeLocal")
	this.Msec = v.GetFloat64("msec")
	this.Timestamp = v.GetInt64("timestamp")
	this.Host = v.GetString("host")
	this.Referer = v.GetString("referer")
	this.UserAgent = v.GetString("userAgent")
	this.Request = v.GetString("request")
	this.ContentType = v.GetString("contentType")
	this.jsonDecode(v.Get("cookie"), &this.Cookie)
	this.jsonDecode(v.Get("arg"), &this.Arg)
	this.Args = v.GetString("args")
	this.QueryString = v.GetString("queryString")
	this.jsonDecode(v.Get("header"), &this.Header)
	this.ServerName = v.GetString("serverName")
	this.ServerPort = v.GetInt("serverPort")
	this.ServerProtocol = v.GetString("serverProtocol")
	this.BackendAddress = v.GetString("backendAddress")
	this.FastcgiAddress = v.GetString("fastcgiAddress")
	this.RequestData = this.toBytes(v.Get("requestData"))
	this.ResponseHeaderData = this.toBytes(v.Get("responseHeaderData"))
	this.ResponseBodyData = this.toBytes(v.Get("responseBodyData"))
	this.jsonDecode(v.Get("errors"), &this.Errors)
	this.HasErrors = v.GetInt("hasErrors") > 0
	this.jsonDecode(v.Get("extend"), &this.Extend)
	this.jsonDecode(v.Get("attrs"), &this.Attrs)
}

// 获取数据库列值
func (this *AccessLog) DBColumns() maps.Map {
	if this.Id.IsZero() {
		this.Id = shared.NewObjectId()
	}
	hasErrors := 0
	if this.HasErrors {
		hasErrors = 1
	}
	return maps.Map{
		"_id":                this.Id.Hex(),
		"serverId":           this.ServerId,
		"backendId":          this.BackendId,
		"locationId":         this.LocationId,
		"fastcgiId":          this.FastcgiId,
		"rewriteId":          this.RewriteId,
		"teaVersion":         this.TeaVersion,
		"remoteAddr":         this.RemoteAddr,
		"remotePort":         this.RemotePort,
		"remoteUser":         this.RemoteUser,
		"requestURI":         this.RequestURI,
		"requestPath":        this.RequestPath,
		"requestLength":      this.RequestLength,
		"requestTime":        this.RequestTime,
		"requestMethod":      this.RequestMethod,
		"requestFilename":    this.RequestFilename,
		"scheme":             this.Scheme,
		"proto":              this.Proto,
		"bytesSent":          this.BytesSent,
		"bodyBytesSent":      this.BodyBytesSent,
		"status":             this.Status,
		"statusMessage":      this.StatusMessage,
		"sentHeader":         this.jsonEncode(this.SentHeader),
		"timeISO8601":        this.TimeISO8601,
		"timeLocal":          this.TimeLocal,
		"msec":               this.Msec,
		"timestamp":          this.Timestamp,
		"host":               this.Host,
		"referer":            this.Referer,
		"userAgent":          this.UserAgent,
		"request":            this.Request,
		"contentType":        this.ContentType,
		"cookie":             this.jsonEncode(this.Cookie),
		"arg":                this.jsonEncode(this.Arg),
		"args":               this.Args,
		"queryString":        this.QueryString,
		"header":             this.jsonEncode(this.Header),
		"serverName":         this.ServerName,
		"serverPort":         this.ServerPort,
		"serverProtocol":     this.ServerProtocol,
		"backendAddress":     this.BackendAddress,
		"fastcgiAddress":     this.FastcgiAddress,
		"requestData":        this.RequestData,
		"responseHeaderData": this.ResponseHeaderData,
		"responseBodyData":   this.ResponseBodyData,
		"errors":             this.jsonEncode(this.Errors),
		"hasErrors":          hasErrors,
		"extend":             this.jsonEncode(this.Extend),
		"attrs":              this.jsonEncode(this.Attrs),
	}
}

func (this *AccessLog) jsonEncode(v interface{}) []byte {
	if v == nil {
		return []byte("null")
	}
	value := reflect.ValueOf(v)
	if !value.IsValid() || value.IsZero() {
		return []byte("null")
	}
	data, _ := ffjson.Marshal(v)
	return data
}

func (this *AccessLog) jsonDecode(data interface{}, vPtr interface{}) {
	if data == nil {
		return
	}
	b, ok := data.([]byte)
	if ok {
		_ = ffjson.Unmarshal(b, vPtr)
	}
	s, ok := data.(string)
	if ok {
		_ = ffjson.Unmarshal([]byte(s), vPtr)
	}
}

func (this *AccessLog) toBytes(v interface{}) []byte {
	if v == nil {
		return nil
	}
	b, ok := v.([]byte)
	if ok {
		return b
	}
	s, ok := v.(string)
	if ok {
		return []byte(s)
	}
	return nil
}

func (this *AccessLog) parseMime() {
	if this.Extend == nil {
		this.Extend = &AccessLogExtend{}
	}

	semicolonIndex := strings.Index(this.ContentType, ";")
	if semicolonIndex == -1 {
		this.Extend.File.MimeType = this.ContentType
		this.Extend.File.Charset = ""
		return
	}

	this.Extend.File.MimeType = this.ContentType[:semicolonIndex]
	match := charsetReg.FindStringSubmatch(this.ContentType)
	if len(match) > 0 {
		this.Extend.File.Charset = strings.ToUpper(match[1])
	} else {
		this.Extend.File.Charset = ""
	}
}

func (this *AccessLog) parseExtension() {
	if this.Extend == nil {
		this.Extend = &AccessLogExtend{}
	}
	ext := filepath.Ext(this.RequestPath)
	if len(ext) == 0 {
		this.Extend.File.Extension = ""
	} else {
		this.Extend.File.Extension = strings.ToLower(ext[1:])
	}
}

func (this *AccessLog) parseUserAgent() {
	// MDN上的参考：https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent
	// 浏览器集合测试：http://www.browserscope.org/

	if len(this.UserAgent) == 0 {
		return
	}

	cacheKey := []byte(this.UserAgent)
	item := userAgentGrid.Read(cacheKey)
	var result *uaparser.UserAgent = nil
	if item != nil {
		if item.ValueInterface == nil {
			return
		}
		result = item.ValueInterface.(*uaparser.UserAgent)
	} else {
		userAgent, found := userAgentParser.Parse(this.UserAgent)
		if found {
			result = userAgent
		}
		userAgentGrid.WriteInterface(cacheKey, userAgent, 3600)
	}

	if result != nil {
		if this.Extend == nil {
			this.Extend = &AccessLogExtend{}
		}
		this.Extend.Client = AccessLogClient{}
		if result.Browser != nil {
			this.Extend.Client.Browser = AccessLogClientBrowser{
				Family: result.Browser.Family,
				Major:  result.Browser.Major,
				Minor:  result.Browser.Minor,
				Patch:  result.Browser.Patch,
			}
		}
		if result.OS != nil {
			this.Extend.Client.OS = AccessLogClientOS{
				Family:     result.OS.Family,
				Major:      result.OS.Major,
				Minor:      result.OS.Minor,
				Patch:      result.OS.Patch,
				PatchMinor: result.OS.PatchMinor,
			}
		}
		if result.Device != nil {
			this.Extend.Client.Device = AccessLogClientDevice{
				Family: result.Device.Family,
				Brand:  result.Device.Brand,
				Model:  result.Device.Model,
			}
		}
	}
}

func (this *AccessLog) parseGeoIP() {
	if teageo.DB == nil {
		return
	}

	if len(this.RemoteAddr) == 0 {
		return
	}

	// 是否为本地
	// 参考 https://tools.ietf.org/html/rfc1918
	if this.RemoteAddr == "127.0.0.1" || strings.HasPrefix(this.RemoteAddr, "10.") || strings.HasPrefix(this.RemoteAddr, "192.168.") || strings.HasPrefix(this.RemoteAddr, "172.16.") {
		return
	}

	// 参考：https://dev.maxmind.com/geoip/geoip2/geolite2/
	record, err := teageo.IP2City(this.RemoteAddr, true)
	if err != nil {
		logs.Error(err)
		return
	}

	if record == nil {
		return
	}

	if this.Extend == nil {
		this.Extend = &AccessLogExtend{}
	}
	this.Extend.Geo.Location.AccuracyRadius = record.Location.AccuracyRadius
	this.Extend.Geo.Location.MetroCode = record.Location.MetroCode
	this.Extend.Geo.Location.TimeZone = record.Location.TimeZone
	this.Extend.Geo.Location.Latitude = record.Location.Latitude
	this.Extend.Geo.Location.Longitude = record.Location.Longitude

	if len(record.Country.Names) > 0 {
		name, found := record.Country.Names["zh-CN"]
		if found {
			this.Extend.Geo.Region = teageo.ConvertName(name)
		}
	}

	if len(record.Subdivisions) > 0 && len(record.Subdivisions[0].Names) > 0 {
		name, found := record.Subdivisions[0].Names["zh-CN"]
		if found {
			this.Extend.Geo.State = teageo.ConvertName(name)
		}
	}

	if len(record.City.Names) > 0 {
		name, found := record.City.Names["zh-CN"]
		if found {
			this.Extend.Geo.City = name
		}
	}
}
