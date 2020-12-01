package teawaf

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/types"
	"regexp"
	"strings"
	"time"
)

type IPAction = string

const (
	IPActionAccept IPAction = "accept"
	IPActionReject IPAction = "reject"
)

// ip table
type IPTable struct {
	Id       string   `yaml:"id" json:"id"`
	On       bool     `yaml:"on" json:"on"`
	IP       string   `yaml:"ip" json:"ip"`             // single ip, cidr, ip range, TODO support *
	Port     string   `yaml:"port" json:"port"`         // single port, range, *
	Action   IPAction `yaml:"action" json:"action"`     // accept, reject
	TimeFrom int64    `yaml:"timeFrom" json:"timeFrom"` // from timestamp
	TimeTo   int64    `yaml:"timeTo" json:"timeTo"`     // zero means forever
	Remark   string   `yaml:"remark" json:"remark"`

	// port
	minPort int
	maxPort int

	minPortWildcard bool
	maxPortWildcard bool

	ports []int

	// ip
	ipRange *shared.IPRangeConfig
}

func NewIPTable() *IPTable {
	return &IPTable{
		On: true,
		Id: rands.HexString(16),
	}
}

func (this *IPTable) Init() error {
	// parse port
	if teautils.RegexpDigitNumber.MatchString(this.Port) {
		this.minPort = types.Int(this.Port)
		this.maxPort = types.Int(this.Port)
	} else if regexp.MustCompile(`[:-]`).MatchString(this.Port) {
		pieces := regexp.MustCompile(`[:-]`).Split(this.Port, 2)
		if pieces[0] == "*" {
			this.minPortWildcard = true
		} else {
			this.minPort = types.Int(pieces[0])
		}
		if pieces[1] == "*" {
			this.maxPortWildcard = true
		} else {
			this.maxPort = types.Int(pieces[1])
		}
	} else if strings.Contains(this.Port, ",") {
		pieces := strings.Split(this.Port, ",")
		for _, piece := range pieces {
			piece = strings.TrimSpace(piece)
			if len(piece) > 0 {
				this.ports = append(this.ports, types.Int(piece))
			}
		}
	} else if this.Port == "*" {
		this.minPortWildcard = true
		this.maxPortWildcard = true
	}

	// parse ip
	if len(this.IP) > 0 {
		ipRange, err := shared.ParseIPRange(this.IP)
		if err != nil {
			return err
		}
		this.ipRange = ipRange
	}

	return nil
}

// check ip
func (this *IPTable) Match(ip string, port int) (isMatched bool) {
	if !this.On {
		return
	}

	now := time.Now().Unix()
	if this.TimeFrom > 0 && now < this.TimeFrom {
		return
	}
	if this.TimeTo > 0 && now > this.TimeTo {
		return
	}

	if !this.matchPort(port) {
		return
	}

	if !this.matchIP(ip) {
		return
	}

	return true
}

func (this *IPTable) matchPort(port int) bool {
	if port == 0 {
		return false
	}
	if this.minPortWildcard {
		if this.maxPortWildcard {
			return true
		}
		if this.maxPort >= port {
			return true
		}
	}
	if this.maxPortWildcard {
		if this.minPortWildcard {
			return true
		}
		if this.minPort <= port {
			return true
		}
	}
	if (this.minPort > 0 || this.maxPort > 0) && this.minPort <= port && this.maxPort >= port {
		return true
	}
	if len(this.ports) > 0 {
		return lists.ContainsInt(this.ports, port)
	}
	return false
}

func (this *IPTable) matchIP(ip string) bool {
	if this.ipRange == nil {
		return false
	}
	return this.ipRange.Contains(ip)
}
