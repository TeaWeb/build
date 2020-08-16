package teacluster

type ActionCode = int8

const (
	ActionCodeSuccess  ActionCode = 1
	ActionCodeFail     ActionCode = 2
	ActionCodeRegister ActionCode = 3
	ActionCodeNotify   ActionCode = 4
	ActionCodePush     ActionCode = 5
	ActionCodePull     ActionCode = 6
	ActionCodePing     ActionCode = 7
	ActionCodeSync     ActionCode = 8
	ActionCodeSum      ActionCode = 9
	ActionCodeRun      ActionCode = 10
)
