package teaconfigs

import "testing"

func TestSharedWAFList(t *testing.T) {
	wafList := SharedWAFList()

	// add
	wafList.AddFile("abc")
	err := wafList.Save()
	if err != nil {
		t.Fatal(err)
	}

	// remove
	wafList.RemoveFile("abc")
	err = wafList.Save()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(wafList.Files)
}

func TestWAFList_FindAllConfigs(t *testing.T) {
	wafList := SharedWAFList()
	for _, config := range wafList.FindAllConfigs() {
		t.Log("Name:", config.Name, "Id:", config.Id)
	}
}
