package teawaf

import (
	"fmt"
	"strings"
	"testing"
)

func TestRuleOperator_Markdown(t *testing.T) {
	result := []string{}
	for _, def := range AllRuleOperators {
		row := "## " + def.Name + "\n"
		row += "符号：`" + def.Code + "`\n"
		row += "描述：" + def.Description + "\n"
		result = append(result, row)
	}

	fmt.Print(strings.Join(result, "\n") + "\n")
}
