package teaplugins

import (
	"bufio"
	"bytes"
	"github.com/TeaWeb/plugin/pkg/messages"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

var plugins = []*Plugin{}
var pluginsLocker = &sync.Mutex{}

var requestFilters = []func(req []byte) (result []byte, willContinue bool){}
var HasRequestFilters = false
var responseFilters = []func(resp []byte) (result []byte, willContinue bool){}
var HasResponseFilters = false

func Register(plugin *Plugin) {
	pluginsLocker.Lock()
	plugins = append(plugins, plugin)
	pluginsLocker.Unlock()
}

func Plugins() []*Plugin {
	return plugins
}

func FilterRequest(request *http.Request) (resultReq *http.Request, willContinue bool) {
	if !HasRequestFilters {
		return request, true
	}

	data, err := httputil.DumpRequest(request, true)
	if err != nil {
		logs.Error(err)
		return request, true
	}

	defer func() {
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(data)))
		if err != nil {
			logs.Error(err)
			return
		}

		resultReq = req
	}()

	for _, f := range requestFilters {
		result, willContinue := f(data)

		data = result

		if !willContinue {
			return resultReq, false
		}
	}

	return resultReq, true
}

func FilterResponse(response *http.Response) (resultResp *http.Response) {
	if !HasResponseFilters {
		return response
	}

	data, err := httputil.DumpResponse(response, true)
	if err != nil {
		logs.Error(err)
		return response
	}

	defer func() {
		resp, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(data)), nil)
		if err != nil {
			logs.Error(err)
			return
		}

		resultResp = resp
	}()

	for _, f := range responseFilters {
		result, willContinue := f(data)

		data = result

		if !willContinue {
			return resultResp
		}
	}

	return resultResp
}

// 刷新所有插件的App
func ReloadAllApps() {
	for _, loader := range loaders {
		message := new(messages.ReloadAppsAction)
		loader.Write(message)
	}
}

// 加载插件
var loaders = []*Loader{}

func load() {
	logs.Println("[plugin]load plugins")
	dir := Tea.Root + Tea.DS + "plugins"
	files.NewFile(dir).Range(func(file *files.File) {
		if !strings.HasSuffix(file.Name(), ".tea") && !strings.HasSuffix(file.Name(), ".tea.exe") {
			return
		}

		logs.Println("[plugin][loader]load plugin '" + file.Name() + "'")

		loader := NewLoader(file.Path())
		go func() {
			err := loader.Load()
			if err != nil {
				logs.Println("[plugin][" + file.Name() + "]failed:" + err.Error())
			}
		}()
		loaders = append(loaders, loader)
	})
}
