package v

import "github.com/gogf/gf/v2/frame/g"

// QuoteIdent 按数据库方言引用标识符（支持 table.column 形式）。
// PostgreSQL 会将未加引号的标识符折叠为小写，手工拼接的 JOIN ON 条件、
// 带表名限定的查询字段必须经此引用，避免 camelCase 列名失效。
func QuoteIdent(group, ident string) string {
	return g.DB(group).GetCore().QuoteString(ident)
}
