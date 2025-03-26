package vck_config

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gfsnotify"
)

func WatchVckConfig() {
	var (
		path     = gfile.Pwd() + "/manifest/config"
		ctx      = context.Background()
		logger   = g.Log()
		callback = func(event *gfsnotify.Event) {

			// if event.IsCreate() {
			// 	logger.Debug(ctx, "创建文件 : ", event.Path)
			// }
			if event.IsWrite() {
				if gfile.Basename(event.Path) == "vck.yaml" {
					logger.Debug(ctx, "Vck配置变化 : ", event.Path)
					// putAdminAtEtcd()

				}
				putAdminAtEtcd()

			}
			if event.IsRemove() {
				logger.Debug(ctx, "删除文件 : ", event.Path)
				logger.Debug(ctx, "Vck配置变化 : ", event.Path)
				putAdminAtEtcd()

			}
			// if event.IsRename() {
			// 	logger.Debug(ctx, "重命名文件 : ", event.Path)
			// }
			// if event.IsChmod() {
			// 	logger.Debug(ctx, "修改权限 : ", event.Path)
			// }
		}
	)
	g.Log().Debug(ctx, "开始监听配置文件", path)
	_, err := gfsnotify.Add(path, callback, gfsnotify.WatchOption{})
	if err != nil {
		logger.Fatal(ctx, err)
	} else {
		select {}
	}
}
