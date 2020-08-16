package main

import (
	_ "github.com/TeaWeb/build/internal/ext"
	"github.com/TeaWeb/build/internal/teaweb"
	_ "github.com/iwind/TeaGo/bootstrap"
)

func main() {
	teaweb.Start()
}
