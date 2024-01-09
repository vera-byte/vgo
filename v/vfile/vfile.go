package vfile

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/v/vconfig"
)

type Driver interface {
	New() Driver
	GetMode() (data interface{}, err error)
	Upload(ctx g.Ctx) (string, error)
}

var (
	// FileMap is the map for registered file drivers.
	FileMap = map[string]Driver{}
)

func NewFile() (d Driver) {
	if driver, ok := FileMap[vconfig.Config.File.Mode]; ok {
		return driver.New()
	}
	errorMsg := "\n"
	errorMsg += `无法找到指定文件上传类型 "%s"`
	errorMsg += `，您是否拼写错误了类型名称 "%s" 或者忘记导入上传支持包？`
	errorMsg += `参考:https://github.com/vera-byte/vgo/tree/master/contrib/files`
	err := gerror.Newf(errorMsg, vconfig.Config.File.Mode, vconfig.Config.File.Mode)

	panic(err)

}

// Register registers custom file driver to v.
func Register(name string, driver Driver) error {
	FileMap[name] = driver
	return nil
}

// func init() {
// 	// Register("local", &Local{})
// 	// Register("oss", &Oss{})
// 	file, err := NewFile()
// 	if err != nil {
// 		panic(err)
// 	}
// 	File = file
// }
