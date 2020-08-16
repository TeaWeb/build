package teacluster

import "testing"

func TestManager_Start(t *testing.T) {
	err := SharedManager.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestManager_PullItems(t *testing.T) {
	SharedManager.PullItems()
}
