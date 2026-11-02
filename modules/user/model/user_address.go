package model

import "github.com/vera-byte/vgo/v"

const TableNameUserAddress = "user_address"

// UserAddress 收货地址（参照 cool-admin-midway user 模块 UserAddressEntity）
type UserAddress struct {
	*v.Model
	UserID   uint   `json:"userId"`   // 用户ID
	Contact  string `json:"contact"`  // 联系人
	Phone    string `json:"phone"`    // 手机号
	Province string `json:"province"` // 省
	City     string `json:"city"`     // 市
	District string `json:"district"` // 区
	Address  string `json:"address"`  // 地址
	IsDefault bool  `json:"isDefault"` // 是否默认
}

// TableName UserAddress 的表名
func (*UserAddress) TableName() string {
	return TableNameUserAddress
}

// NewUserAddress 创建一个 UserAddress 模型实例
func NewUserAddress() *UserAddress {
	return &UserAddress{Model: v.NewModel()}
}
