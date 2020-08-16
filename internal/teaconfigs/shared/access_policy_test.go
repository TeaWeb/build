package shared

import (
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/utils/time"
	"log"
	"sync"
	"testing"
	"time"
)

func TestAPIAccessPolicy(t *testing.T) {
	printTime := func(t1 int64) {
		log.Println(timeutil.Format("Y-m-d H:i:s", time.Unix(t1, 0)))
	}

	printTime2 := func(t1 int64, t2 int64) {
		log.Println(timeutil.Format("Y-m-d H:i:s", time.Unix(t1, 0)), timeutil.Format("Y-m-d H:i:s", time.Unix(t2, 0)))
	}

	p := AccessPolicy{}
	p.Traffic.On = true
	p.Traffic.Second.On = true
	p.Traffic.Second.Duration = 1

	p.Traffic.Minute.On = true
	p.Traffic.Minute.Duration = 1

	p.Traffic.Hour.On = true
	p.Traffic.Hour.Duration = 1

	p.Traffic.Day.On = true
	p.Traffic.Day.Duration = 1

	p.Traffic.Month.On = true
	p.Traffic.Month.Duration = 1

	for {
		time.Sleep(1000 * time.Millisecond)

		log.Println("===")
		p.IncreaseTraffic()

		log.Println("seconds:", p.Traffic.Second.Used)
		printTime(p.Traffic.Second.FromTime)

		log.Println("minutes:", p.Traffic.Minute.Used)
		printTime2(p.Traffic.Minute.FromTime, p.Traffic.Minute.ToTime)

		/**log.Println("hours:", p.Traffic.Hour.Used)
		printTime2(p.Traffic.Hour.fromTime, p.Traffic.Hour.toTime)

		log.Println("days:", p.Traffic.Minute.Used)
		printTime2(p.Traffic.Day.fromTime, p.Traffic.Day.toTime)

		log.Println("months:", p.Traffic.Minute.Used)
		printTime2(p.Traffic.Month.fromTime, p.Traffic.Month.toTime)**/

		break
	}
}

func TestAPIAccessPolicySecond(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	p := AccessPolicy{}
	p.Traffic.On = true
	p.Traffic.Second.On = false
	{
		_, isAllowed := p.AllowTraffic()
		a.IsTrue(isAllowed)
	}

	p.Traffic.Second.On = true
	{
		_, isAllowed := p.AllowTraffic()
		a.IsFalse(isAllowed)
	}

	p.Traffic.Second.On = true
	p.Traffic.Second.Duration = 1
	p.Traffic.Second.Total = 1
	{
		_, isAllowed := p.AllowTraffic()
		a.IsTrue(isAllowed)
	}

	p.Traffic.Second.On = true
	p.Traffic.Second.Duration = 1
	p.Traffic.Second.Total = 2
	p.IncreaseTraffic()
	p.IncreaseTraffic()
	{
		_, isAllowed := p.AllowTraffic()
		a.IsFalse(isAllowed)
	}

	time.Sleep(1 * time.Second)
	{
		_, isAllowed := p.AllowTraffic()
		a.IsTrue(isAllowed)
	}
}

func TestAPIAccessPolicyMinute(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	p := AccessPolicy{}
	p.Traffic.On = true
	p.Traffic.Minute.On = true
	{
		_, isAllowed := p.AllowTraffic()
		a.IsFalse(isAllowed)
	}

	p.Traffic.Minute.Total = 1
	{
		_, isAllowed := p.AllowTraffic()
		a.IsFalse(isAllowed)
	}

	p.Traffic.Minute.Duration = 1
	{
		_, isAllowed := p.AllowTraffic()
		a.IsTrue(isAllowed)
	}

	p.IncreaseTraffic()
	{
		_, isAllowed := p.AllowTraffic()
		a.IsFalse(isAllowed)
	}

	//time.Sleep(61 * time.Second)
	//a.IsTrue(p.AllowTraffic())
}

func BenchmarkAccessPolicyPerformance(b *testing.B) {
	p := AccessPolicy{}

	for i := 0; i < b.N; i++ {
		locker := sync.Mutex{}
		locker.Lock()

		p.Traffic.On = true
		p.Traffic.Second.On = true
		p.Traffic.Second.Duration = 1
		p.Traffic.Second.Total = 1
		p.Traffic.Minute.On = true
		p.Traffic.Minute.Duration = 1
		p.Traffic.Minute.Total = 1
		p.Traffic.Hour.On = true
		p.Traffic.Hour.Duration = 1
		p.Traffic.Hour.Total = 1
		p.Traffic.Day.On = true
		p.Traffic.Day.Duration = 1
		p.Traffic.Day.Total = 1
		p.Traffic.Month.On = true
		p.Traffic.Month.Duration = 1
		p.Traffic.Month.Total = 1
		p.AllowTraffic()
		p.IncreaseTraffic()

		locker.Unlock()
	}
}

func TestAccessPolicy_AllowAccess(t *testing.T) {
	p := AccessPolicy{}
	p.Access.On = true

	{
		client := NewClientConfig()
		client.IP = "192.168.1.100"
		client.On = false
		p.Access.AllowOn = false
		p.Access.AddAllow(client)
	}

	{
		client := NewClientConfig()
		client.IP = "192.168.1.101"
		client.On = true
		p.Access.DenyOn = true
		p.Access.AddDeny(client)
	}

	p.Validate()
	t.Log(p.AllowAccess("192.168.1.100"))
}
