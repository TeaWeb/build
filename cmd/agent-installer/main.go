package main

import (
	"flag"
	"fmt"
	"github.com/TeaWeb/agentinstaller/pkg/installers"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
)

// agentinstaller -id=ID -key=KEY -master=http://xxxx:7777 -dir=xxx
func main() {
	var idArg string
	var keyArg string
	var masterArg string
	var dirArg string
	flag.StringVar(&idArg, "id", "", "ID")
	flag.StringVar(&keyArg, "key", "", "Key")
	flag.StringVar(&masterArg, "master", "", "Master Address")
	flag.StringVar(&dirArg, "dir", "", "Install dir")
	flag.Parse()

	installer := installers.NewInstaller()
	installer.Id = idArg
	installer.Key = keyArg
	installer.Master = masterArg
	installer.Dir = dirArg

	isInstalled, err := installer.Start()
	if err != nil {
		fmt.Print(stringutil.JSONEncode(maps.Map{
			"isInstalled": isInstalled,
			"err":         err.Error(),
			"ip":          installer.IP,
		}))
	} else {
		fmt.Print(stringutil.JSONEncode(maps.Map{
			"isInstalled": isInstalled,
			"err":         nil,
			"ip":          installer.IP,
		}))
	}
}
