package model

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/janartodesk/domain-design/lists/domain"
	"github.com/uptrace/bun"
)

// Subscriber is a database model for a list subscriber.
type Subscriber struct {
	PK           uuid.UUID `bun:"pk,pk"`
	EmailAddress string    `bun:"email"`
	Version      uint32    `bun:"version"`

	bun.BaseModel `bun:"subscribers"`
}

// CreateSubscriber creates a subscriber.
func CreateSubscriber(db bun.IDB, subscriber *domain.Subscriber) error {
	if _, err := db.NewInsert().Model(&Subscriber{
		PK:           subscriber.PK,
		EmailAddress: string(subscriber.EmailAddress),
		Version:      subscriber.Version,
	}).Exec(context.Background()); err != nil {
		return err
	}

	return nil
}

// UpdateSubscriber updates a subscriber.
func UpdateSubscriber(db bun.IDB, subscriber *domain.Subscriber) error {
	res, err := db.NewUpdate().Model(&Subscriber{
		PK:           subscriber.PK,
		EmailAddress: string(subscriber.EmailAddress),
		Version:      subscriber.Version,
	}).Where(
		"pk = ? AND version = ?",
		subscriber.PK,
		subscriber.Version-1,
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

// GetSubscriber returns a subscriber.
func GetSubscriber(db bun.IDB, pk uuid.UUID) (*domain.Subscriber, error) {
	model := Subscriber{
		PK: pk,
	}

	if err := db.NewSelect().Model(&model).WherePK().Scan(context.Background()); err != nil {
		return nil, err
	}

	return &domain.Subscriber{
		PK:           model.PK,
		EmailAddress: domain.EmailAddress(model.EmailAddress),
		Version:      model.Version,
	}, nil
}

// GetSubscriberByEmailAddress returns a subscriber by their email address.
func GetSubscriberByEmailAddress(db bun.IDB, addr domain.EmailAddress) (*domain.Subscriber, error) {
	model := Subscriber{}

	if err := db.NewSelect().Model(&model).Where(
		"email_addr = ?",
		addr,
	).Scan(context.Background()); err != nil {
		return nil, err
	}

	return &domain.Subscriber{
		PK:           model.PK,
		EmailAddress: domain.EmailAddress(model.EmailAddress),
		Version:      model.Version,
	}, nil
}

// ListSubscribers returns a list of subscribers.
func ListSubscribers(db bun.IDB, offset, limit uint32) ([]*domain.Subscriber, error) {
	model := []Subscriber{}

	if err := db.NewSelect().Model(&model).Scan(context.Background()); err != nil {
		return nil, err
	}

	res := []*domain.Subscriber{}

	for _, subscriber := range model {
		res = append(res, &domain.Subscriber{
			PK:           subscriber.PK,
			EmailAddress: domain.EmailAddress(subscriber.EmailAddress),
			Version:      subscriber.Version,
		})
	}

	return res, nil
}
