package mongo

import (
	"fmt"
	"github.com/iwind/TeaGo/actions"
	"time"
)

type InstallStatusAction actions.Action

func (this *InstallStatusAction) Run(params struct{}) {
	this.Data["status"] = installStatus
	this.Data["percent"] = installPercent

	timeLeft := ""
	seconds := int64(0)
	if installPercent > 0 && installPercent < 100 {
		seconds = time.Now().Unix() - installStartAt
		if seconds > 0 {
			leftSeconds := int(float64(100-installPercent) * float64(seconds) / float64(installPercent))
			if leftSeconds > 60 {
				timeLeft = fmt.Sprintf("%dm%ds", leftSeconds/60, leftSeconds%60)
			} else if leftSeconds > 0 {
				timeLeft = fmt.Sprintf("%ds", leftSeconds)
			}
		}
	}
	this.Data["timeLeft"] = timeLeft
	this.Data["timeSeconds"] = seconds

	this.Success()
}
