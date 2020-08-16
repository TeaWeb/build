package locations

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Prefix("/proxy/locations").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(proxy.Helper)).
			Get("", new(IndexAction)).
			GetPost("/add", new(AddAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/delete", new(DeleteAction)).
			Post("/moveUp", new(MoveUpAction)).
			Post("/moveDown", new(MoveDownAction)).
			Post("/move", new(MoveAction)).
			Get("/detail", new(DetailAction)).
			Get("/headers", new(HeadersAction)).
			Get("/rewrite", new(RewriteAction)).
			Get("/fastcgi", new(FastcgiAction)).
			GetPost("/access", new(AccessAction)).
			Get("/cache", new(CacheAction)).
			Post("/updateCache", new(UpdateCacheAction)).
			Get("/waf", new(WafAction)).
			Post("/waf/update", new(WafUpdateAction)).
			Post("/test", new(TestAction)).
			Get("/export", new(ExportAction)).
			Post("/duplicate", new(DuplicateAction)).
			GetPost("/import", new(ImportAction)).
			EndAll()
	})
}
