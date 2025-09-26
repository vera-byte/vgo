package openapi

import (
	_ "github.com/vera-byte/vgo/modules/openapi/controller/admin"
	_ "github.com/vera-byte/vgo/modules/openapi/controller/open"
	_ "github.com/vera-byte/vgo/modules/openapi/middleware"
)

// // init 模块初始化函数
// // 功能: 初始化开放平台模块，注册配置和控制器
// func init() {
// 	var (
// 		ctx = gctx.GetInitCtx()
// 	)
// 	g.Log().Debug(ctx, "module openapi init start ...")

// 	// 初始化数据模型
// 	// v.FillInitData(ctx, "openapi", &model.OpenapiApp{})
// 	// v.FillInitData(ctx, "openapi", &model.OpenapiSignLog{})

// 	g.Log().Debug(ctx, "module openapi init finished ...")
// }
