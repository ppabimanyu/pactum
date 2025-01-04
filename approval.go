package pactum

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"reflect"
)

func _createApproval(db *gorm.DB, refCode string, approvalType ApprovalType, approvalObject interface{}, tag string, reqBy, reqHost, reqNote *string) (uint64, error) {
	approvalObjectBytes, err := json.Marshal(approvalObject)
	if err != nil {
		return 0, err
	}
	approval := &ApprovalModel{
		Tag:            tag,
		RefCode:        refCode,
		ApprovalType:   approvalType,
		ApprovalObject: approvalObjectBytes,
		ReqBy:          reqBy,
		ReqHost:        reqHost,
		ReqNote:        reqNote,
	}
	if err := db.Create(approval).Error; err != nil {
		return 0, err
	}
	return approval.ID, nil
}

func _getApprovalID(obj any) (*uint64, error) {
	// Get the type and value of the struct
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("expected a pointer to a struct")
	}

	// Dereference the pointer to get the struct
	val = val.Elem()
	typ = typ.Elem()

	var approvalIDField bool
	var approvalIDValue *uint64

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()

		if field.Name == "ApprovalID" && field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint64 {
			approvalIDField = true
			approvalID := value.(*uint64)
			approvalIDValue = approvalID
		}
	}

	if !approvalIDField {
		return nil, errors.New("field ApprovalID not found")
	}

	return approvalIDValue, nil
}

func _unsetApprovalID(obj any) {
	// Get the type and value of the struct
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return
	}

	// Dereference the pointer to get the struct
	val = val.Elem()
	typ = typ.Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.Name == "ApprovalID" && field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint64 {
			val.Field(i).Set(reflect.Zero(field.Type))
		}

		if field.Name == "Approval" && field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			val.Field(i).Set(reflect.Zero(field.Type))
		}
	}
}

// ApprovalUpdate creates an approval record for updating an object.
// It marshals the approval object to JSON, creates an approval record in the database,
// and returns the approval ID and any error encountered.
// Parameters:
//   - tx: the database transaction.
//   - refCode: the reference code for the approval.
//   - tag: the tag associated with the approval.
//   - approvalObject: the object to be approved.
//   - reqBy: the user who requested the approval.
//   - reqHost: the host from which the request was made.
//   - reqNote: additional notes for the request.
//
// Returns:
//   - uint64: the ID of the created approval record.
//   - error: an error if the operation fails, otherwise nil.
func ApprovalUpdate(tx *gorm.DB, refCode string, tag string, approvalObject interface{}, reqBy, reqHost, reqNote *string) (uint64, error) {
	return _createApproval(tx, refCode, ApprovalTypeUpdate, approvalObject, tag, reqBy, reqHost, reqNote)
}

// ApprovalCreate creates an approval record for creating an object.
// It marshals the approval object to JSON, creates an approval record in the database,
// and returns the approval ID and any error encountered.
// Parameters:
//   - tx: the database transaction.
//   - refCode: the reference code for the approval.
//   - tag: the tag associated with the approval.
//   - approvalObject: the object to be approved.
//   - reqBy: the user who requested the approval.
//   - reqHost: the host from which the request was made.
//   - reqNote: additional notes for the request.
//
// Returns:
//   - uint64: the ID of the created approval record.
//   - error: an error if the operation fails, otherwise nil.
func ApprovalCreate(tx *gorm.DB, refCode string, tag string, approvalObject interface{}, reqBy, reqHost, reqNote *string) (uint64, error) {
	return _createApproval(tx, refCode, ApprovalTypeCreate, approvalObject, tag, reqBy, reqHost, reqNote)
}

// ApprovalDelete creates an approval record for deleting an object.
// It marshals the approval object to JSON, creates an approval record in the database,
// and returns the approval ID and any error encountered.
// Parameters:
//   - tx: the database transaction.
//   - refCode: the reference code for the approval.
//   - tag: the tag associated with the approval.
//   - approvalObject: the object to be approved.
//   - reqBy: the user who requested the approval.
//   - reqHost: the host from which the request was made.
//   - reqNote: additional notes for the request.
//
// Returns:
//   - uint64: the ID of the created approval record.
//   - error: an error if the operation fails, otherwise nil.
func ApprovalDelete(tx *gorm.DB, refCode string, tag string, approvalObject interface{}, reqBy, reqHost, reqNote *string) (uint64, error) {
	return _createApproval(tx, refCode, ApprovalTypeDelete, approvalObject, tag, reqBy, reqHost, reqNote)
}

// Approve handles the approval of an approval object.
// It retrieves the ApprovalID from the object, finds the corresponding approval record,
// parses the approval object, unsets the ApprovalID from the object, updates or deletes
// the approval object based on the approval type, adds metadata to the approval record,
// creates an audit record, and deletes the approval record.
// Parameters:
//   - tx: the database transaction.
//   - approvalObject: the object containing the approval information.
//   - approvalBy: the user who approved the object.
//   - approvalHost: the host from which the approval was made.
//   - approvalNote: additional notes for the approval.
//
// Returns:
//   - error: an error if the operation fails, otherwise nil.
func Approve(tx *gorm.DB, approvalObject any, approvalBy, approvalHost, approvalNote *string) error {
	// Get ApprovalID from the object
	approvalID, err := _getApprovalID(approvalObject)
	if err != nil {
		return err
	}
	if approvalID == nil {
		return ErrStatusNotPending
	}

	// Find the approval record
	p, err := _findApprovalByID(tx, *approvalID)
	if err != nil {
		return err
	}

	// Parse the approval object
	if err := p.ParseApprovalObject(approvalObject); err != nil {
		return err
	}

	// Unset the ApprovalID from the object
	_unsetApprovalID(approvalObject)

	switch p.ApprovalType {
	// Update the approval object based on the approval type
	case ApprovalTypeCreate, ApprovalTypeUpdate:
		if err := _update(tx, approvalObject); err != nil {
			return err
		}
	// Delete the approval record
	case ApprovalTypeDelete:
		if err := _delete(tx, approvalObject); err != nil {
			return err
		}
	}

	// Add metadata to the approval record
	p.ApprovalBy = approvalBy
	p.ApprovalHost = approvalHost
	p.ApprovalNote = approvalNote

	// Create audit record
	if err := p._createAudit(tx); err != nil {
		return err
	}

	// Delete the approval record
	if err := p._delete(tx); err != nil {
		return err
	}

	return nil
}

// Reject handles the rejection of an approval object.
// It retrieves the ApprovalID from the object, finds the corresponding approval record,
// unsets the ApprovalID from the object, updates the object, and deletes the approval record.
// Parameters:
//   - tx: the database transaction.
//   - approvalObject: the object containing the approval information.
//
// Returns:
//   - error: an error if the operation fails, otherwise nil.
func Reject(tx *gorm.DB, approvalObject any) error {
	// Get ApprovalID from the object
	approvalID, err := _getApprovalID(approvalObject)
	if err != nil {
		return err
	}
	if approvalID == nil {
		return ErrStatusNotPending
	}

	// Find the approval record
	p, err := _findApprovalByID(tx, *approvalID)
	if err != nil {
		return err
	}

	// Unset the ApprovalID from the object
	_unsetApprovalID(approvalObject)
	if err := _update(tx, approvalObject); err != nil {
		return err
	}

	// Delete the approval record
	if err := p._delete(tx); err != nil {
		return err
	}

	return nil
}
