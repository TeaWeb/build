package teaplugins

import (
	"encoding/binary"
	"errors"
	"github.com/TeaWeb/plugin/pkg/messages"
	"github.com/iwind/TeaGo/logs"
	"io"
	"path/filepath"
	"reflect"
)

// 插件加载器
type Loader struct {
	path   string
	plugin *Plugin

	methods   map[string]reflect.Method
	thisValue reflect.Value

	writer PipeInterface

	debug bool
}

type PipeInterface interface {
	Read([]byte) (n int, err error)
	Write([]byte) (n int, err error)
}

func NewLoader(path string) *Loader {
	loader := &Loader{
		path:    path,
		methods: map[string]reflect.Method{},
	}

	// 当前methods
	t := reflect.TypeOf(loader)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		loader.methods[method.Name] = method
	}

	loader.thisValue = reflect.ValueOf(loader)

	return loader
}

func (this *Loader) Debug() {
	this.debug = true
}

func (this *Loader) pipe(reader PipeInterface, writer PipeInterface) {
	buf := make([]byte, 1024)
	msgData := []byte{}
	for {
		if this.debug {
			logs.Println("[plugin][" + this.shortFileName() + "]try to read buf")
		}

		n, err := reader.Read(buf)

		if n > 0 {
			msgData = append(msgData, buf[:n]...)

			if this.debug {
				logs.Println("[plugin]["+this.shortFileName()+"]len:", len(msgData), ",", "read msg data:", string(msgData))
			}

			msgLen := uint32(len(msgData))
			h := uint32(24) // header length

			if msgLen > h { // 数据组成方式： | actionLen[8] | dataLen[8] | action | data[len-8]
				id := binary.BigEndian.Uint32(msgData[:8])
				l1 := binary.BigEndian.Uint32(msgData[8:16])
				l2 := binary.BigEndian.Uint32(msgData[16:24])

				if msgLen >= h+l1+l2 { // 数据已经完整了
					action := string(msgData[h : h+l1])
					valueData := msgData[h+l1 : h+l1+l2]

					msgData = msgData[h+l1+l2:]

					ptr, err := messages.Unmarshal(action, valueData)
					if err != nil {
						logs.Println("[plugin]["+this.shortFileName()+"[unmarshal message error:", err.Error())
						continue
					}

					err = this.CallAction(ptr, id)
					if err != nil {
						logs.Println("[plugin]["+this.shortFileName()+"]call action error:", err.Error())
						continue
					}
				}
			}
		}

		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				logs.Println("[plugin][" + this.shortFileName() + "]break:" + err.Error())
			}
			break
		}
	}
}

func (this *Loader) CallAction(ptr interface{}, messageId uint32) error {
	action, ok := ptr.(messages.ActionInterface)
	if !ok {
		return errors.New("ptr should be an action")
	}
	action.SetMessageId(messageId)

	method, found := this.methods["Action"+action.Name()]
	if !found {
		return errors.New("handler for '" + action.Name() + "' not found")

	}
	method.Func.Call([]reflect.Value{this.thisValue, reflect.ValueOf(action)})
	return nil
}

func (this *Loader) ActionRegisterPlugin(action *messages.RegisterPluginAction) {
	if this.plugin != nil {
		logs.Println("[plugin][" + this.shortFileName() + "]load only one plugin from one file")
		return
	}

	// 添加到插件中
	if action.Plugin != nil {
		p := action.Plugin
		p2 := NewPlugin()
		p2.IsExternal = true
		p2.Name = p.Name
		p2.Description = p.Description
		p2.Code = p.Code
		p2.Site = p.Site
		p2.Developer = p.Developer
		p2.Date = p.Date
		p2.Version = p.Version
		p2.HasRequestFilter = p.HasRequestFilter
		p2.HasResponseFilter = p.HasResponseFilter

		// request filter
		if p.HasRequestFilter {
			requestFilters = append(requestFilters, func(data []byte) (result []byte, willContinue bool) {
				action := &messages.FilterRequestAction{
					Data: data,
				}
				this.Write(action)

				respAction := messages.ActionQueue.Wait(action)
				r, ok := respAction.(*messages.FilterRequestAction)
				if ok {
					return r.Data, r.Continue
				} else {
					return action.Data, true
				}
			})
			HasRequestFilters = true
		}

		// response filter
		if p.HasResponseFilter {
			responseFilters = append(responseFilters, func(data []byte) (result []byte, willContinue bool) {
				action := &messages.FilterResponseAction{
					Data: data,
				}
				this.Write(action)

				respAction := messages.ActionQueue.Wait(action)
				r, ok := respAction.(*messages.FilterResponseAction)
				if ok {
					return r.Data, r.Continue
				} else {
					return action.Data, true
				}
			})
			HasResponseFilters = true
		}

		Register(p2)

		this.plugin = p2
	}

	// 发送启动信息
	this.Write(&messages.StartAction{})
}

func (this *Loader) ActionFilterRequest(action *messages.FilterRequestAction) {
	messages.ActionQueue.Notify(action)
}

func (this *Loader) ActionFilterResponse(action *messages.FilterResponseAction) {
	messages.ActionQueue.Notify(action)
}

func (this *Loader) Write(action messages.ActionInterface) error {
	msg := messages.NewActionMessage(action)
	msg.Id = action.MessageId()
	data, err := msg.Marshal()
	if err != nil {
		return err
	}
	action.SetMessageId(msg.Id)
	if this.writer == nil {
		return nil
	}
	_, err = this.writer.Write(data)
	return err
}

// 插件文件短名称
func (this *Loader) shortFileName() string {
	return filepath.Base(this.path)
}
