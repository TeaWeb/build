package teaagents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"testing"
)

func TestApp_Run(t *testing.T) {
	config := agents.NewAppConfig()

	app := NewApp(config)
	t.Log(app)
}
