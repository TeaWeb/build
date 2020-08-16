package apiutils

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/pquerna/ffjson/ffjson"
)

// 校验用户
func ValidateUser(actionPtr actions.ActionWrapper) {
	action := actionPtr.Object()
	action.AddHeader("Content-Type", "application/json; charset=utf-8")

	key, found := action.Param("TeaKey")
	if !found || len(key) == 0 {
		Fail(actionPtr, "Authenticate Failed 001")
	}

	user := configs.SharedAdminConfig().FindUserWithKey(key)
	if user == nil {
		Fail(actionPtr, "Authenticate Failed 002")
	}
}

// 错误提示
func Fail(actionPtr actions.ActionWrapper, message string) {
	action := actionPtr.Object()
	action.ResponseWriter.WriteHeader(400)
	action.Fail(message)
}

// 成功并返回数据
func Success(actionPtr actions.ActionWrapper, data interface{}) {
	if actionPtr.Object().HasParam("TeaPretty") {
		dataBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			Fail(actionPtr, err.Error())
			return
		}
		actionPtr.Object().Write(dataBytes)
	} else {
		dataBytes, err := ffjson.Marshal(data)
		if err != nil {
			Fail(actionPtr, err.Error())
			return
		}
		actionPtr.Object().Write(dataBytes)
	}
}

// 成功
func SuccessOK(actionPtr actions.ActionWrapper) {
	actionPtr.Object().WriteJSON(maps.Map{
		"ok": 1,
	})
}
