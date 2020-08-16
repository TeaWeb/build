package checkpoints

import (
	"fmt"
	"strings"
	"testing"
)

func TestFindCheckpointDefinition_Markdown(t *testing.T) {
	result := []string{}
	for _, def := range AllCheckpoints {
		row := "## " + def.Name + "\n* 前缀：`${" + def.Prefix + "}`\n* 描述：" + def.Description
		if def.HasParams {
			row += "\n* 是否有子参数：YES"

			paramOptions := def.Instance.ParamOptions()
			if paramOptions != nil && len(paramOptions.Options) > 0 {
				row += "\n* 可选子参数"
				for _, option := range paramOptions.Options {
					row += "\n   * `" + option.Name + "`：值为 `" + option.Value + "`"
				}
			}
		} else {
			row += "\n* 是否有子参数：NO"
		}
		row += "\n"
		result = append(result, row)
	}

	fmt.Print(strings.Join(result, "\n") + "\n")
}
