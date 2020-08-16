package teaagents

import "github.com/TeaWeb/build/internal/teaconfigs/agents"

type App struct {
	config *agents.AppConfig
}

func NewApp(config *agents.AppConfig) *App {
	return &App{
		config: config,
	}
}
