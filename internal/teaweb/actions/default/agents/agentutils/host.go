package agentutils

import (
	"fmt"
	"github.com/iwind/TeaGo/types"
	"net"
	"regexp"
	"strings"
	"time"
)

// 分析主机规则
func ParseHostRules(rules string, max int) (result []string) {
	rules = strings.TrimSpace(rules)
	if len(rules) == 0 {
		return
	}
	lines := regexp.MustCompile("[\n\r]").Split(rules, -1)
	rangeReg := regexp.MustCompile("\\[\\s*(\\d+)\\s*-\\s*(\\d+)\\s*]")
	leadingZeroReg := regexp.MustCompile("^0+")
	ipReg := regexp.MustCompile("^\\d+\\.\\d+\\.\\d+\\.\\d+$")
	ipLeadingZeroReg := regexp.MustCompile("\\b0+")

	count := 0
	for _, line := range lines {
		if max > 0 && count >= max {
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if rangeReg.MatchString(line) {
			lineArray := []string{line}
			for {
				newLineArray := []string{}
				for _, line1 := range lineArray {
					matches := rangeReg.FindStringSubmatch(line1)
					if len(matches) == 0 {
						break
					}
					fromString := matches[1]
					toString := matches[2]
					from := types.Int(fromString)
					to := types.Int(toString)
					lineLen := 0
					if leadingZeroReg.MatchString(fromString) {
						lineLen = len(fromString)
					}
					if from > to {
						break
					}
					indexes := rangeReg.FindStringIndex(line1)

					shouldBreak := false
					for i := from; i <= to; i ++ {
						newLine := line1[:indexes[0]] + fmt.Sprintf("%0"+fmt.Sprintf("%d", lineLen)+"d", i) + line1[indexes[1]:]

						newLineArray = append(newLineArray, newLine)

						if !rangeReg.MatchString(newLine) {
							count ++
							if max > 0 && count >= max {
								shouldBreak = true
								break
							}
						}
					}

					if shouldBreak {
						break
					}
				}

				if len(newLineArray) > 0 {
					lineArray = newLineArray
				} else {
					break
				}
			}

			// 如果是IP去掉leading zeros
			for _, line2 := range lineArray {
				if ipReg.MatchString(line2) {
					line2 = ipLeadingZeroReg.ReplaceAllString(line2, "")
				}
				result = append(result, line2)
			}
		} else {
			result = append(result, line)
			count ++
		}
	}

	return
}

// 检查主机设置
func CheckHostConnectivity(host string, port int, timeout time.Duration) (cost time.Duration, b bool) {
	before := time.Now()
	defer func() {
		cost = time.Since(before)
	}()

	dialer := net.Dialer{
		Timeout: timeout,
	}
	conn, err := dialer.Dial("tcp", host+":"+fmt.Sprintf("%d", port))
	if err != nil {
		return
	}
	conn.Close()
	b = true
	return
}
