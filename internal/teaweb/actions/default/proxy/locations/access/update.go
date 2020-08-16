package access

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction actions.Action

// 修改访问控制
func (this *UpdateAction) Run(params struct {
	ServerId   string
	LocationId string
}) {
	_, location := locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "access")

	if location.AccessPolicy != nil {
		this.Data["policy"] = location.AccessPolicy
	} else {
		this.Data["policy"] = shared.NewAccessPolicy()
	}

	this.Show()
}

// 保存访问控制
func (this *UpdateAction) RunPost(params struct {
	ServerId   string
	LocationId string

	TrafficOn         bool
	TrafficTotalOn    bool
	TrafficTotalTotal int64

	TrafficSecondOn       bool
	TrafficSecondTotal    int64
	TrafficSecondDuration int64

	TrafficMinuteOn       bool
	TrafficMinuteTotal    int64
	TrafficMinuteDuration int64

	TrafficHourOn       bool
	TrafficHourTotal    int64
	TrafficHourDuration int64

	TrafficDayOn       bool
	TrafficDayTotal    int64
	TrafficDayDuration int64

	TrafficMonthOn       bool
	TrafficMonthTotal    int64
	TrafficMonthDuration int64

	AccessOn            bool
	AccessAllowOn       bool
	AccessAllowIPValues []string

	AccessDenyOn       bool
	AccessDenyIPValues []string

	Must *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	location := server.FindLocation(params.LocationId)
	if location == nil {
		this.Fail("找不到Location")
	}
	policy := location.AccessPolicy
	if policy == nil {
		policy = shared.NewAccessPolicy()
		location.AccessPolicy = policy
	}

	// 控制信息
	policy.Traffic.On = params.TrafficOn
	policy.Traffic.Total.On = params.TrafficTotalOn
	policy.Traffic.Total.Total = params.TrafficTotalTotal

	policy.Traffic.Second.On = params.TrafficSecondOn
	policy.Traffic.Second.Total = params.TrafficSecondTotal
	policy.Traffic.Second.Duration = params.TrafficSecondDuration

	policy.Traffic.Minute.On = params.TrafficMinuteOn
	policy.Traffic.Minute.Total = params.TrafficMinuteTotal
	policy.Traffic.Minute.Duration = params.TrafficMinuteDuration

	policy.Traffic.Hour.On = params.TrafficHourOn
	policy.Traffic.Hour.Total = params.TrafficHourTotal
	policy.Traffic.Hour.Duration = params.TrafficHourDuration

	policy.Traffic.Day.On = params.TrafficDayOn
	policy.Traffic.Day.Total = params.TrafficDayTotal
	policy.Traffic.Day.Duration = params.TrafficDayDuration

	policy.Traffic.Month.On = params.TrafficMonthOn
	policy.Traffic.Month.Total = params.TrafficMonthTotal
	policy.Traffic.Month.Duration = params.TrafficMonthDuration

	policy.Access.On = params.AccessOn
	policy.Access.AllowOn = params.AccessAllowOn
	policy.Access.Allow = []*shared.ClientConfig{}
	for _, ip := range params.AccessAllowIPValues {
		if len(ip) > 0 {
			client := shared.NewClientConfig()
			client.IP = ip
			policy.Access.AddAllow(client)
		}
	}

	policy.Access.DenyOn = params.AccessDenyOn
	policy.Access.Deny = []*shared.ClientConfig{}
	for _, ip := range params.AccessDenyIPValues {
		if len(ip) > 0 {
			client := shared.NewClientConfig()
			client.IP = ip
			policy.Access.AddDeny(client)
		}
	}

	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success("保存成功")
}
