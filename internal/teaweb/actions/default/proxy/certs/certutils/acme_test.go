package certutils

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"testing"
)

func TestRenewACMECerts(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	RenewACMECerts()
}
