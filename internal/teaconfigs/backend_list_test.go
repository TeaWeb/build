package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestBackendList_CloneBackendList(t *testing.T) {
	a := assert.NewAssertion(t)

	backendList := new(BackendList)
	backendList.AddBackend(&BackendConfig{
		Id: "1",
	})
	backendList.AddBackend(&BackendConfig{
		Id: "2",
	})
	backendList.AddBackend(&BackendConfig{
		Id: "3",
	})
	clone := backendList.CloneBackendList()
	a.IsTrue(len(clone.Backends) == 3)
	clone.DeleteBackend("2")
	a.IsTrue(len(clone.Backends) == 2)
	t.Log(len(backendList.Backends))
	a.IsTrue(len(backendList.Backends) == 3)
}
