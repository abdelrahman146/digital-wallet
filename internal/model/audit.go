package model

import "digital-wallet/pkg/types"

const (
	AuditOperationCreate = "CREATE"
	AuditOperationUpdate = "UPDATE"
	AuditOperationDelete = "DELETE"
)

type Audit struct {
	ID        string       `gorm:"column:id;primaryKey;default:uuid_generate_v4()" json:"id"`
	Operation string       `gorm:"column:operation" json:"operation"`
	Table     string       `gorm:"column:table_name" json:"table"`
	RecordID  string       `gorm:"column:record_id" json:"recordId"`
	Actor     string       `gorm:"column:actor" json:"actor"`
	ActorID   *string      `gorm:"column:actor_id" json:"actorId"`
	Remark    *string      `gorm:"column:remarks" json:"remarks"`
	OldRecord *types.JSONB `gorm:"column:old_record" json:"oldRecord"`
	NewRecord types.JSONB  `gorm:"column:new_record" json:"newRecord"`
}

func (Audit) TableName() string {
	return "audit"
}