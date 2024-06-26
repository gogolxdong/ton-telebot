// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameJetton = "jetton"

// Jetton jetton表
type Jetton struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement:true;comment:id" json:"id"`                            // id
	Address    string    `gorm:"column:address;not null;comment:jetton地址" json:"address"`                                 // jetton地址
	CreateTime time.Time `gorm:"column:create_time;not null;default:current_timestamp();comment:创建时间" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;not null;default:current_timestamp();comment:更新时间" json:"update_time"` // 更新时间
}

// TableName Jetton's table name
func (*Jetton) TableName() string {
	return TableNameJetton
}
