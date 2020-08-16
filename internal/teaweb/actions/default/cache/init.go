package cache

import (
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/cache").
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Get("", new(IndexAction)).
			GetPost("/createPolicy", new(CreatePolicyAction)).
			Post("/deletePolicy", new(DeletePolicyAction)).
			GetPost("/updatePolicy", new(UpdatePolicyAction)).
			GetPost("/testPolicy", new(TestPolicyAction)).
			GetPost("/statPolicy", new(StatPolicyAction)).
			Get("/policy", new(PolicyAction)).
			GetPost("/cleanPolicy", new(CleanPolicyAction)).
			GetPost("/refreshPolicy", new(RefreshPolicyAction)).
			EndAll()
	})
}
