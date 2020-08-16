package teautils

import (
	"fmt"
	"strconv"
)

// 命令帮助
type CommandHelp struct {
	product       string
	version       string
	usage         string
	options       []*CommandHelpOption
	appendStrings []string
}

func NewCommandHelp() *CommandHelp {
	return &CommandHelp{}
}

type CommandHelpOption struct {
	Code        string
	Description string
}

// 产品
func (this *CommandHelp) Product(product string) *CommandHelp {
	this.product = product
	return this
}

// 版本
func (this *CommandHelp) Version(version string) *CommandHelp {
	this.version = version
	return this
}

// 使用方法
func (this *CommandHelp) Usage(usage string) *CommandHelp {
	this.usage = usage
	return this
}

// 选项
func (this *CommandHelp) Option(code string, description string) *CommandHelp {
	this.options = append(this.options, &CommandHelpOption{
		Code:        code,
		Description: description,
	})
	return this
}

// 附加内容
func (this *CommandHelp) Append(appendString string) *CommandHelp {
	this.appendStrings = append(this.appendStrings, appendString)
	return this
}

// 打印
func (this *CommandHelp) Print() {
	fmt.Println(this.product + " v" + this.version)
	fmt.Println("Usage:", "\n   "+this.usage)

	if len(this.options) > 0 {
		fmt.Println("")
		fmt.Println("Options:")

		spaces := 20
		max := 40
		for _, option := range this.options {
			l := len(option.Code)
			if l < max && l > spaces {
				spaces = l + 4
			}
		}

		for _, option := range this.options {
			if len(option.Code) > max {
				fmt.Println("")
				fmt.Println("  " + option.Code)
				option.Code = ""
			}

			fmt.Printf("  %-"+strconv.Itoa(spaces)+"s%s\n", option.Code, ": "+option.Description)
		}
	}

	if len(this.appendStrings) > 0 {
		fmt.Println("")
		for _, s := range this.appendStrings {
			fmt.Println(s)
		}
	}
}
