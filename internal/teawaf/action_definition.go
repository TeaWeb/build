package teawaf

import "reflect"

// action definition
type ActionDefinition struct {
	Name        string
	Code        ActionString
	Description string
	Instance    ActionInterface
	Type        reflect.Type
}
