package forms

import (
	"encoding/json"
	"errors"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/robertkrimen/otto"
	"net/http"
	"reflect"
	"strings"
)

type Form struct {
	Namespace     string            `yaml:"namespace" json:"namespace"`
	Groups        []*Group          `yaml:"groups" json:"groups"`
	Javascript    string            `yaml:"javascript" json:"javascript"`
	CSS           string            `yaml:"css" json:"css"`
	ValidateCode  string            `yaml:"validateCode" json:"validateCode"`
	ComposedAttrs map[string]string `yaml:"composedAttrs" json:"composedAttrs"`
	vm            *otto.Otto
}

func FormDecode(data []byte, form *Form) error {
	m := map[string]interface{}{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	namespace, found := m["namespace"]
	if found {
		form.Namespace = types.String(namespace)
	}

	validateCode, found := m["validateCode"]
	if found {
		form.ValidateCode = types.String(validateCode)
	}

	groups, found := m["groups"]
	form.Groups = []*Group{}
	if found && groups != nil {
		if reflect.TypeOf(groups).Kind() == reflect.Slice {
			lists.Each(groups, func(k int, v interface{}) {
				if v == nil {
					return
				}
				group := &Group{}
				form.Groups = append(form.Groups, group)

				vMap := maps.NewMap(v)
				group.Namespace = vMap.GetString("namespace")

				elements := vMap.GetSlice("elements")
				lists.Each(elements, func(k int, v interface{}) {
					vMap := maps.NewMap(v)
					classType := vMap.GetString("classType")
					instance, found := allElementTypes[classType]
					if !found || instance == nil {
						return
					}

					obj := reflect.New(reflect.TypeOf(instance).Elem()).Interface().(ElementInterface)
					err := teautils.MapToObjectJSON(vMap, obj)
					if err != nil {
						logs.Error(err)
					}
					group.Elements = append(group.Elements, obj)
				})
			})
		}
	}

	return nil
}

func NewForm(namespace string) *Form {
	vm := otto.New()
	_, err := vm.Run(`
function FieldError(field, message) {
	return {
		"field": field,
		"error": message
	};
}
`)
	if err != nil {
		logs.Error(err)
	}
	return &Form{
		Namespace: namespace,
		vm:        vm,
	}
}

func (this *Form) NewGroup() *Group {
	group := &Group{
		Namespace: this.Namespace,
	}
	this.Groups = append(this.Groups, group)
	return group
}

func (this *Form) Compose() {
	this.CSS = "<style type=\"text/css\">\n"
	this.Javascript = ""

	for _, g := range this.Groups {
		g.ComposedAttrs = this.ComposedAttrs
		g.Compose()

		for _, e := range g.Elements {
			s := e.Super()

			if len(s.CSS) > 0 {
				this.CSS += s.CSS + "\n\n"
			}
			if len(s.Javascript) > 0 {
				this.Javascript += s.Javascript + "\n\n"
			}
		}
	}
	this.CSS += "</style>\n"
}

func (this *Form) Init(values map[string]interface{}) {
	eMap := map[string]ElementInterface{} // code => []element
	for _, g := range this.Groups {
		for _, e := range g.Elements {
			s := e.Super()
			eMap[s.Code] = e

			// 处理不在values中的元素的init
			if len(s.InitCode) > 0 {
				_, ok := values[s.Code]
				if !ok {
					newValue, err := this.vm.Run("(function() { var values = " + stringutil.JSONEncode(values) + ";" + s.InitCode + "})()")
					if err != nil {
						logs.Error(err)
					} else if !newValue.IsUndefined() {
						newValueInterface, err := newValue.Export()
						if err == nil {
							s.Value = newValueInterface
						}
					}
				}
			}
		}
	}

	// 处理在values中的元素的init
	for k, v := range values {
		e, found := eMap[k]
		if !found {
			continue
		}

		// 初始化值
		superElement := e.Super()
		if len(superElement.InitCode) > 0 {
			newValue, err := this.vm.Run("(function() { var value = " + stringutil.JSONEncode(v) + "; var values = " + stringutil.JSONEncode(values) + ";" + superElement.InitCode + "})()")
			if err != nil {
				logs.Error(err)
			} else if !newValue.IsUndefined() {
				newValueInterface, err := newValue.Export()
				if err == nil {
					v = newValueInterface
				}
			}
		}

		superElement.Value = v
	}
}

func (this *Form) ApplyRequest(req *http.Request) (values map[string]interface{}, errField string, err error) {
	values = map[string]interface{}{}
	for _, g := range this.Groups {
		for _, e := range g.Elements {
			superElement := e.Super()

			v, skip, err := e.ApplyRequest(req)
			if skip {
				continue
			}
			if err != nil {
				return values, "", err
			}

			// 校验字段值
			if len(superElement.ValidateCode) > 0 {
				newValue, err := this.vm.Run("(function() { var value = " + stringutil.JSONEncode(v) + ";" + superElement.ValidateCode + "})()")
				if err != nil {
					return values, superElement.Namespace + "_" + superElement.Code, errors.New(strings.Replace(err.Error(), "Error:", "", -1))
				}
				if !newValue.IsUndefined() {
					newValueInterface, err := newValue.Export()
					if err == nil {
						v = newValueInterface
					}
				}
			}

			values[superElement.Code] = v
		}
	}

	// 全局校验
	if len(this.ValidateCode) > 0 {
		newValues, err := this.vm.Run("(function() { var values = " + stringutil.JSONEncode(values) + ";" + this.ValidateCode + "})()")
		if err != nil {
			return values, "", err
		}
		if !newValues.IsUndefined() {
			newValueInterface, err := newValues.Export()
			if err == nil {
				m := map[string]interface{}{}
				if reflect.TypeOf(newValueInterface).Kind() == reflect.Map {
					mm := maps.NewMap(newValueInterface)
					if mm.Has("error") {
						field := mm.GetString("field")
						errorString := mm.GetString("error")
						return values, this.Namespace + "_" + field, errors.New(errorString)
					}
					m = mm
				}
				values = m
			}
		}
	}

	return values, "", nil
}

func (this *Form) Encode() ([]byte, error) {
	for _, g := range this.Groups {
		for _, e := range g.Elements {
			v := reflect.ValueOf(e).Elem().FieldByName("Element")
			if !v.IsValid() {
				logs.Error(errors.New("element should has 'Element' field"))
				continue
			}
			element := v.Interface().(Element)

			element.ClassType = reflect.TypeOf(e).Elem().Name()

			// namespace
			element.Namespace = this.Namespace
			v.Set(reflect.ValueOf(element))
		}
	}

	data, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (this *Form) EncodePretty() ([]byte, error) {
	data, err := this.Encode()
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	data, err = json.MarshalIndent(m, "", "   ")
	if err != nil {
		return nil, err
	}
	return data, nil
}
