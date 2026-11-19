package v

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/vera-byte/vgo/v/vconfig"
)

type IController interface {
	Add(ctx context.Context, req *AddReq) (res *BaseRes, err error)
	Delete(ctx context.Context, req *DeleteReq) (res *BaseRes, err error)
	Update(ctx context.Context, req *UpdateReq) (res *BaseRes, err error)
	Info(ctx context.Context, req *InfoReq) (res *BaseRes, err error)
	List(ctx context.Context, req *ListReq) (res *BaseRes, err error)
	Page(ctx context.Context, req *PageReq) (res *BaseRes, err error)
}
type Controller struct {
	Prefix  string     `json:"prefix"`
	Api     g.ArrayStr `json:"api"`
	Service IService   `json:"service"`
}

type AddReq struct {
	g.Meta `path:"/add" method:"POST" summary:"新增" tags:"通用CRUD"`
	Items  []g.Map `json:"items" body:"*"` // 请求体为单条对象或对象数组（对象自动包装为数组）
}

type DeleteReq struct {
	g.Meta `path:"/delete" method:"POST" summary:"删除" tags:"通用CRUD"`
	Ids    []int `json:"ids" v:"required#请选择要删除的数据"`
}

type UpdateReq struct {
	g.Meta `path:"/update" method:"POST" summary:"更新" tags:"通用CRUD"`
	Items  []g.Map `json:"items" body:"*"` // 请求体为单条对象或对象数组（对象自动包装为数组）
}

type InfoReq struct {
	g.Meta `path:"/info" method:"GET" summary:"详情" tags:"通用CRUD"`
	Id     int `json:"id" v:"integer|required#请选择要查询的数据"`
}

// type InfoRes struct {
// 	*BaseRes
// 	Data interface{} `json:"data"`
// }

type ListReq struct {
	g.Meta `path:"/list" method:"POST" summary:"列表" tags:"通用CRUD"`
	Order  string `json:"order"`
	Sort   string `json:"sort"`
}

type PageReq struct {
	g.Meta         `path:"/page" method:"POST" summary:"分页" tags:"通用CRUD"`
	Page           int    `d:"1" json:"page"`     // 页码
	Size           int    `d:"15" json:"size"`    //每页条数
	Order          string `json:"order"`          // 排序字段
	Sort           string `json:"sort"`           // 排序方式 asc desc
	IsExport       bool   `json:"isExport"`       // 是否导出
	MaxExportLimit int    `json:"maxExportLimit"` // 最大导出条数,不传或者小于等于0则不限制
}

func (c *Controller) Add(ctx context.Context, req *AddReq) (res *BaseRes, err error) {
	if !garray.NewStrArrayFrom(c.Api).Contains("Add") {
		g.RequestFromCtx(ctx).Response.Status = 404
		return nil, nil
	}
	if len(req.Items) == 0 {
		return nil, gerror.New("数据不能为空")
	}
	ids := make([]any, 0, len(req.Items))
	for _, item := range req.Items {
		// 与 TS add() 一致：新增剥掉 id，不允许外部指定
		delete(item, "id")
		param := g.MapStrAny(item)
		if _, err = c.Service.ServiceAddOrUpdate(ctx, "Add", param); err != nil {
			return Fail(err.Error()), err
		}
		// 与 TS add() 一致：id 回写进 param 后收集
		ids = append(ids, param["id"])
	}
	// 单条返回 {id}，批量返回 {idList}（idList 批量回写）
	if len(ids) == 1 {
		return Ok(g.Map{"id": ids[0]}), nil
	}
	return Ok(g.Map{"idList": ids}), nil
}
func (c *Controller) Delete(ctx context.Context, req *DeleteReq) (res *BaseRes, err error) {
	if garray.NewStrArrayFrom(c.Api).Contains("Delete") {
		err := c.Service.ModifyBefore(ctx, "Delete", g.RequestFromCtx(ctx).GetMap())
		if err != nil {
			return nil, err
		}

		data, err := c.Service.ServiceDelete(ctx, req)
		if err != nil {
			return Fail(err.Error()), err
		}
		if err = c.Service.ModifyAfter(ctx, "Delete", g.RequestFromCtx(ctx).GetMap()); err != nil {
			return nil, err
		}
		return Ok(data), err
	}
	g.RequestFromCtx(ctx).Response.Status = 404
	return nil, nil
}
func (c *Controller) Update(ctx context.Context, req *UpdateReq) (res *BaseRes, err error) {
	if !garray.NewStrArrayFrom(c.Api).Contains("Update") {
		g.RequestFromCtx(ctx).Response.Status = 404
		return nil, nil
	}
	if len(req.Items) == 0 {
		return nil, gerror.New("数据不能为空")
	}
	ids := make([]any, 0, len(req.Items))
	for _, item := range req.Items {
		param := g.MapStrAny(item)
		if _, err = c.Service.ServiceAddOrUpdate(ctx, "Update", param); err != nil {
			return Fail(err.Error()), err
		}
		ids = append(ids, param["id"])
	}
	// 单条返回 {id}，批量返回 {idList}（idList 批量回写）
	if len(ids) == 1 {
		return Ok(g.Map{"id": ids[0]}), nil
	}
	return Ok(g.Map{"idList": ids}), nil
}
func (c *Controller) Info(ctx context.Context, req *InfoReq) (res *BaseRes, err error) {
	if garray.NewStrArrayFrom(c.Api).Contains("Info") {
		data, err := c.Service.ServiceInfo(ctx, req)
		return Ok(data), err
	}
	g.RequestFromCtx(ctx).Response.Status = 404
	return nil, nil
}
func (c *Controller) List(ctx context.Context, req *ListReq) (res *BaseRes, err error) {
	if garray.NewStrArrayFrom(c.Api).Contains("List") {
		data, err := c.Service.ServiceList(ctx, req)
		return Ok(data), err
	}
	g.RequestFromCtx(ctx).Response.Status = 404
	return nil, nil
}
func (c *Controller) Page(ctx context.Context, req *PageReq) (res *BaseRes, err error) {
	if garray.NewStrArrayFrom(c.Api).Contains("Page") {
		data, err := c.Service.ServicePage(ctx, req)
		return Ok(data), err
	}
	g.RequestFromCtx(ctx).Response.Status = 404
	return nil, nil
}

// 注册控制器到路由
func RegisterController(c IController) {
	var sController = &Controller{}
	gconv.Struct(c, &sController)
	// EPS 字段探测延后到 LoadModelInfo（server 启动钩子）执行，
	// 避免任意子命令在 init 阶段触库
	if sController.Service != nil {
		epsModels[sController.Prefix] = sController.Service.GetModel()
		// 回填所属服务实例，保证编排方法（ServiceAddOrUpdate）链上的子类覆写虚分派生效
		if o, ok := sController.Service.(interface{ SetOwner(IService) }); ok {
			o.SetOwner(sController.Service)
		}
	}
	g.Server().Group(
		sController.Prefix, func(group *ghttp.RouterGroup) {
			group.Middleware(MiddlewareHandlerResponse)
			group.Bind(
				c,
			)
		})
	// body:"*" 字段只接受 JSON 数组，单对象 body 在此包装为单元素数组（兼容 any | any[]）
	g.Server().BindMiddleware(sController.Prefix+"/add", normalizeObjectBodyToArray)
	g.Server().BindMiddleware(sController.Prefix+"/update", normalizeObjectBodyToArray)
}

// normalizeObjectBodyToArray 将单对象 JSON body 包装为数组，使 add/update 同时兼容
// 单条对象与对象数组；数组 body 原样放行。
func normalizeObjectBodyToArray(r *ghttp.Request) {
	body := bytes.TrimSpace(r.GetBody())
	if len(body) > 0 && body[0] == '{' {
		wrapped := make([]byte, 0, len(body)+2)
		wrapped = append(wrapped, '[')
		wrapped = append(wrapped, body...)
		wrapped = append(wrapped, ']')
		r.Body = io.NopCloser(bytes.NewReader(wrapped))
		r.ReloadParam()
	}
	r.Middleware.Next()
}

// ColumnInfo 表字段信息
type ColumnInfo struct {
	Comment      string `json:"comment"`
	Length       string `json:"length"`
	Nullable     bool   `json:"nullable"`
	PropertyName string `json:"propertyName"`
	Type         string `json:"type"`
}

// ModelInfo 路由prefix	对应的model信息
var ModelInfo = make(map[string][]*ColumnInfo)

// epsModels 已注册路由 -> 模型（EPS 探测延后执行的登记表）
var epsModels = make(map[string]IModel)

func init() {
	RegisterStartupHook("eps-model-info", LoadModelInfo)
}

// LoadModelInfo 探测并填充 ModelInfo（幂等；server 启动时经启动钩子调用）
func LoadModelInfo(ctx context.Context) {
	if !vconfig.Config.Eps {
		return
	}
	for prefix, model := range epsModels {
		if model == nil {
			continue
		}
		ModelInfo[prefix] = getModelInfo(ctx, prefix, model)
	}
}

// getModelInfo 获取模型信息
func getModelInfo(ctx g.Ctx, prefix string, model IModel) (columns []*ColumnInfo) {
	fields, err := g.DB(model.GroupName()).TableFields(ctx, model.TableName())
	if err != nil {
		// 表不存在或数据库暂不可用时（如全新环境先执行迁移）跳过EPS字段信息，不阻断启动
		g.Log().Warningf(ctx, "获取表 %s 字段信息失败: %v", model.TableName(), err)
		return nil
	}
	// g.Log().Info(ctx, "fields", fields)
	sortedFields := garray.NewArraySize(len(fields), len(fields))
	for _, field := range fields {
		sortedFields.Set(field.Index, field)
	}
	for _, field := range sortedFields.Slice() {
		if field.(*gdb.TableField).Name == "deleted_at" {
			continue
		}
		var comment string
		if field.(*gdb.TableField).Comment != "" {
			comment = field.(*gdb.TableField).Comment
		} else {
			comment = field.(*gdb.TableField).Name
		}
		// 去除 type中的长度
		var length string
		if strings.Contains(field.(*gdb.TableField).Type, "(") {
			length = field.(*gdb.TableField).Type[strings.Index(field.(*gdb.TableField).Type, "(")+1 : strings.Index(field.(*gdb.TableField).Type, ")")]
		}
		columnType := gstr.Replace(field.(*gdb.TableField).Type, "("+length+")", "")
		column := &ColumnInfo{
			Comment:      comment,
			Length:       "",
			Nullable:     field.(*gdb.TableField).Null,
			PropertyName: field.(*gdb.TableField).Name,
			Type:         columnType,
		}
		columns = append(columns, column)
	}
	return
}
