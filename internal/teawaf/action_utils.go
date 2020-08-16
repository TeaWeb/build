package teawaf

import (
	"github.com/iwind/TeaGo/maps"
	"reflect"
)

var AllActions = []*ActionDefinition{
	{
		Name:     "阻止",
		Code:     ActionBlock,
		Instance: new(BlockAction),
	},
	{
		Name:     "允许通过",
		Code:     ActionAllow,
		Instance: new(AllowAction),
	},
	{
		Name:     "允许并记录日志",
		Code:     ActionLog,
		Instance: new(LogAction),
	},
	{
		Name:     "Captcha验证码",
		Code:     ActionCaptcha,
		Instance: new(CaptchaAction),
	},
	{
		Name:     "跳到下一个规则分组",
		Code:     ActionGoGroup,
		Instance: new(GoGroupAction),
		Type:     reflect.TypeOf(new(GoGroupAction)).Elem(),
	},
	{
		Name:     "跳到下一个规则集",
		Code:     ActionGoSet,
		Instance: new(GoSetAction),
		Type:     reflect.TypeOf(new(GoSetAction)).Elem(),
	},
}

func FindActionInstance(action ActionString, options maps.Map) ActionInterface {
	for _, def := range AllActions {
		if def.Code == action {
			if def.Type != nil {
				// create new instance
				ptrValue := reflect.New(def.Type)
				instance := ptrValue.Interface().(ActionInterface)

				if len(options) > 0 {
					count := def.Type.NumField()
					for i := 0; i < count; i++ {
						field := def.Type.Field(i)
						tag, ok := field.Tag.Lookup("yaml")
						if ok {
							v, ok := options[tag]
							if ok && reflect.TypeOf(v) == field.Type {
								ptrValue.Elem().FieldByName(field.Name).Set(reflect.ValueOf(v))
							}
						}
					}
				}

				return instance
			}

			// return shared instance
			return def.Instance
		}
	}
	return nil
}

func FindActionName(action ActionString) string {
	for _, def := range AllActions {
		if def.Code == action {
			return def.Name
		}
	}
	return ""
}
