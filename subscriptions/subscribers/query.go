package subscribers

import (
	"github.com/go-pg/pg/orm"
	"github.com/gofrs/uuid"
)

// Get gets a subscriber.
func Get(db orm.DB, organizationPK uuid.UUID, emailAddress string) (*Subscriber, error) {
	subscriber := &Subscriber{}

	// TODO: the email address needs normalization.
	if err := db.Model(subscriber).Where("organization_pk = ? AND email_address = ?", organizationPK, emailAddress).Select(); err != nil {
		return nil, err
	}

	return subscriber, nil
}
