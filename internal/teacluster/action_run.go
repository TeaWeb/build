package teacluster

import (
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type RunAction struct {
	Action

	Cmd  string
	Data maps.Map
}

func (this *RunAction) Name() string {
	return "run"
}

func (this *RunAction) Execute() error {
	switch this.Cmd {
	case "cache.refresh":
		this.runCacheRefresh()
	case "cache.clean":
		this.runCacheClean()
	}
	return nil
}

func (this *RunAction) TypeId() int8 {
	return ActionCodeRun
}

// clean cache with prefixes
func (this *RunAction) runCacheRefresh() {
	filename := this.Data.GetString("filename")
	policy := shared.NewCachePolicyFromFile(filename)
	if policy == nil {
		logs.Println("[cluster][cache.refresh]can not find policy with '" + filename + "'")
		return
	}

	manager := teacache.FindCachePolicyManager(filename)
	if manager == nil {
		manager = teacache.NewManagerFromConfig(policy)
		defer func() {
			_ = manager.Close()
		}()
	}
	prefixes := this.Data.GetSlice("prefixes")
	prefixStrings := []string{}
	for _, prefix := range prefixes {
		prefixStrings = append(prefixStrings, types.String(prefix))
	}
	_, err := manager.DeletePrefixes(prefixStrings)
	if err != nil {
		logs.Println("[cluster][cache.refresh]delete prefixes: " + err.Error())
		return
	}
}

// clean cache
func (this *RunAction) runCacheClean() {
	filename := this.Data.GetString("filename")
	policy := shared.NewCachePolicyFromFile(filename)
	if policy == nil {
		logs.Println("[cluster][cache.clean]can not find policy with '" + filename + "'")
		return
	}

	manager := teacache.FindCachePolicyManager(filename)
	if manager == nil {
		manager = teacache.NewManagerFromConfig(policy)
		defer func() {
			_ = manager.Close()
		}()
	}
	err := manager.Clean()
	if err != nil {
		logs.Println("[cluster][cache.clean]clean: " + err.Error())
		return
	}
}
