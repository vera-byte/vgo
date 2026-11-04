package v1

import "github.com/gogf/gf/v2/frame/g"

// DemoSampleWelcomeReq 示例自定义接口请求
type DemoSampleWelcomeReq struct {
	g.Meta `path:"/welcome" method:"GET" summary:"示例接口" tags:"示例"`
}
