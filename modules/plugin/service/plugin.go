package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/vera-byte/vgo/modules/plugin/model"
	"github.com/vera-byte/vgo/v"
)

// PluginService 插件信息服务（参照 cool-admin-midway PluginService；
// 插件本体为可执行程序，经 stdio JSON-RPC 与宿主通信）
type PluginService struct {
	*v.Service
}

// NewPluginService 创建插件服务实例
func NewPluginService() *PluginService {
	return &PluginService{
		Service: &v.Service{
			Model: model.NewPluginInfo(),
			PageQueryOp: &v.QueryOp{
				Select:     `plugin_info."id",plugin_info."name",plugin_info."keyName",plugin_info."hook",plugin_info."version",plugin_info."status",plugin_info."readme",plugin_info."author",plugin_info."logo",plugin_info."description",plugin_info."pluginJson",plugin_info."config",plugin_info."createTime",plugin_info."updateTime"`,
				AddOrderby: g.MapStrStr{"id": "DESC"},
			},
		},
	}
}

// pluginZipData ZIP 插件包解析结果
type pluginZipData struct {
	PluginJSON map[string]interface{} `json:"pluginJson"`
	Readme     string                 `json:"readme"`
	Logo       string                 `json:"logo"`
	Bin        []byte                 `json:"-"` // 可执行文件内容
	BinPath    string                 `json:"binPath"`
}

// binEntry 按平台优先级寻找可执行文件（含 Windows .exe）
func binEntry(files map[string]*zip.File) *zip.File {
	plat := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	for _, name := range []string{
		fmt.Sprintf("bin/%s/plugin.exe", plat),
		fmt.Sprintf("bin/%s/plugin", plat),
		"bin/plugin.exe",
		"bin/plugin",
		"plugin.exe",
		"plugin",
	} {
		if f, ok := files[name]; ok {
			return f
		}
	}
	return nil
}

// binName 平台可执行文件名（Windows 需 .exe 后缀）
func binName() string {
	if runtime.GOOS == "windows" {
		return "plugin.exe"
	}
	return "plugin"
}

// ParseZip 解析插件 ZIP 包：plugin.json / readme / logo / 可执行文件
func ParseZip(raw []byte) (*pluginZipData, error) {
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return nil, gerror.Wrap(err, "插件包不是有效的 ZIP")
	}
	files := map[string]*zip.File{}
	for _, f := range zr.File {
		files[strings.TrimPrefix(f.Name, "./")] = f
	}
	read := func(name string) (string, error) {
		f, ok := files[name]
		if !ok {
			return "", fmt.Errorf("插件包缺少文件: %s", name)
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		b, err := io.ReadAll(rc)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	readB64 := func(name string) (string, error) {
		f, ok := files[name]
		if !ok {
			return "", fmt.Errorf("插件包缺少文件: %s", name)
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		b, err := io.ReadAll(rc)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(b), nil
	}

	pluginJSONRaw, err := read("plugin.json")
	if err != nil {
		return nil, gerror.Wrap(err, "插件信息不完整")
	}
	var pj map[string]interface{}
	if err = json.Unmarshal([]byte(pluginJSONRaw), &pj); err != nil {
		return nil, gerror.Wrap(err, "plugin.json 解析失败")
	}
	out := &pluginZipData{PluginJSON: pj}

	if name := gconv.String(pj["readme"]); name != "" {
		out.Readme, _ = read(name)
	}
	if name := gconv.String(pj["logo"]); name != "" {
		out.Logo, _ = readB64(name)
	}
	bin := binEntry(files)
	if bin == nil {
		return nil, gerror.New("插件包缺少可执行文件（bin/<goos>-<goarch>/plugin 或 bin/plugin）")
	}
	bf, err := bin.Open()
	if err != nil {
		return nil, err
	}
	defer bf.Close()
	if out.Bin, err = io.ReadAll(bf); err != nil {
		return nil, err
	}
	out.BinPath = bin.Name
	return out, nil
}

// CheckResult 安装前检查结果
type CheckResult struct {
	Type    int    `json:"type"`    // 0-信息不完整 1-已存在(覆盖) 2-同名Hook插件已启用 3-检查通过
	Message string `json:"message"`
}

// Check 安装前检查
func (s *PluginService) Check(raw []byte) (*CheckResult, error) {
	data, err := ParseZip(raw)
	if err != nil {
		return &CheckResult{Type: 0, Message: "插件信息不完整：" + err.Error()}, nil
	}
	key := gconv.String(data.PluginJSON["key"])
	record, _ := v.DBM(model.NewPluginInfo()).Where("keyName", key).One()
	if !record.IsEmpty() && record["hook"].String() == "" {
		return &CheckResult{Type: 1, Message: "插件已存在，继续安装将覆盖"}, nil
	}
	if !record.IsEmpty() && record["hook"].String() != "" && record["status"].Int() == 1 {
		return &CheckResult{Type: 2, Message: "已存在同名Hook插件，你可以继续安装，但是多个相同的Hook插件只能同时开启一个"}, nil
	}
	return &CheckResult{Type: 3, Message: "检查通过"}, nil
}

// pluginDir 插件安装目录
func pluginDir(key string) string {
	return filepath.Join("data", "plugins", key)
}

// Install 安装插件（ZIP 包；force 跳过检查）
func (s *PluginService) Install(ctx context.Context, raw []byte, force bool) (*CheckResult, error) {
	check, err := s.Check(raw)
	if err != nil {
		return nil, err
	}
	if check.Type != 3 && !force {
		return check, nil
	}
	data, err := ParseZip(raw)
	if err != nil {
		return nil, err
	}
	key := gconv.String(data.PluginJSON["key"])
	if key == "" || key == "plugin" {
		return nil, gerror.New("插件key不能为空且不能为plugin，请更换其他key")
	}

	// 落盘可执行文件
	dir := pluginDir(key)
	if err = gfile.Mkdir(dir); err != nil {
		return nil, err
	}
	binPath := filepath.Join(dir, binName())
	if err = os.WriteFile(binPath, data.Bin, 0o755); err != nil {
		return nil, gerror.Wrap(err, "插件可执行文件写入失败")
	}

	pluginJSONRaw, _ := json.Marshal(data.PluginJSON)
	configJSON, _ := json.Marshal(data.PluginJSON["config"])
	fields := g.Map{
		"name":        gconv.String(data.PluginJSON["name"]),
		"keyName":     key,
		"version":     gconv.String(data.PluginJSON["version"]),
		"author":      gconv.String(data.PluginJSON["author"]),
		"hook":        gconv.String(data.PluginJSON["hook"]),
		"readme":      data.Readme,
		"logo":        data.Logo,
		"description": gconv.String(data.PluginJSON["description"]),
		"pluginJson":  string(pluginJSONRaw),
		"config":      string(configJSON),
		"status":      1,
	}

	existing, _ := v.DBM(model.NewPluginInfo()).Where("keyName", key).One()
	if !existing.IsEmpty() {
		fields["status"] = existing["status"].Int()
		merged := mergeConfig(gconv.Map(existing["config"]), gconv.Map(data.PluginJSON["config"]))
		cb, _ := json.Marshal(merged)
		fields["config"] = string(cb)
		_, err = v.DBM(model.NewPluginInfo()).Where("id", existing["id"]).Data(fields).Update()
	} else {
		_, err = v.DBM(model.NewPluginInfo()).Data(fields).Insert()
	}
	if err != nil {
		return nil, err
	}
	s.ReloadKey(key)
	return &CheckResult{Type: 3, Message: "安装成功"}, nil
}

// mergeConfig 合并配置：插件包默认值在下，既有用户配置覆盖
func mergeConfig(old, fresh map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range fresh {
		out[k] = v
	}
	for k, v := range old {
		out[k] = v
	}
	return out
}

// ServiceUpdate 更新后联动热加载/移除（同 Hook 互斥）
func (s *PluginService) ServiceUpdate(ctx context.Context, req *v.UpdateReq) (data interface{}, err error) {
	data, err = s.Service.ServiceUpdate(ctx, req)
	if err != nil {
		return
	}
	id := g.RequestFromCtx(ctx).Get("id").Uint()
	if id == 0 {
		record, _ := v.DBM(model.NewPluginInfo()).Where("keyName", g.RequestFromCtx(ctx).Get("keyName")).One()
		id = record["id"].Uint()
	}
	if id > 0 {
		err = s.AddOrUpdate(ctx, id)
	}
	return
}

// ServiceDelete 删除插件（不走回收站，连同安装目录清理）
func (s *PluginService) ServiceDelete(ctx context.Context, req *v.DeleteReq) (data interface{}, err error) {
	ids := gconv.Ints(g.RequestFromCtx(ctx).Get("ids").Slice())
	err = s.Delete(ctx, ids)
	return
}

// AddOrUpdate 新增或更新后联动热加载/移除
func (s *PluginService) AddOrUpdate(ctx context.Context, id uint) error {
	record, err := v.DBM(model.NewPluginInfo()).Where("id", id).One()
	if err != nil {
		return err
	}
	if record.IsEmpty() {
		return gerror.New("插件不存在")
	}
	key := record["keyName"].String()
	hook := record["hook"].String()
	if record["status"].Int() == 1 {
		if hook != "" {
			// 启用 Hook 插件时禁用同 Hook 的其他插件
			_, _ = v.DBM(model.NewPluginInfo()).
				Where("hook", hook).Where("status", 1).
				Where("id <> ?", id).
				Data(g.Map{"status": 0}).Update()
		}
		s.ReloadKey(key)
	} else {
		GetPluginCenter().Remove(key, false)
		if hook != "" {
			GetPluginCenter().Remove(hook, false)
		}
	}
	return nil
}

// Delete 删除插件（不走回收站）
func (s *PluginService) Delete(ctx context.Context, ids []int) error {
	for _, id := range ids {
		record, _ := v.DBM(model.NewPluginInfo()).Where("id", id).One()
		if record.IsEmpty() {
			continue
		}
		key := record["keyName"].String()
		hook := record["hook"].String()
		GetPluginCenter().Remove(key, false)
		if hook != "" {
			GetPluginCenter().Remove(hook, false)
		}
		_ = os.RemoveAll(pluginDir(key))
		if _, err := v.DBM(model.NewPluginInfo()).Unscoped().Where("id", id).Delete(); err != nil {
			return err
		}
	}
	return nil
}

// ReloadKey 重新加载单个插件（启用/安装后热更新：起进程并 ready 握手）
func (s *PluginService) ReloadKey(key string) {
	record, _ := v.DBM(model.NewPluginInfo()).Where("keyName", key).One()
	if record.IsEmpty() || record["status"].Int() != 1 {
		GetPluginCenter().Remove(key, false)
		return
	}
	binPath := filepath.Join(pluginDir(key), binName())
	if !gfile.Exists(binPath) {
		g.Log().Warningf(context.Background(), "插件[%s]可执行文件不存在: %s", key, binPath)
		return
	}
	info := BuildPluginInfoData(record.Map(), pluginDir(key))
	// 停旧起新（同 Hook 键下重复加载时覆盖）
	GetPluginCenter().Remove(info.Hook, false)
	GetPluginCenter().Remove(key, false)
	if err := LoadProcessPlugin(key, info.Hook, binPath, info); err != nil {
		g.Log().Warningf(context.Background(), "插件[%s]加载失败: %v", key, err)
	}
}

// InitAll 启动初始化：加载全部启用插件
func (s *PluginService) InitAll() {
	records, err := v.DBM(model.NewPluginInfo()).Where("status", 1).All()
	if err != nil {
		return
	}
	for _, record := range records {
		s.ReloadKey(record["keyName"].String())
	}
}

// BuildPluginInfoData 从数据库行构造插件信息（含 @env 多环境配置解析）
func BuildPluginInfoData(row map[string]interface{}, baseDir string) *PluginInfoData {
	var cfg map[string]interface{}
	if err := gconv.Struct(row["config"], &cfg); err != nil || cfg == nil {
		cfg = map[string]interface{}{}
	}
	isMulti := false
	for k := range cfg {
		if strings.HasPrefix(k, "@") {
			isMulti = true
			break
		}
	}
	if isMulti {
		sub := gconv.Map(cfg["@local"])
		if sub != nil {
			cfg = sub
		}
	}
	var pj map[string]interface{}
	_ = gconv.Struct(row["pluginJson"], &pj)
	return &PluginInfoData{
		Name:        gconv.String(row["name"]),
		Key:         gconv.String(row["keyName"]),
		Hook:        gconv.String(row["hook"]),
		Singleton:   gconv.Bool(pj["singleton"]),
		Version:     gconv.String(row["version"]),
		Description: gconv.String(row["description"]),
		Author:      gconv.String(row["author"]),
		BaseDir:     baseDir,
		Config:      cfg,
	}
}
