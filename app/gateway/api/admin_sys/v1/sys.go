package v1

import "github.com/gogf/gf/v2/frame/g"

type LogPageReq struct {
	g.Meta `path:"page" method:"post" sm:"日志分页" tags:"系统"`
}
type LogPageRes struct {
}
