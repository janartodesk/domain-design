package subscriptions

import (
	"github.com/go-pg/pg/orm"
	"github.com/gofrs/uuid"
)

// Get a subscription.
func Get(db orm.DB, subscriberPK, listPK uuid.UUID) (*Subscription, error) {
	subscription := &Subscription{}

	if err := db.Model(subscription).Where("subscriber_pk = ? AND list_pk = ?", subscriberPK, listPK).Select(); err != nil {
		return nil, err
	}

	return subscription, nil
}
