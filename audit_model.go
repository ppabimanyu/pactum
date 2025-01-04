package pactum

import (
	"encoding/json"
	"os"
	"time"
)

const (
	_auditTableName = "t_audit"
)

type AuditModel struct {
	ID             uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApprovalID     uint64         `gorm:"column:approval_id;not null" json:"approval_id"`
	RefCode        string         `gorm:"column:ref_code;size:36;not null" json:"ref_code"`
	Tag            string         `gorm:"column:tag;size:24;not null" json:"tag"`
	ApprovalType   ApprovalType   `gorm:"column:approval_type;size:6;not null" json:"approval_type"`
	ApprovalStatus ApprovalStatus `gorm:"column:approval_status;size:9;not null" json:"approval_status"`

	ApprovalObject json.RawMessage `gorm:"type:text;serializer:json" json:"-"`

	ReqBy   *string    `gorm:"column:request_by;size:100" json:"request_by"`
	ReqHost *string    `gorm:"column:request_host;size:100" json:"request_host"`
	ReqAt   *time.Time `gorm:"column:request_at;size:100" json:"request_at"`
	ReqNote *string    `gorm:"column:request_note" json:"request_note"`

	ApprovalBy   *string    `gorm:"column:approval_by;size:100" json:"approval_by"`
	ApprovalHost *string    `gorm:"column:approval_host;size:100" json:"approval_host"`
	ApprovalAt   *time.Time `gorm:"column:approval_at;size:100" json:"approval_at"`
	ApprovalNote *string    `gorm:"column:approval_note" json:"approval_note"`
}

func (a *AuditModel) TableName() string {
	return os.Getenv("DB_PREFIX") + _auditTableName
}
