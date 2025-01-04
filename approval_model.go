package pactum

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"os"
	"time"
)

var (
	ErrInvalidApprovalType = errors.New("invalid approval type")
	ErrInvalidApprovalData = errors.New("invalid approval data")
	ErrStatusNotPending    = errors.New("status not pending")
)

const (
	_approvalTableName = "t_approval"
)

// ApprovalModel represents the structure of an approval record in the database.
// It contains information about the approval, such as the reference code, tag,
// approval type, approval object, request details, and approval details.
type ApprovalModel struct {
	ID             uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RefCode        string         `gorm:"column:ref_code;size:36;not null" json:"ref_code"`
	Tag            string         `gorm:"column:tag;size:24;not null" json:"tag"`
	ApprovalType   ApprovalType   `gorm:"column:approval_type;size:6;not null" json:"approval_type"`
	ApprovalStatus ApprovalStatus `gorm:"column:approval_status;size:8;not null" json:"approval_status"`

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

func (p *ApprovalModel) TableName() string {
	return os.Getenv("DB_PREFIX") + _approvalTableName
}

// ParseApprovalObject unmarshal the JSON-encoded approval object into the provided object.
// Parameters:
//   - obj: the object into which the approval object will be unmarshalled.
//
// Returns:
//   - error: an error if the unmarshalling fails, otherwise nil.
func (p *ApprovalModel) ParseApprovalObject(obj any) error {
	if p != nil {
		return json.Unmarshal(p.ApprovalObject, obj)
	}
	return ErrInvalidApprovalData
}

func (p *ApprovalModel) BeforeCreate(_ *gorm.DB) error {
	if !p.ApprovalType.IsValid() {
		return ErrInvalidApprovalType
	}
	if p.ReqAt == nil {
		now := time.Now()
		p.ReqAt = &now
	}
	if p.ReqBy == nil {
		system := DefaultUser
		p.ReqBy = &system
	}
	if p.ReqHost == nil {
		host := DefaultHost
		p.ReqHost = &host
	}
	return nil
}

func (p *ApprovalModel) _createAudit(tx *gorm.DB) error {
	if p != nil {
		audit := AuditModel{
			ApprovalID:     p.ID,
			RefCode:        p.RefCode,
			Tag:            p.Tag,
			ApprovalType:   p.ApprovalType,
			ApprovalStatus: p.ApprovalStatus,
			ApprovalObject: p.ApprovalObject,
			ReqBy:          p.ReqBy,
			ReqHost:        p.ReqHost,
			ReqAt:          p.ReqAt,
			ReqNote:        p.ReqNote,
			ApprovalBy:     p.ApprovalBy,
			ApprovalHost:   p.ApprovalHost,
			ApprovalAt:     p.ApprovalAt,
			ApprovalNote:   p.ApprovalNote,
		}
		return tx.Create(&audit).Error
	}
	return ErrInvalidApprovalData
}

func (p *ApprovalModel) _createApprovalLog(tx *gorm.DB) error {
	if p != nil {
		log := ApprovalLogModel{
			ApprovalID:     p.ID,
			RefCode:        p.RefCode,
			Tag:            p.Tag,
			ApprovalType:   p.ApprovalType,
			ApprovalStatus: p.ApprovalStatus,
			ApprovalObject: p.ApprovalObject,
			ReqBy:          p.ReqBy,
			ReqHost:        p.ReqHost,
			ReqAt:          p.ReqAt,
			ReqNote:        p.ReqNote,
			ApprovalBy:     p.ApprovalBy,
			ApprovalHost:   p.ApprovalHost,
			ApprovalAt:     p.ApprovalAt,
			ApprovalNote:   p.ApprovalNote,
		}
		return tx.Create(&log).Error
	}
	return ErrInvalidApprovalData
}

func (p *ApprovalModel) _delete(tx *gorm.DB) error {
	if p != nil {
		return tx.Delete(p).Error
	}
	return ErrInvalidApprovalData
}

func (p *ApprovalModel) IsPending() bool {
	if p != nil {
		return p.ApprovalStatus == ApprovalStatusPending
	}
	return false
}
