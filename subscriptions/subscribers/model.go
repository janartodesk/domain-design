package subscribers

import (
	"github.com/gofrs/uuid"
)

// Subscriber represents a business entity consented to receive content from us.
type Subscriber struct {
	PK             uuid.UUID `sql:"pk,pk"`
	OrganizationPK uuid.UUID `sql:"organization_pk"`
	EmailAddress   string    `sql:"email_address"`
	IsForgotten    bool      `sql:"is_forgotten"`

	tableName struct{} `sql:"subscribers"`
}
