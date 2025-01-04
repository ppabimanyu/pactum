package pactum

import (
	"errors"
	"gorm.io/gorm"
)

func _update(tx *gorm.DB, data interface{}) error {
	return tx.Model(data).Select("*").Updates(data).Error
}

func _delete(tx *gorm.DB, data interface{}) error {
	return tx.Delete(data).Error
}

func _findApprovalByID(tx *gorm.DB, id uint64) (*ApprovalModel, error) {
	var p ApprovalModel
	if err := tx.Where("id = ?", id).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidApprovalData
		}
		return nil, err
	}
	return &p, nil
}
