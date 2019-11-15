package subscriptions

import (
	"github.com/gofrs/uuid"
)

// Subscription represents an entry in a subscriber list.
type Subscription struct {
	PK           uuid.UUID              `sql:"pk,pk"`
	SubscriberPK uuid.UUID              `sql:"subscriber_pk"`
	ListPK       uuid.UUID              `sql:"list_pk"`
	EmailAddress string                 `sql:"email_address"`
	Data         map[string]interface{} `sql:"data"`
	Version      uint32                 `sql:"version"`

	tableName struct{} `sql:"subscriptions"`
}
