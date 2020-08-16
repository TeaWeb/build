package teaconfigs

import (
	"testing"
)

func TestNode(t *testing.T) {
	config := SharedNodeConfig()
	t.Log(config)

	config = SharedNodeConfig()
	t.Log(config)
}
