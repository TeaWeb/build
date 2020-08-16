package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

// Check Point
type CheckpointInterface interface {
	// initialize
	Init()

	// is request?
	IsRequest() bool

	// get request value
	RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error)

	// get response value
	ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error)

	// param option list
	ParamOptions() *ParamOptions

	// options
	Options() []OptionInterface

	// start
	Start()

	// stop
	Stop()
}
