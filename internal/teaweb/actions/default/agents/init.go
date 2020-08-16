package agents

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAgent,
			}).
			Helper(new(Helper)).
			Prefix("/agents").
			Get("", new(IndexAction)).
			GetPost("/addAgent", new(AddAgentAction)).
			GetPost("/delete", new(DeleteAction)).
			Post("/deleteAgents", new(DeleteAgentsAction)).
			Post("/move", new(MoveAction)).
			Get("/menu", new(MenuAction)).
			GetPost("/form", new(FormAction)).
			EndAll()
	})
}
