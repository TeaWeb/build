package teaconfigs

import "testing"

func TestServerList(t *testing.T) {
	list, err := SharedServerList()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list.Files)
}

func TestServerList_RemoveServer(t *testing.T) {
	list, err := SharedServerList()
	if err != nil {
		t.Fatal(err)
	}
	list.RemoveServer("server.www.proxy.conf")
	err = list.Save()
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerList_AddServer(t *testing.T) {
	list, err := SharedServerList()
	if err != nil {
		t.Fatal(err)
	}
	list.AddServer("server.www.proxy1.conf")
	err = list.Save()
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerList_FindAllServers(t *testing.T) {
	list, err := SharedServerList()
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range list.FindAllServers() {
		t.Log(s.Filename)
	}
}
