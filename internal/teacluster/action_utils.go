package teacluster

import (
	"errors"
	"fmt"
	"github.com/iwind/TeaGo/logs"
	"github.com/vmihailenco/msgpack"
	"reflect"
	"sync/atomic"
)

var actionMap = map[int8]reflect.Type{}
var actionId uint64

func RegisterActionType(actions ...ActionInterface) {
	for _, action := range actions {
		typeId := action.TypeId()
		_, ok := actionMap[typeId]
		if ok {
			logs.Error(errors.New("action type '" + fmt.Sprintf("%d", typeId) + "' already exist"))
			continue
		}

		msgpack.RegisterExt(typeId, action)
		actionMap[typeId] = reflect.TypeOf(action).Elem()
	}
}

func FindActionInstance(typeId int8) ActionInterface {
	t, ok := actionMap[typeId]
	if !ok {
		return nil
	}
	return reflect.New(t).Interface().(ActionInterface)
}

func GenerateActionId() uint64 {
	return atomic.AddUint64(&actionId, 1)
}
