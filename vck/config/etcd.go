package vck_config

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	etcd3 "go.etcd.io/etcd/client/v3"
)

// ChainableEtcdClient 用于链式调用的结构体
type ChainableEtcdClient struct {
	Client *etcd3.Client
}

var (
	EtcdManager *ChainableEtcdClient
)

// 初始化 Etcd 客户端并返回 ChainableEtcdClient
func NewChainableEtcdClient() *ChainableEtcdClient {
	if EtcdManager != nil {
		return EtcdManager
	}
	var ctx = gctx.GetInitCtx()
	conf, err := g.Cfg("etcd").Get(ctx, "etcd.address")
	if err == nil {
		var address = conf.Strings()
		g.Log().Debug(ctx, "Etcd配置中心开启")
		client, err := etcd3.New(etcd3.Config{
			Endpoints:   address,         // Etcd 服务地址
			DialTimeout: 5 * time.Second, // 连接超时
		})
		if err != nil {
			return nil
		}
		EtcdManager = &ChainableEtcdClient{Client: client}
		return EtcdManager

	}

	return nil

}

// 链式调用：从 Etcd 获取配置信息
func (c *ChainableEtcdClient) GetConfig(key string) (*g.Var, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 从 Etcd 获取配置值
	resp, err := c.Client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		err = gerror.Newf("etcd  %s not found", key)
		g.Log().Warning(ctx, err)
		return nil, err
	}

	return g.NewVar(resp.Kvs[0].Value), nil

}

// 链式调用：从 Etcd 写入信息
func (c *ChainableEtcdClient) PutConfig(key string, val interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 从 Etcd 获取配置值
	_, err := c.Client.Put(ctx, key, gconv.String(val))
	return err

}
