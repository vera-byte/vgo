package v

import (
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
)

type vFunc interface {
	// Func handler
	Func(ctx g.Ctx, param string) (err error)
	// IsSingleton 是否单例,当为true时，只能有一个任务在执行,在注意函数为计划任务时使用
	IsSingleton() bool
	// IsAllWorker 是否所有worker都执行
	IsAllWorker() bool
}

// FuncMap 函数列表
var FuncMap = make(map[string]vFunc)

// RegisterFunc 注册函数
func RegisterFunc(name string, f vFunc) {
	FuncMap[name] = f
}

// GetFunc 获取函数
func GetFunc(name string) vFunc {
	return FuncMap[name]
}

// RunFunc 运行函数
func RunFunc(ctx g.Ctx, funcstring string) (err error) {
	l := gstr.Pos(funcstring, "(")
	r := gstr.Pos(funcstring, ")")
	if l < 0 || r < l {
		return gerror.New("函数格式错误:" + funcstring)
	}
	funcName := gstr.SubStr(funcstring, 0, l)
	funcParam := gstr.SubStr(funcstring, l+1, r-l-1)
	if _, ok := FuncMap[funcName]; !ok {
		err = gerror.New("函数不存在:" + funcName)
		return
	}
	if !FuncMap[funcName].IsAllWorker() {
		// 检查当前是否为主进程, 如果不是主进程, 则不执行
		if ProcessFlag != CacheManager.MustGetOrSet(ctx, "v:masterflag", ProcessFlag, 60*time.Second).String() {
			g.Log().Debug(ctx, "当前进程不是主进程, 不执行单例函数", funcName)
			return
		}
	}
	err = FuncMap[funcName].Func(ctx, funcParam)
	return
}

// ClusterRunFunc 集群运行函数,如果是单机模式, 则直接运行函数
func ClusterRunFunc(ctx g.Ctx, funcstring string) (err error) {
	if IsRedisMode {
		conn, err := g.Redis("v").Conn(ctx)
		if err != nil {
			return err
		}
		defer conn.Close(ctx)
		_, err = conn.Do(ctx, "publish", "v:func", funcstring)
		return err
	} else {
		return RunFunc(ctx, funcstring)
	}
}

// ListenFunc 监听函数：订阅 v:func 频道执行集群分发的函数调用。
// 订阅态不跨连接保留，连接失效后必须重连并重新订阅。
func ListenFunc(ctx g.Ctx) {
	if !IsRedisMode {
		panic(gerror.New("集群模式下, 请使用Redis作为缓存"))
	}
	for {
		if err := listenFuncOnce(ctx); err != nil {
			g.Log().Error(ctx, "Redis 订阅连接中断，10秒后重连:", err)
			time.Sleep(10 * time.Second)
		}
	}
}

// listenFuncOnce 建立订阅连接并持续消费，返回错误即代表连接失效
func listenFuncOnce(ctx g.Ctx) error {
	conn, err := g.Redis("v").Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	if _, err = conn.Do(ctx, "subscribe", "v:func"); err != nil {
		return err
	}
	for {
		data, err := conn.Receive(ctx)
		if err != nil {
			return err
		}
		if data != nil {
			dataMap := data.MapStrStr()
			if dataMap["Kind"] == "subscribe" {
				continue
			}
			if dataMap["Channel"] == "v:func" {
				g.Log().Debug(ctx, "执行函数", dataMap["Payload"])
				if e := RunFunc(ctx, dataMap["Payload"]); e != nil {
					g.Log().Error(ctx, "执行函数失败", e)
				}
			}
		}
	}
}
