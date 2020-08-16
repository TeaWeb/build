package teacluster

import "testing"

func TestRegisterActionType(t *testing.T) {
	a := FindActionInstance(1)
	t.Log(a)
	t.Log(a.BaseAction())
	a.BaseAction().RequestId = 123
	t.Log(a)
}

func TestGenerateActionId(t *testing.T) {
	t.Log(GenerateActionId())
	t.Log(GenerateActionId())
}
