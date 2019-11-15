package subscribers

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-pg/pg/orm"
)

// Create a subscriber.
func Create(db orm.DB, subscriber *Subscriber) error {
	subscriber.IsForgotten = false

	if err := normalize(subscriber); err != nil {
		return err
	} else if err := validate(subscriber); err != nil {
		return err
	}

	if _, err := db.Model(subscriber).Insert(); err != nil {
		return err
	}

	return nil
}

// Forget the subscriber.
func Forget(db orm.DB, subscriber *Subscriber) error {
	if subscriber.IsForgotten {
		return errors.New("subscriber has already been forgotten")
	}

	subscriber.EmailAddress = "forgotten@smaily.com"
	subscriber.IsForgotten = true

	if err := validate(subscriber); err != nil {
		return err
	}

	if r, err := db.Model(subscriber).Where("pk = ?pk").Update(); err != nil {
		return err
	} else if r.RowsAffected() == 0 {
		return errors.New("subscriber was not found")
	}

	return nil
}

func normalize(subscriber *Subscriber) error {
	subscriber.EmailAddress = strings.TrimSpace(subscriber.EmailAddress)
	subscriber.EmailAddress = strings.ToLower(subscriber.EmailAddress)

	return nil
}

func validate(subscriber *Subscriber) error {
	return validation.ValidateStruct(
		subscriber,
		validation.Field(&subscriber.PK, validation.Required),
		validation.Field(&subscriber.OrganizationPK, validation.Required),
		validation.Field(&subscriber.EmailAddress, validation.Required, is.Email),
		validation.Field(&subscriber.IsForgotten, validation.Required),
	)
}
