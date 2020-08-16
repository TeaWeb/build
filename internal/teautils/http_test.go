package teautils

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
	"time"
)

func TestNewHTTPClient(t *testing.T) {
	a := assert.NewAssertion(t)

	client := NewHTTPClient(1 * time.Second)
	a.IsTrue(client.Timeout == 1*time.Second)

	client2 := NewHTTPClient(1 * time.Second)
	a.IsTrue(client != client2)
}

func TestSharedHTTPClient(t *testing.T) {
	a := assert.NewAssertion(t)

	_ = SharedHttpClient(2 * time.Second)
	_ = SharedHttpClient(3 * time.Second)

	client := SharedHttpClient(1 * time.Second)
	a.IsTrue(client.Timeout == 1*time.Second)

	client2 := SharedHttpClient(1 * time.Second)
	a.IsTrue(client == client2)

	t.Log(timeoutClientMap)
}
