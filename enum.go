package pactum

type ApprovalType string

const (
	// ApprovalTypeCreate represents the approval type for creating an object.
	ApprovalTypeCreate ApprovalType = "CREATE"
	// ApprovalTypeUpdate represents the approval type for updating an object.
	ApprovalTypeUpdate ApprovalType = "UPDATE"
	// ApprovalTypeDelete represents the approval type for deleting an object.
	ApprovalTypeDelete ApprovalType = "DELETE"
)

// IsValid checks if the ApprovalType is one of the predefined valid types.
// Returns:
//   - bool: true if the ApprovalType is valid, otherwise false.
func (p ApprovalType) IsValid() bool {
	switch p {
	case ApprovalTypeCreate, ApprovalTypeUpdate, ApprovalTypeDelete:
		return true
	}
	return false
}

type ApprovalStatus string

const (
	// ApprovalStatusPending represents the status of an approval that is pending.
	ApprovalStatusPending ApprovalStatus = "PENDING"
	// ApprovalStatusApproved represents the status of an approval that is approved.
	ApprovalStatusApproved ApprovalStatus = "APPROVED"
	// ApprovalStatusRejected represents the status of an approval that is rejected.
	ApprovalStatusRejected ApprovalStatus = "REJECTED"
)

const (
	DefaultHost = "LOCALHOST"
	DefaultUser = "SYSTEM"
)
