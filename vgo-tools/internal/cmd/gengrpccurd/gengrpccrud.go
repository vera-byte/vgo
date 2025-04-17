package gengrpccurd

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/olekukonko/tablewriter"
	"github.com/vera-byte/vgo/vgo-tools/internal/consts"
	"github.com/vera-byte/vgo/vgo-tools/utility/mlog"
	"github.com/vera-byte/vgo/vgo-tools/utility/utils"
)

const (
	CGenGrpcCrudConfig   = `vgocli.gen.grpc.crud` // 配置项 Key
	defaultPackageSuffix = `api`                  // 默认生成文件路径后缀
)

func init() {
	// 设置标签，供 gtag 使用
	gtag.Sets(g.MapStrStr{
		`CGenGrpcCrudConfig`: CGenGrpcCrudConfig,
	})
}

// 默认数据库字段类型到 proto 类型的映射关系
var defaultTypeMapping = map[DBFieldTypeName]CustomAttributeType{
	"string":       {Type: "string"},
	"date":         {Type: "google.protobuf.Timestamp", Import: "google/protobuf/timestamp.proto"},
	"datetime":     {Type: "google.protobuf.Timestamp", Import: "google/protobuf/timestamp.proto"},
	"int":          {Type: "int32"},
	"uint":         {Type: "uint32"},
	"int64":        {Type: "int64"},
	"uint64":       {Type: "uint64"},
	"[]int":        {Type: "repeated int32"},
	"[]int64":      {Type: "repeated int64"},
	"[]uint64":     {Type: "repeated uint64"},
	"int64-bytes":  {Type: "repeated int64"},
	"uint64-bytes": {Type: "repeated uint64"},
	"float32":      {Type: "float"},
	"float64":      {Type: "double"},
	"[]byte":       {Type: "bytes"},
	"bool":         {Type: "bool"},
}

type (
	// 主结构体
	CGenGrpcCrud struct{}

	// CLI 输入结构体
	CGenGrpcCrudInput struct {
		g.Meta            `name:"grpc-crud" config:"{CGenPbEntityConfig}" brief:"{CGenPbEntityBrief}" eg:"{CGenPbEntityEg}" ad:"{CGenPbEntityAd}"`
		SapceName         string                                   `name:"sapceName"              short:"sn"  d:"app"`
		Version           string                                   `name:"version"              short:"v"  d:"v1"`
		Path              string                                   `name:"path"              short:"p"  brief:"{CGenPbEntityBriefPath}" d:"manifest/protobuf"`
		Package           string                                   `name:"package"           short:"k"  brief:"{CGenPbEntityBriefPackage}"`
		GoPackage         string                                   `name:"goPackage"           short:"g"  brief:"{CGenPbEntityBriefGoPackage}"`
		Link              string                                   `name:"link"              short:"l"  brief:"{CGenPbEntityBriefLink}"`
		Tables            string                                   `name:"tables"            short:"t"  brief:"{CGenPbEntityBriefTables}"`
		Prefix            string                                   `name:"prefix"            short:"f"  brief:"{CGenPbEntityBriefPrefix}"`
		RemovePrefix      string                                   `name:"removePrefix"      short:"r"  brief:"{CGenPbEntityBriefRemovePrefix}"`
		RemoveFieldPrefix string                                   `name:"removeFieldPrefix" short:"rf" brief:"{CGenPbEntityBriefRemoveFieldPrefix}"`
		TablesEx          string                                   `name:"tablesEx"          short:"x"  brief:"{CGenDaoBriefTablesEx}"`
		NameCase          string                                   `name:"nameCase"          short:"n"  brief:"{CGenPbEntityBriefNameCase}" d:"Camel"`
		JsonCase          string                                   `name:"jsonCase"          short:"j"  brief:"{CGenPbEntityBriefJsonCase}" d:"none"`
		Option            string                                   `name:"option"            short:"o"  brief:"{CGenPbEntityBriefOption}"`
		Crud              string                                   `name:"crud" short:"c"  d:"Add,Update,Info,Delete,Page,List"`
		TypeMapping       map[DBFieldTypeName]CustomAttributeType  `name:"typeMapping"  short:"y"  brief:"{CGenPbEntityBriefTypeMapping}"  orphan:"true"`
		FieldMapping      map[DBTableFieldName]CustomAttributeType `name:"fieldMapping" short:"fm" brief:"{CGenPbEntityBriefFieldMapping}" orphan:"true"`
	}
	// 输出结构体
	CGenGrpcCrudOutput struct{}

	// 内部调用用结构体（带表名、数据库等）
	CGenPbEntityInternalInput struct {
		CGenGrpcCrudInput
		DB           gdb.DB
		TableName    string
		NewTableName string
	}

	DBTableFieldName = string
	DBFieldTypeName  = string

	// 字段类型及其对应 proto 类型和 import
	CustomAttributeType struct {
		Type   string
		Import string
	}
)

// 主入口：执行 CRUD 代码生成
func (c CGenGrpcCrud) GenGrpcCrud(ctx context.Context, in CGenGrpcCrudInput) (out *CGenGrpcCrudOutput, err error) {
	var config = g.Cfg()

	if config.Available(ctx) {
		v := config.MustGet(ctx, CGenGrpcCrudConfig)
		if v.IsSlice() {
			for i := 0; i < len(v.Interfaces()); i++ {
				genGrpcCrud(ctx, i, in)
			}
		} else {
			genGrpcCrud(ctx, -1, in)
		}
	} else {
		genGrpcCrud(ctx, -1, in)
	}

	mlog.Print("✅完成!")
	return
}

// 单个配置生成处理逻辑
func genGrpcCrud(ctx context.Context, index int, in CGenGrpcCrudInput) {
	var db gdb.DB
	var err error
	in.Path = fmt.Sprintf("%s/%s/%s", in.Path, in.SapceName, in.Version)
	// 多个配置时读取对应项
	if index >= 0 {
		err = g.Cfg().MustGet(ctx, fmt.Sprintf(`%s.%d`, CGenGrpcCrudConfig, index)).Scan(&in)
		if err != nil {
			mlog.Fatalf(`无效配置 "%s": %+v`, CGenGrpcCrudConfig, err)
		}
	}

	// 如果没设置 package，自动从 go.mod 计算
	if in.Package == "" {
		modName := utils.GetImportPath(gfile.Pwd())
		in.Package = modName + "/" + defaultPackageSuffix
	}

	// 要移除的表前缀数组
	removePrefixArray := gstr.SplitAndTrim(in.RemovePrefix, ",")
	excludeTables := gset.NewStrSetFrom(gstr.SplitAndTrim(in.TablesEx, ","))

	// 初始化数据库连接
	if in.Link != "" {
		match, _ := gregex.MatchString(`([a-z]+):(.+)`, in.Link)
		if len(match) == 3 {
			tempGroup := gtime.TimestampNanoStr()
			gdb.AddConfigNode(tempGroup, gdb.ConfigNode{
				Type: gstr.Trim(match[1]),
				Link: in.Link,
			})
			db, _ = gdb.Instance(tempGroup)
		}
	} else {
		db = g.DB()
	}
	if db == nil {
		mlog.Fatal("❌数据库初始化失败")
	}

	// 获取表名
	var tableNames []string
	if in.Tables != "" {
		tableNames = gstr.SplitAndTrim(in.Tables, ",")
	} else {
		tableNames, err = db.Tables(context.TODO())
		if err != nil {
			mlog.Fatalf("❌获取表失败: \n %v", err)
		}
	}

	// 合并默认类型映射
	if in.TypeMapping == nil {
		in.TypeMapping = defaultTypeMapping
	} else {
		for k, v := range defaultTypeMapping {
			if _, ok := in.TypeMapping[k]; !ok {
				in.TypeMapping[k] = v
			}
		}
	}

	// 遍历每张表生成 .proto 文件
	for _, tableName := range tableNames {
		if excludeTables.Contains(tableName) {
			continue
		}
		newTableName := tableName
		for _, v := range removePrefixArray {
			newTableName = gstr.TrimLeftStr(newTableName, v, 1)
		}
		generatePbEntityContentFile(ctx, CGenPbEntityInternalInput{
			CGenGrpcCrudInput: in,
			DB:                db,
			TableName:         tableName,
			NewTableName:      newTableName,
		})
	}
}

// 生成 .proto 文件内容并保存
func generatePbEntityContentFile(ctx context.Context, in CGenPbEntityInternalInput) {
	fieldMap, err := in.DB.TableFields(ctx, in.TableName)
	if err != nil {
		mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", in.TableName, err)
	}

	newTableName := in.Prefix + in.NewTableName
	tableNameCamelCase := gstr.CaseCamel(newTableName)
	tableNameSnakeCase := gstr.CaseSnake(newTableName)
	// 生成消息体
	entityMessageDefine, appendImports := generateEntityMessageDefinition(tableNameCamelCase, fieldMap, in)
	// 生成crud调用方法
	entityRpcInvoke, entityFieldPriKey := generateEntityRpcInvoke(tableNameCamelCase, fieldMap, in)
	// 生成RpcService实现
	g.Dump(entityRpcInvoke, entityFieldPriKey)

	fileName := gstr.Trim(tableNameSnakeCase, "-_.")
	path := filepath.FromSlash(gfile.Join(in.Path, fmt.Sprintf("%s.crud", fileName)+".proto"))

	// 处理导入
	packageImportsArray := garray.NewStrArray()
	for _, appendImport := range appendImports {
		importStr := fmt.Sprintf(`import "%s";`, appendImport)
		if packageImportsArray.Search(importStr) == -1 {
			packageImportsArray.Append(importStr)
		}
	}

	if in.GoPackage == "" {
		in.GoPackage = fmt.Sprintf("%s/%s/%s/protobuf", in.Package, in.Version, in.SapceName)
	}

	// 渲染模板
	entityContent := gstr.ReplaceByMap(getTplPbEntityContent(""), g.MapStrStr{
		"{Imports}":             packageImportsArray.Join("\n"),
		"{PackageName}":         gfile.Basename(in.Package),
		"{GoPackage}":           in.GoPackage,
		"{OptionContent}":       in.Option,
		"{EntityMessage}":       entityMessageDefine,
		"{RpcService}":          generatedServiceContent(tableNameCamelCase, entityRpcInvoke),
		"{RpcServiceImplement}": generatedRpcServiceImplementContent(tableNameCamelCase, entityRpcInvoke, entityFieldPriKey, in),
	})

	// 写入 proto 文件
	if err := gfile.PutContents(path, strings.TrimSpace(entityContent)); err != nil {
		mlog.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		mlog.Print("generated:", gfile.RealPath(path))
	}
}

// ParseRpcLine 从 rpc 定义中解析出方法名、请求、响应结构体
func ParseRpcLine(rpcLine string) (name, req, res string) {
	// 示例匹配：rpc Delete(xxxxxxDeleteRpcInvoke) returns (xxxxxxDeleteRpcRes)
	re := regexp.MustCompile(`rpc\s+(\w+)\s*\(\s*(\w+)\s*\)\s+returns\s+\(\s*(\w+)\s*\)`)
	matches := re.FindStringSubmatch(rpcLine)
	if len(matches) == 4 {
		return matches[1], matches[2], matches[3]
	}
	return "", "", ""
}

// 生成Service内容
func generatedServiceContent(sn string, s map[string]string) string {

	sContent := ""
	for name, rpc := range s {
		sContent += fmt.Sprintf("    // %s\n    %s;\n", name, rpc)
	}
	return gstr.ReplaceByMap(consts.TemplateServiceContent, g.MapStrStr{
		"{ServiceName}": sn,
		"{Rpc}":         sContent,
	})
}

// 生成RpcService实现
func generatedRpcServiceImplementContent(table string, s map[string]string, priKey *gdb.TableField, in CGenPbEntityInternalInput) string {
	var (
		sContent = ""
	)
	for _, rawRpc := range s {
		name, req, res := ParseRpcLine(rawRpc)
		switch name {
		case "Add":
			sContent += generateEntityCreate(req, res, table)
		case "Info":
			sContent += generateEntityFindOne(req, res, priKey, table, in)
		case "Delete":
			sContent += generateEntityDelets(req, res, priKey, in)
		case "Page":
			sContent += generateEntityPage(req, res, table)
		case "List":
			sContent += generateEntityFindAll(req, res, table)
		case "Update":
			sContent += generateEntityUpdate(req, res, priKey, in)
		}
	}
	return sContent
}

// 注册调用方法
func generateEntityRpcInvoke(entityName string, fieldMap map[string]*gdb.TableField, in CGenPbEntityInternalInput) (map[string]string, *gdb.TableField) {
	var (
		curds             = gstr.Split(in.Crud, ",")
		entityFieldPriKey *gdb.TableField
		rpc               = map[string]string{}
	)
	if len(curds) == 0 {
		mlog.Fatal("crud 为空")
	}
	for _, field := range fieldMap {
		if field.Key == "pri" {
			entityFieldPriKey = field
			break
		}
	}
	if entityFieldPriKey == nil {
		mlog.Fatalf("%s 表中没有主键,请至少设置一个主键", entityName)
	}
	for _, method := range curds {
		buffer := bytes.NewBuffer(nil)
		// 使用 tablewriter 美化格式
		tw := tablewriter.NewWriter(buffer)
		tw.SetBorder(false)
		tw.SetRowLine(false)
		tw.SetAutoWrapText(false)
		tw.SetColumnSeparator("")
		tw.Render()
		stContent := buffer.String()
		buffer.Reset()
		buffer.WriteString(fmt.Sprintf("rpc %s(%s%sRpcInvoke) returns (%s%sRpcRes)", method, entityName, method, entityName, method))
		buffer.WriteString(stContent)
		rpc[method] = buffer.String()
	}
	return rpc, entityFieldPriKey
}

// 构建 rpc查询方法(对应Info接口)
func generateEntityFindOne(req string, res string, priKey *gdb.TableField, table string, in CGenPbEntityInternalInput) string {
	entityMessageDefine, _ := generateEntityMessageDefinition(req, map[string]*gdb.TableField{"id": priKey}, in)
	return gstr.ReplaceByMap(consts.TemplateRpcServiceImplement, g.MapStrStr{
		"{RpcInvoke}": entityMessageDefine,
		"{RpcRes}":    res,
		"{ResField}":  fmt.Sprintf("    %s %s = 1;\n", table, table),
	})
}

// 构建 rpc更新方法(对应Update接口)
func generateEntityUpdate(req string, res string, priKey *gdb.TableField, in CGenPbEntityInternalInput) string {
	temp := priKey
	entityMessageDefine, _ := generateEntityMessageDefinition(req, map[string]*gdb.TableField{"id": temp}, in)
	return gstr.ReplaceByMap(consts.TemplateRpcServiceImplement, g.MapStrStr{
		"{RpcInvoke}": entityMessageDefine,
		"{RpcRes}":    res,
		"{ResField}":  "",
	})
}

// 构建 rpc查询方法(对应List接口)
func generateEntityFindAll(req string, res string, table string) string {
	return gstr.ReplaceByMap(consts.TemplateRpcServiceImplement, g.MapStrStr{
		"{RpcInvoke}": fmt.Sprintf("message %s {}", req),
		"{RpcRes}":    res,
		"{ResField}":  fmt.Sprintf("repeated %s %s = 1;\n", table, table),
	})
}

// 构建 rpc查询方法(对应delete接口)
func generateEntityDelets(req string, res string, priKey *gdb.TableField, in CGenPbEntityInternalInput) string {
	entityMessageDefine, _ := generateEntityMessageDefinition(req, map[string]*gdb.TableField{priKey.Name: priKey}, in)
	return gstr.ReplaceByMap(consts.TemplateRpcServiceImplement, g.MapStrStr{
		"{RpcInvoke}": entityMessageDefine,
		"{RpcRes}":    res,
		"{ResField}":  "",
	})
}

// 构建 rpc查询方法(对应add接口)
func generateEntityCreate(req string, res string, table string) string {
	return gstr.ReplaceByMap(consts.TemplateRpcServiceImplement, g.MapStrStr{
		"{RpcInvoke}": fmt.Sprintf(`message %s {
			 %s %s = 1;
		}`, req, table, table),
		"{RpcRes}":   res,
		"{ResField}": "",
	})
}

// 构建 rpc查询方法(对应page接口)
func generateEntityPage(req string, res string, table string) string {
	return gstr.ReplaceByMap(consts.TemplateRpcServiceImplement, g.MapStrStr{
		"{RpcInvoke}": fmt.Sprintf("message %s {%s %s = 1;}", req, table, table),
		"{RpcRes}":    res,
		"{ResField}":  "",
	})
}

// 构建 message 定义体
func generateEntityMessageDefinition(entityName string, fieldMap map[string]*gdb.TableField, in CGenPbEntityInternalInput) (string, []string) {
	var (
		appendImports []string
		buffer        = bytes.NewBuffer(nil)
		array         = make([][]string, len(fieldMap))
		names         = sortFieldKeyForPbEntity(fieldMap)
	)

	for index, name := range names {
		var imports string
		array[index], imports = generateMessageFieldForPbEntity(index+1, fieldMap[name], in)
		if imports != "" {
			appendImports = append(appendImports, imports)
		}
	}

	// 使用 tablewriter 美化格式
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()

	stContent := buffer.String()
	stContent = regexp.MustCompile(`\s+\n`).ReplaceAllString(gstr.Replace(stContent, "  #", ""), "\n")

	buffer.Reset()
	buffer.WriteString(fmt.Sprintf("message %s {\n", entityName))
	buffer.WriteString(stContent)
	buffer.WriteString("}")
	return buffer.String(), appendImports
}

// 获取 proto 模板内容
func getTplPbEntityContent(tplEntityPath string) string {
	if tplEntityPath != "" {
		return gfile.GetContents(tplEntityPath)
	}
	return consts.TemplatePbEntityMessageContent
}

// 生成单个字段的 proto 字段定义
func generateMessageFieldForPbEntity(index int, field *gdb.TableField, in CGenPbEntityInternalInput) (attrLines []string, appendImport string) {
	var (
		localTypeNameStr string
		localTypeName    gdb.LocalType
		comment          string
		jsonTagStr       string
		err              error
		ctx              = gctx.GetInitCtx()
	)

	// 获取字段类型
	if in.TypeMapping != nil {
		localTypeName, err = in.DB.CheckLocalTypeForField(ctx, field.Type, nil)
		if err != nil {
			panic(err)
		}
		if localTypeName != "" {
			if typeMapping, ok := in.TypeMapping[strings.ToLower(string(localTypeName))]; ok {
				localTypeNameStr = typeMapping.Type
				appendImport = typeMapping.Import
			}
		}
	}
	if localTypeNameStr == "" {
		localTypeNameStr = "string"
	}

	// 处理注释
	comment = gstr.ReplaceByArray(field.Comment, g.SliceStr{"\n", " ", "\r", " "})
	comment = gstr.Trim(comment)
	comment = gstr.Replace(comment, `\n`, " ")
	comment, _ = gregex.ReplaceString(`\s{2,}`, ` `, comment)

	// json tag
	if jsonTagName := formatCase(field.Name, in.JsonCase); jsonTagName != "" {
		jsonTagStr = fmt.Sprintf(`[json_name = "%s"]`, jsonTagName)
		if index < 10 {
			jsonTagStr = "   " + jsonTagStr
		} else if index < 100 {
			jsonTagStr = "  " + jsonTagStr
		} else {
			jsonTagStr = " " + jsonTagStr
		}
	}

	// 移除字段前缀
	removeFieldPrefixArray := gstr.SplitAndTrim(in.RemoveFieldPrefix, ",")
	newFiledName := field.Name
	for _, v := range removeFieldPrefixArray {
		newFiledName = gstr.TrimLeftStr(newFiledName, v, 1)
	}

	// 字段级别映射
	if in.FieldMapping != nil {
		if typeMapping, ok := in.FieldMapping[fmt.Sprintf("%s.%s", in.TableName, newFiledName)]; ok {
			localTypeNameStr = typeMapping.Type
			appendImport = typeMapping.Import
		}
	}

	return []string{
		"    #" + localTypeNameStr,
		" #" + formatCase(newFiledName, in.NameCase),
		" #= " + gconv.String(index) + jsonTagStr + ";",
		" #" + fmt.Sprintf(`// %s`, comment),
	}, appendImport
}

// 格式化字段命名方式
func formatCase(str, caseStr string) string {
	if caseStr == "none" {
		return ""
	}
	return gstr.CaseConvert(str, gstr.CaseTypeMatch(caseStr))
}

// 按字段顺序排序
func sortFieldKeyForPbEntity(fieldMap map[string]*gdb.TableField) []string {
	names := make(map[int]string)
	for _, field := range fieldMap {
		names[field.Index] = field.Name
	}
	result := make([]string, len(names))
	i, j := 0, 0
	for {
		if len(names) == 0 {
			break
		}
		if val, ok := names[i]; ok {
			result[j] = val
			j++
			delete(names, i)
		}
		i++
	}
	return result
}
