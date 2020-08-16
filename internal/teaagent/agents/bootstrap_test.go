package teaagents

import "testing"

func TestOSName(t *testing.T) {
	t.Log(retrieveOSName())
	t.Log(retrieveOSNameBase64())

	// from cache
	t.Log(retrieveOSName())
	t.Log(retrieveOSNameBase64())
}
