package plugins

import (
	"github.com/TeaWeb/build/internal/teaplugins"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
	pluginArray := []maps.Map{}
	for _, p := range teaplugins.Plugins() {
		if !p.IsExternal {
			continue
		}

		pluginArray = append(pluginArray, maps.Map{
			"name":        p.Name,
			"developer":   p.Developer,
			"version":     p.Version,
			"site":        p.Site,
			"date":        p.Date,
			"description": p.Description,
			"code":        p.Code,
		})
	}

	this.Data["plugins"] = pluginArray

	this.Show()
}
