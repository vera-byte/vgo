// Package packed 是打包资源的占位包。
// vgo-tools build 会向本目录生成 build_pack_data.go（gres 资源注册），
// 空实现保证未执行打包时 blank-import 仍然可编译。
package packed
