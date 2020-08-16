package teaproxy

// 清理Path中的多余的字符
func CleanPath(path string) string {
	l := len(path)
	if l == 0 {
		return "/"
	}
	result := []byte{'/'}
	isSlash := true
	for i := 0; i < l; i++ {
		if path[i] == '\\' || path[i] == '/' {
			if !isSlash {
				isSlash = true
				result = append(result, '/')
			}
		} else {
			isSlash = false
			result = append(result, path[i])
		}
	}
	return string(result)
}
