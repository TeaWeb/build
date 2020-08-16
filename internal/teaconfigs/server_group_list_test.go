package teaconfigs

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestSharedServerGroupList(t *testing.T) {
	logs.PrintAsJSON(SharedServerGroupList(), t)
}
