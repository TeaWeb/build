package teageo

// 纠正名称
func ConvertName(name string) string {
	switch name {
	case "台湾":
		name = "中国台湾"
	case "香港":
		name = "中国香港"
	case "澳门":
		name = "中国澳门"
	case "闽":
		name = "福建省"
	case "河南":
		name = "河南省"
	case "重庆":
		name = "重庆市"
	case "安徽":
		name = "安徽省"
	case "上海":
		name = "上海市"
	case "辽宁":
		name = "辽宁省"
	case "贵州":
		name = "贵州省"
	case "湖南":
		name = "湖南省"
	case "海南":
		name = "海南省"
	case "江西":
		name = "江西省"
	case "广东":
		name = "广东省"
	case "北京":
		name = "北京市"
	case "陕西":
		name = "陕西省"
	}
	return name
}
