package po

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int64     `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// 包含这个类型的字段就是软删除
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
