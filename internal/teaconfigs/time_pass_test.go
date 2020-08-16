package teaconfigs

import (
	timeutil "github.com/iwind/TeaGo/utils/time"
	"testing"
	"time"
)

func TestTimePastSeconds(t *testing.T) {
	t.Log("current:", timeutil.Format("Y-m-d H:i:s"))
	for _, m := range AllTimePasts() {
		value := m["value"].(TimePast)
		t.Log(value+":", TimePastUnixTime(value), timeutil.Format("Y-m-d H:i:s", time.Unix(TimePastUnixTime(value), 0)))
	}
}
