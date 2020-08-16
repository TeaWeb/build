package teatesting

import "testing"

func TestRequireServer(t *testing.T) {
	t.Log(RequireHTTPServer())
	t.Log(RequireHTTPServer())
	t.Log(RequireHTTPServer())
}
