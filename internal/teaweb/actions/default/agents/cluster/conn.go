package cluster

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net"
	"regexp"
	"sync"
	"time"
)

type ConnAction actions.Action

// 测试连接
func (this *ConnAction) Run(params struct {
	Hosts []string
	Port  int
}) {
	this.Data["states"] = []interface{}{}

	if len(params.Hosts) == 0 {
		this.Success()
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(params.Hosts))
	states := []maps.Map{}
	statesLocker := sync.Mutex{}
	for _, h := range params.Hosts {
		go func(host string) {
			defer wg.Done()
			cost, b := agentutils.CheckHostConnectivity(host, params.Port, 3*time.Second)

			ip := ""
			name := ""
			if b {
				if regexp.MustCompile("^(\\d+)\\.(\\d+)\\.(\\d+)\\.(\\d+)$").MatchString(host) {
					ip = host
				} else {
					name = host

					ipList, err := net.LookupIP(host)
					if err == nil && len(ipList) > 0 {
						ip = ipList[0].String()
					}
				}
			}

			statesLocker.Lock()
			states = append(states, maps.Map{
				"addr":       host,
				"cost":       cost.Seconds() * 1000,
				"canConnect": b,
				"ip":         ip,
				"name":       name,
			})
			statesLocker.Unlock()
		}(h)
	}
	wg.Wait()

	this.Data["states"] = states

	this.Success()
}
