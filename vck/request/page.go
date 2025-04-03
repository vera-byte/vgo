package vck_request

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/grpc/metadata"
)

type PageReq struct {
	Order string `json:"order" sm:"排序" dc:"排序"`
	Page  int    `json:"page" sm:"页码" dc:"页码"`
	Size  int    `json:"size" sm:"页数" dc:"页数"`
	Sort  string `json:"sort" sm:"排序" dc:"排序"`
}

type Pagination struct {
	Page  int `json:"page" sm:"页码" dc:"页码"`
	Size  int `json:"size" sm:"页数" dc:"页数"`
	Total int `json:"total" sm:"总数" dc:"总数"`
}

// 在grpc中获取传入的分页对象
func GetGrpcPageRequest(ctx context.Context) *PageReq {
	// ✅ 读取 metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}

	// ✅ 获取 Admin 字段
	pageReq := md.Get("PageReq")
	if len(pageReq) == 0 {
		g.Log().Info(ctx, "Page")
		return nil
	}
	page := &PageReq{}
	gconv.Scan(pageReq[0], page)
	return page
}
