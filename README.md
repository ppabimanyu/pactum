![logo.png](logo.png)

# Project Description

## Overview

This project is a Go-based application that manages approval processes. It includes functionalities for creating, updating, and deleting approval records, as well as handling the approval and rejection of these records.

## Key Components

### Models

- **ApprovalModel**: Represents the structure of an approval record in the database. It contains information such as the reference code, tag, approval type, approval object, request details, and approval details.
- **AuditModel**: Represents the structure of an audit record in the database. It contains information such as the approval ID, reference code, tag, approval type, approval status, approval object, request details, and approval details.

### Constants

- **Approval Types**:
    - `ApprovalTypeCreate`: Represents the approval type for creating an object.
    - `ApprovalTypeUpdate`: Represents the approval type for updating an object.
    - `ApprovalTypeDelete`: Represents the approval type for deleting an object.

### Functions

- **ApprovalCreate**: Creates an approval record for creating an object.
- **ApprovalUpdate**: Creates an approval record for updating an object.
- **ApprovalDelete**: Creates an approval record for deleting an object.
- **Approve**: Handles the approval of an approval object.
- **Reject**: Handles the rejection of an approval object.

### Error Handling

- **ErrInvalidApprovalType**: Error for invalid approval type.
- **ErrInvalidApprovalData**: Error for invalid approval data.
- **ErrStatusNotPending**: Error for status not pending.

## Dependencies

- **Go**: The primary programming language used.
- **GORM**: The ORM library used for database interactions.
- **JSON**: Used for encoding and decoding approval objects.

## Environment Variables

- **DB\_PREFIX**: The prefix for database table names.

## Usage

To use this project, ensure you have Go installed and set up the necessary environment variables. The project can be run using standard Go commands.

```go
package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize the database connection
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	// Example usage of ApprovalCreate
	refCode := "REF123"
	tag := "TAG123"
	approvalObject := &ApprovalModel{ /* initialize fields */ }
	reqBy := "user1"
	reqHost := "localhost"
	reqNote := "Initial approval"

	approvalID, err := ApprovalCreate(db, refCode, tag, approvalObject, &reqBy, &reqHost, &reqNote)
	if err != nil {
		fmt.Println("Failed to create approval:", err)
		return
	}
	fmt.Println("Approval created with ID:", approvalID)

	// Example usage of ApprovalUpdate
	approvalObject = &ApprovalModel{ /* initialize fields */ }
	approvalID, err = ApprovalUpdate(db, refCode, tag, approvalObject, &reqBy, &reqHost, &reqNote)
	if err != nil {
		fmt.Println("Failed to update approval:", err)
		return
	}
	fmt.Println("Approval updated with ID:", approvalID)

	// Example usage of ApprovalDelete
	approvalObject = &ApprovalModel{ /* initialize fields */ }
	approvalID, err = ApprovalDelete(db, refCode, tag, approvalObject, &reqBy, &reqHost, &reqNote)
	if err != nil {
		fmt.Println("Failed to delete approval:", err)
		return
	}
	fmt.Println("Approval deleted with ID:", approvalID)

	// Example usage of Approve
	err = Approve(db, approvalObject, &reqBy, &reqHost, &reqNote)
	if err != nil {
		fmt.Println("Failed to approve:", err)
		return
	}
	fmt.Println("Approval approved")

	// Example usage of Reject
	err = Reject(db, approvalObject)
	if err != nil {
		fmt.Println("Failed to reject:", err)
		return
	}
	fmt.Println("Approval rejected")
}
```

## License

This project is licensed under the MIT License.