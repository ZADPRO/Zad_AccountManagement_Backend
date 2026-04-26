package internal

import "time"

type AuditModel struct {
	CreatedAt time.Time  `json:"createdAt" gorm:"column:createdat"`
	CreatedBy int32      `json:"createdBy" gorm:"column:createdby"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updatedat"`
	UpdatedBy int32      `json:"updatedBy" gorm:"column:updatedby"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"column:deletedat"`
	DeletedBy *string    `json:"deletedBy,omitempty" gorm:"column:deletedby"`
}