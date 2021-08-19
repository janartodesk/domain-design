package model

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/janartodesk/domain-design/lists/domain"
	"github.com/uptrace/bun"
)

// Subscription is a database model for a list subscription.
type Subscription struct {
	PK           uuid.UUID              `bun:"pk,pk"`
	ListPK       uuid.UUID              `bun:"list_pk,pk"`
	SubscriberPK uuid.UUID              `bun:"subscriber_pk"`
	EmailAddress string                 `bun:"email"`
	Data         map[string]interface{} `bun:"data"`
	Version      uint32                 `bun:"version"`

	bun.BaseModel `bun:"subscriptions"`
}

// CreateSubscription creates a subscription.
func CreateSubscription(db bun.IDB, subscription *domain.Subscription) error {
	if _, err := db.NewInsert().Model(&Subscription{
		PK:           subscription.PK,
		EmailAddress: string(subscription.EmailAddress),
		Version:      subscription.Version,
	}).Exec(context.Background()); err != nil {
		return err
	}

	return nil
}

// UpdateSubscription updates a subscription.
func UpdateSubscription(db bun.IDB, subscription *domain.Subscription) error {
	res, err := db.NewUpdate().Model(&Subscription{
		PK:           subscription.PK,
		EmailAddress: string(subscription.EmailAddress),
		Version:      subscription.Version,
	}).Where(
		"pk = ? AND version = ?",
		subscription.PK,
		subscription.Version-1,
	).Exec(context.Background())

	if err != nil {
		return err
	}

	if c, err := res.RowsAffected(); err != nil {
		return err
	} else if c == 0 {
		return errors.New("precondition failed")
	}

	return nil
}

// GetSubscription returns a subscription.
func GetSubscription(db bun.IDB, listPK, pk uuid.UUID) (*domain.Subscription, error) {
	model := Subscription{
		PK:     pk,
		ListPK: listPK,
	}

	if err := db.NewSelect().Model(&model).WherePK().Scan(context.Background()); err != nil {
		return nil, err
	}

	return &domain.Subscription{
		PK:           model.PK,
		EmailAddress: domain.EmailAddress(model.EmailAddress),
		Version:      model.Version,
	}, nil
}

// GetSubscriptionForSubscriber returns a subscription for a subscriber in a list.
func GetSubscriptionForSubscriber(db bun.IDB, listPK, subscriberPK uuid.UUID) (*domain.Subscription, error) {
	model := Subscription{}

	if err := db.NewSelect().Model(&model).Where(
		"list_pk = ? AND subscriber_pk = ?",
		listPK,
		subscriberPK,
	).Scan(context.Background()); err != nil {
		return nil, err
	}

	return &domain.Subscription{
		PK:           model.PK,
		EmailAddress: domain.EmailAddress(model.EmailAddress),
		Version:      model.Version,
	}, nil
}

// ListSubscriptions returns a list of subscriptions.
func ListSubscriptions(db bun.IDB, offset, limit uint32) ([]*domain.Subscription, error) {
	model := []Subscription{}

	if err := db.NewSelect().Model(&model).Scan(context.Background()); err != nil {
		return nil, err
	}

	res := []*domain.Subscription{}

	for _, subscription := range model {
		res = append(res, &domain.Subscription{
			PK:           subscription.PK,
			SubscriberPK: subscription.SubscriberPK,
			ListPK:       subscription.ListPK,
			EmailAddress: domain.EmailAddress(subscription.EmailAddress),
			Data:         subscription.Data,
			Version:      subscription.Version,
		})
	}

	return res, nil
}
