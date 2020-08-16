package scripts

import (
	"encoding/json"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/stats"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/caches"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"github.com/robertkrimen/otto"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var engineCache = caches.NewFactory()

// 脚本引擎
type Engine struct {
	vm           *otto.Otto
	chartOptions []maps.Map
	widgetCodes  map[string]maps.Map // "code" => { name, ..., definition:FUNCTION CODE }
	output       []string
	context      *Context
	mongoEnabled bool
}

// 获取新引擎
func NewEngine() *Engine {
	engine := &Engine{
		chartOptions: []maps.Map{},
		widgetCodes:  map[string]maps.Map{},
	}
	engine.init()
	return engine
}

// 设置MongoDB是否可用
func (this *Engine) SetMongo(b bool) {
	this.mongoEnabled = b
}

// 设置上下文信息
func (this *Engine) SetContext(context *Context) {
	this.context = context

	if context.Server != nil {
		runningServer := teaproxy.SharedManager.FindServer(context.Server.Id)

		options := map[string]interface{}{
			"isOn":        context.Server.On,
			"id":          context.Server.Id,
			"name":        context.Server.Name,
			"filename":    context.Server.Filename,
			"description": context.Server.Description,
			"listen":      context.Server.Listen,
			"http":        context.Server.Http,
			"backends": lists.Map(context.Server.Backends, func(k int, v interface{}) interface{} {
				backend := v.(*teaconfigs.BackendConfig)

				if runningServer != nil {
					runningBackend := runningServer.FindBackend(backend.Id)
					if runningBackend != nil {
						backend.IsDown = runningBackend.IsDown
					}
				}

				return map[string]interface{}{
					"isOn":     backend.On,
					"weight":   backend.Weight,
					"id":       backend.Id,
					"isDown":   backend.IsDown,
					"isBackup": backend.IsBackup,
					"address":  backend.Address,
					"code":     backend.Code,
				}
			}),
			"locations": lists.Map(context.Server.Locations, func(k int, v interface{}) interface{} {
				location := v.(*teaconfigs.LocationConfig)
				err := location.Validate()
				if err != nil {
					logs.Error(err)
				}
				locationOptions := map[string]interface{}{
					"id":          location.Id,
					"isOn":        location.On,
					"pattern":     location.PatternString(),
					"cachePolicy": location.CachePolicy,
					"fastcgi": lists.Map(location.Fastcgi, func(k int, v interface{}) interface{} {
						fastcgi := v.(*teaconfigs.FastcgiConfig)
						return map[string]interface{}{
							"id":   fastcgi.Id,
							"isOn": fastcgi.On,
							"pass": fastcgi.Pass,
						}
					}),
					"rewrite": lists.Map(location.Rewrite, func(k int, v interface{}) interface{} {
						rewrite := v.(*teaconfigs.RewriteRule)
						return map[string]interface{}{
							"id":      rewrite.Id,
							"isOn":    rewrite.On,
							"pattern": rewrite.Pattern,
							"replace": rewrite.Replace,
						}
					}),
					"root":    location.Root,
					"index":   location.Index,
					"headers": location.Headers,
					"backends": lists.Map(location.Backends, func(k int, v interface{}) interface{} {
						backend := v.(*teaconfigs.BackendConfig)

						if runningServer != nil {
							runningBackend := runningServer.FindBackend(backend.Id)
							if runningBackend != nil {
								backend.IsDown = runningBackend.IsDown
							}
						}

						return map[string]interface{}{
							"isOn":     backend.On,
							"weight":   backend.Weight,
							"id":       backend.Id,
							"isDown":   backend.IsDown,
							"isBackup": backend.IsBackup,
							"address":  backend.Address,
							"code":     backend.Code,
						}
					}),
				}
				if location.Websocket != nil && location.Websocket.On {
					locationOptions["websocket"] = maps.Map{
						"isOn": true,
					}
				} else {
					locationOptions["websocket"] = nil
				}
				return locationOptions
			}),
		}

		if context.Server.SSL != nil {
			options["ssl"] = maps.Map{
				"isOn":   context.Server.SSL.On,
				"listen": context.Server.SSL.Listen,
			}
		} else {
			options["ssl"] = maps.Map{
				"isOn":   false,
				"listen": []string{},
			}
		}

		_, err := this.vm.Run(`context.server = new http.Server(` + stringutil.JSONEncode(options) + `);`)
		if err != nil {
			logs.Error(err)
		}
	}

	// 可供使用的特性
	features := []string{}
	if this.mongoEnabled {
		features = append(features, "mongo")
	}
	features = append(features, runtime.GOOS)
	features = append(features, runtime.GOARCH)
	_, err := this.vm.Run(`context.features=` + stringutil.JSONEncode(features) + `;`)
	if err != nil {
		logs.Error(err)
	}
}

// 初始化
func (this *Engine) init() {
	this.vm = otto.New()
	err := this.vm.Set("callConsoleLog", this.callConsoleLog)
	if err != nil {
		logs.Error(err)
	}

	_, err = this.vm.Run("console.log = callConsoleLog;")
	if err != nil {
		logs.Error(err)
	}

	this.loadLib("libs/array.js")
	this.loadLib("libs/times.js")
	this.loadLib("libs/caches.js")
	this.loadLib("libs/mutex.js")
	this.loadLib("libs/server.logs.js")
	this.loadLib("libs/server.stat.js")
	this.loadLib("libs/http.js")
	this.loadLib("libs/colors.js")
	this.loadLib("libs/widgets.js")
	this.loadLib("libs/charts.js")
	this.loadLib("libs/charts.menu.js")
	this.loadLib("libs/charts.gauge.js")
	this.loadLib("libs/charts.html.js")
	this.loadLib("libs/charts.line.js")
	this.loadLib("libs/charts.pie.js")
	this.loadLib("libs/charts.progress.js")
	this.loadLib("libs/charts.stackbar.js")
	this.loadLib("libs/charts.url.js")
	this.loadLib("libs/charts.table.js")
	this.loadLib("libs/context.js")

	err = this.vm.Set("callSetCache", this.callSetCache)
	if err != nil {
		logs.Error(err)
	}

	err = this.vm.Set("callGetCache", this.callGetCache)
	if err != nil {
		logs.Error(err)
	}

	err = this.vm.Set("callChartRender", this.callRenderChart)
	if err != nil {
		logs.Error(err)
	}

	err = this.vm.Set("callLogExecuteQuery", this.callLogExecuteQuery)
	if err != nil {
		logs.Error(err)
	}

	err = this.vm.Set("callStatExecuteQuery", this.callStatExecuteQuery)
	if err != nil {
		logs.Error(err)
	}
}

// 运行Widget代码
func (this *Engine) RunCode(code string) error {
	_, err := this.vm.Run(`(function () {` + code + `
	
	widget.callRun();
})();`)
	return err
}

// 获取控制台输出
func (this *Engine) Output() []string {
	if this.output == nil {
		return []string{}
	}
	return this.output
}

// 获取Widget中的图表对象
func (this *Engine) Charts() []maps.Map {
	return this.chartOptions
}

func (this *Engine) callConsoleLog(call otto.FunctionCall) otto.Value {
	values := []string{}
	for _, v := range call.ArgumentList {
		i, err := v.Export()
		if err != nil {
			values = append(values, v.String())
		} else {
			values = append(values, stringutil.JSONEncodePretty(i))
		}
	}
	s := strings.Join(values, ", ")
	//logs.Println("[console]", s)

	this.output = append(this.output, s)

	return otto.UndefinedValue()
}

func (this *Engine) callRenderChart(call otto.FunctionCall) otto.Value {
	obj := call.Argument(0)
	v, err := obj.Export()
	if err != nil {
		logs.Error(err)
		return otto.UndefinedValue()
	}
	m := maps.NewMap(v)

	menus, err := obj.Object().Get("menus")
	if err != nil {
		logs.Error(err)
	} else {
		menusV, err := menus.Export()
		if err != nil {
			logs.Error(err)
		} else {
			m["menus"] = menusV
		}
	}

	options, err := obj.Object().Get("options")
	if err != nil {
		logs.Error(err)
	} else {
		v, err := options.Export()
		if err != nil {
			logs.Error(err)
		} else {
			m["options"] = maps.NewMap(v)
		}
	}

	this.chartOptions = append(this.chartOptions, m)
	return otto.UndefinedValue()
}

func (this *Engine) callLogExecuteQuery(call otto.FunctionCall) otto.Value {
	if this.context == nil {
		logs.Error(errors.New("'context' should not be nil"))
		return otto.UndefinedValue()
	}

	if this.context.Server == nil {
		logs.Error(errors.New("'context.server' should not be nil"))
		return otto.UndefinedValue()
	}

	arg, err := call.Argument(0).Export()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	m := maps.NewMap(arg)

	action := m.GetString("action")
	if len(action) == 0 {
		this.throw(errors.New("'action' should not be empty"))
		return otto.UndefinedValue()
	}

	day := timeutil.Format("Ymd")
	query := teadb.NewQuery(teadb.AccessLogDAO().TableName(day))

	// group
	group := m.Get("group")
	if group != nil {
		logs.Error(errors.New("unsupported method 'group()'"))
		return otto.UndefinedValue()
	}

	// cond
	cond := m.Get("cond")
	if cond != nil && reflect.TypeOf(cond).Kind() == reflect.Map {
		m, ok := cond.(map[string]interface{})
		if ok {
			for field, ops := range m {
				opsMap, ok := ops.(map[string]interface{})
				if ok {
					for op, v := range opsMap {
						query.Op(field, op, v)
					}
				}
			}
		}
	}

	// result
	result := m.Get("result")
	if types.IsSlice(result) {
		lists.Each(result, func(k int, v interface{}) {
			query.Result(types.String(v))
		})
	}

	// timeFrom
	timeFrom := m.GetInt64("timeFrom")
	if timeFrom > 0 {
		query.Gte("timestamp", timeFrom)
	}

	// timeTo
	timeTo := m.GetInt64("timeTo")
	if timeTo > 0 {
		query.Lte("timestamp", timeTo)
	}

	// offset & size
	query.Offset(m.GetInt("offset"))
	query.Limit(m.GetInt("size"))

	// sort
	sorts := m.Get("sorts")
	if sorts != nil {
		sortsMap, ok := sorts.([]map[string]interface{})
		if ok {
			for _, m := range sortsMap {
				for k, v := range m {
					vInt := types.Int(v)
					if len(k) == 0 {
						k = "_id"
					}
					if vInt < 0 {
						query.Desc(k)
					} else {
						query.Asc(k)
					}
				}
			}
		}
	}

	// 开始执行
	var v interface{} = nil
	switch action {
	case "findAll":
		v, err = teadb.AccessLogDAO().QueryAccessLogs(day, this.context.Server.Id, query)
	default:
		logs.Error(errors.New("unsupported action '" + action + "'"))
		return otto.UndefinedValue()
	}
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}

	jsValue, err := this.toValue(v)
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	return jsValue
}

func (this *Engine) callStatExecuteQuery(call otto.FunctionCall) otto.Value {
	if this.context == nil {
		this.throw(errors.New("'context' should not be nil"))
		return otto.UndefinedValue()
	}

	if this.context.Server == nil {
		this.throw(errors.New("'context.server' should not be nil"))
		return otto.UndefinedValue()
	}

	arg, err := call.Argument(0).Export()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	m := maps.NewMap(arg)

	action := m.GetString("action")
	if len(action) == 0 {
		this.throw(errors.New("'action' should not be empty"))
		return otto.UndefinedValue()
	}

	query := teadb.NewQuery(teadb.ServerValueDAO().TableName(this.context.Server.Id))

	// cond
	cond := m.Get("cond")
	if cond != nil && reflect.TypeOf(cond).Kind() == reflect.Map {
		m, ok := cond.(map[string]interface{})
		if ok {
			for field, ops := range m {
				opsMap, ok := ops.(map[string]interface{})
				if ok {
					for op, v := range opsMap {
						query.Op(field, op, v)
					}
				}
			}
		}
	}

	// offset & size
	query.Offset(m.GetInt("offset"))
	query.Limit(m.GetInt("size"))

	// sort
	sorts := m.Get("sorts")
	if sorts != nil {
		sortsMap, ok := sorts.([]map[string]interface{})
		if ok {
			for _, m := range sortsMap {
				for k, v := range m {
					if len(k) == 0 {
						k = "_id"
					}
					vInt := types.Int(v)
					if vInt < 0 {
						query.Desc(k)
					} else {
						query.Asc(k)
					}
				}
			}
		}
	}

	// 开始执行
	var v interface{} = nil
	switch action {
	case "findAll":
		v, err = teadb.ServerValueDAO().QueryValues(query)
	case "find":
		result, err1 := teadb.ServerValueDAO().
			QueryValues(query)
		if err1 != nil {
			err = err1
		}
		if len(result) > 0 {
			v = result[0]
		}
	default:
		logs.Error(errors.New("unsupported action '" + action + "'"))
		return otto.UndefinedValue()
	}
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}

	jsValue, err := this.toValue(v)
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	return jsValue
}

func (this *Engine) callSetCache(call otto.FunctionCall) otto.Value {
	key, err := call.Argument(0).ToString()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}

	key = stringutil.Md5(key)

	value, err := call.Argument(1).Export()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}

	life, err := call.Argument(2).Export()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	lifeSeconds := types.Int64(life)

	engineCache.Set(key, value, time.Duration(lifeSeconds)*time.Second)

	return otto.UndefinedValue()
}

func (this *Engine) callGetCache(call otto.FunctionCall) otto.Value {
	key, err := call.Argument(0).ToString()
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}

	key = stringutil.Md5(key)

	value, found := engineCache.Get(key)
	if !found {
		return otto.UndefinedValue()
	}
	v, err := this.vm.ToValue(value)
	if err != nil {
		this.throw(err)
		return otto.UndefinedValue()
	}
	return v
}

// 加载JS库文件
func (this *Engine) loadLib(file string) {
	path := teautils.WebRoot() + Tea.DS + file
	cacheKey := "libfile://" + path
	code, found := engineCache.Get(cacheKey)
	if !found {
		var err error = nil
		code, err = files.NewFile(path).ReadAllString()
		if err != nil {
			logs.Error(err)
			return
		}
		engineCache.Set(cacheKey, code)
	}

	_, err := this.vm.Run(code)
	if err != nil {
		logs.Error(err)
		return
	}
}

func (this *Engine) toValue(data interface{}) (v otto.Value, err error) {
	if data == nil {
		return this.vm.ToValue(data)
	}

	// *AccessLog
	if _, ok := data.(*accesslogs.AccessLog); ok {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return this.vm.ToValue(data)
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(jsonData, &m)
		if err != nil {
			logs.Error(err)
			return this.vm.ToValue(data)
		}
		return this.vm.ToValue(m)
	}

	// []*AccessLog
	if _, ok := data.([]*accesslogs.AccessLog); ok {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return this.vm.ToValue(data)
		}
		m := []map[string]interface{}{}
		err = json.Unmarshal(jsonData, &m)
		if err != nil {
			logs.Error(err)
			return this.vm.ToValue(data)
		}
		return this.vm.ToValue(m)
	}

	// *Value
	if _, ok := data.(*stats.Value); ok {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return this.vm.ToValue(data)
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(jsonData, &m)
		if err != nil {
			logs.Error(err)
			return this.vm.ToValue(data)
		}
		return this.vm.ToValue(m)
	}

	if _, ok := data.([]interface{}); ok {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return this.vm.ToValue(data)
		}
		m := []map[string]interface{}{}
		err = json.Unmarshal(jsonData, &m)
		if err != nil {
			logs.Error(err)
			return this.vm.ToValue(data)
		}

		return this.vm.ToValue(m)
	}

	if _, ok := data.([]*stats.Value); ok {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return this.vm.ToValue(data)
		}
		m := []map[string]interface{}{}
		err = json.Unmarshal(jsonData, &m)
		if err != nil {
			logs.Error(err)
			return this.vm.ToValue(data)
		}

		return this.vm.ToValue(m)
	}

	return this.vm.ToValue(data)
}

func (this *Engine) throw(err error) {
	if err != nil {
		value, _ := this.vm.Call("new Error", nil, err.Error())
		panic(value)
	}
}
