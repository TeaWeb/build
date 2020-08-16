package teacluster

import (
	"github.com/TeaWeb/build/internal/teacluster/configs"
	"github.com/TeaWeb/build/internal/teaevents"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
)

// cluster -> node
type SyncAction struct {
	Action

	ItemActions []*configs.ItemAction
}

func (this *SyncAction) Name() string {
	return "sync"
}

func (this *SyncAction) Execute() error {
	for _, itemAction := range this.ItemActions {
		logs.Println("[cluster]"+itemAction.Action, "'"+itemAction.ItemId+"'")
		switch itemAction.Action {
		case configs.ItemActionAdd:
			fallthrough
		case configs.ItemActionChange:
			file := files.NewFile(Tea.ConfigFile(itemAction.ItemId))
			dir := file.Parent()
			if !dir.Exists() {
				err := dir.MkdirAll()
				if err != nil {
					logs.Error(err)
					return err
				}
			}
			err := file.Write(itemAction.Item.Data)
			if err != nil {
				logs.Error(err)
			}
		case configs.ItemActionRemove:
			file := files.NewFile(Tea.ConfigFile(itemAction.ItemId))
			if file.Exists() {
				err := file.Delete()
				if err != nil {
					logs.Error(err)
				}
			}
		}
	}

	sumData := SharedManager.BuildSum()

	// write to local file
	file := files.NewFile(Tea.ConfigFile("cluster.sum"))
	err := file.Write(sumData)
	if err != nil {
		logs.Error(err)
	}

	// reload system
	teaevents.Post(teaevents.NewReloadEvent())

	return nil
}

func (this *SyncAction) TypeId() int8 {
	return ActionCodeSync
}
