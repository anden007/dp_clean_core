package misc

import "strings"

// 处理Swagger模版，替换其中接口访问占位符为当前挂载路径
func ProcessSwaggerTemplate(template string, placeholder string, realPath string) (result string) {
	if strings.EqualFold(strings.TrimSpace(realPath), "/") {
		realPath = ""
	}
	return strings.ReplaceAll(template, placeholder, realPath)
}
