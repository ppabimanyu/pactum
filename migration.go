package pactum

import (
	"fmt"
	"gorm.io/gorm"
	"log/slog"
)

// MigrationDB performs the database migration for the specified models.
// It uses GORM's AutoMigrate function to create or update the tables for the models.
// If the migration fails for any model, it logs the error and panics.
//
// Parameters:
// - tx (*gorm.DB): The GORM database transaction object.
func MigrationDB(tx *gorm.DB) {
	approvalModelTable := new(ApprovalModel).TableName()
	if err := tx.AutoMigrate(&ApprovalModel{}); err != nil {
		slog.Error(fmt.Sprintf("Failed to migrate table %s: %v", approvalModelTable, err))
		panic(err)
	} else {
		slog.Info(fmt.Sprintf("Table %s migrated successfully", approvalModelTable))
	}

	approvalLogModelTable := new(ApprovalLogModel).TableName()
	if err := tx.AutoMigrate(&ApprovalLogModel{}); err != nil {
		slog.Error(fmt.Sprintf("Failed to migrate table %s: %v", approvalLogModelTable, err))
		panic(err)
	} else {
		slog.Info(fmt.Sprintf("Table %s migrated successfully", approvalLogModelTable))
	}

	auditModelTable := new(AuditModel).TableName()
	if err := tx.AutoMigrate(&AuditModel{}); err != nil {
		slog.Error(fmt.Sprintf("Failed to migrate table %s: %v", auditModelTable, err))
		panic(err)
	} else {
		slog.Info(fmt.Sprintf("Table %s migrated successfully", auditModelTable))
	}
}
