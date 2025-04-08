package cmd

import (
	_ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"
	_ "github.com/gogf/gf/contrib/drivers/mssql/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/oracle/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/vera-byte/vgo/vgo-tools/internal/cmd/gengatewaycurd"
	"github.com/vera-byte/vgo/vgo-tools/internal/cmd/gengrpccurd"
)

type (
	cGenGrpcCrud    = gengrpccurd.CGenGrpcCrud
	cGenGatewayCrud = gengatewaycurd.CGenGatewayCrud
)
