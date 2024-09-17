package model

import (
	"digital-wallet/pkg/types"
	"time"
)

const (
	AuditOperationCreate = "CREATE"
	AuditOperationUpdate = "UPDATE"
	AuditOperationDelete = "DELETE"
)

type Audit struct {
	ID        string  `gorm:"column:id;primaryKey;default:uuid_generate_v4()" json:"id"`
	Operation string  `gorm:"column:operation" json:"operation"`
	Table     string  `gorm:"column:table_name" json:"table"`
	RecordID  string  `gorm:"column:record_id" json:"recordId"`
	Actor     string  `gorm:"column:actor" json:"actor"`
	ActorID   string  `gorm:"column:actor_id" json:"actorId"`
	Remarks   *string `gorm:"column:remarks" json:"remarks"`
	// @swaggertype object
	OldRecord *types.JSONB `gorm:"column:old_record" json:"oldRecord"`
	// @swaggertype object
	NewRecord *types.JSONB `gorm:"column:new_record" json:"newRecord"`
	CreatedAt time.Time    `gorm:"column:created_at" json:"createdAt"`
}

func (Audit) TableName() string {
	return "audit"
}
