package notices

import (
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/utils/time"
	"testing"
	"time"
)

func TestNoticeMediaConfig_ShouldNotify(t *testing.T) {
	a := assert.NewAssertion(t)

	media := NewNoticeMediaConfig()
	media.On = false
	a.IsFalse(media.ShouldNotify(0))

	media.On = true
	a.IsTrue(media.ShouldNotify(0))

	media.TimeFrom = "00:00:00"
	media.TimeTo = timeutil.Format("H:i:s", time.Now().Add(-1*time.Second))
	a.IsFalse(media.ShouldNotify(0))

	media.TimeFrom = "00:00:00"
	media.TimeTo = "23:59:59"
	a.IsTrue(media.ShouldNotify(0))

	media.TimeFrom = "00:00:00"
	media.TimeTo = "23:59:59"
	media.RateCount = 5
	media.RateMinutes = 1
	a.IsFalse(media.ShouldNotify(10))
	a.IsTrue(media.ShouldNotify(4))
}
