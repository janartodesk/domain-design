package subscriptions

import (
	"errors"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-pg/pg/orm"
)

// Create a subscription.
func Create(db orm.DB, subscription *Subscription) error {
	subscription.Version = 1

	if err := normalize(subscription); err != nil {
		return err
	} else if err := validate(subscription); err != nil {
		return err
	}

	if _, err := db.Model(subscription).Insert(); err != nil {
		return err
	}

	return nil
}

// Cancel a subscription.
func Cancel(db orm.DB, subscription *Subscription) error {
	if r, err := db.Model(subscription).Where("subscriber_pk = ?subscriber_pk AND list_pk = ?list_pk AND version = ?version").Delete(); err != nil {
		return err
	} else if r.RowsAffected() == 0 {
		return errors.New("subscription was not found")
	}

	return nil
}

// Update the subscription with new data.
func Update(db orm.DB, subscription *Subscription, data map[string]interface{}) error {
	if reflect.DeepEqual(subscription.Data, data) {
		return errors.New("data did not change")
	}

	subscription.Data = data
	subscription.Version++

	if err := validate(subscription); err != nil {
		return err
	}

	if r, err := db.Model(subscription).Where("subscriber_pk = ?subscriber_pk AND list_pk = ?list_pk AND version = ?version").Update(); err != nil {
		return err
	} else if r.RowsAffected() == 0 {
		return errors.New("subscription was not found")
	}

	return nil
}

// PartialUpdate updates the subscription data partially.
func PartialUpdate(db orm.DB, subscription *Subscription, data map[string]interface{}) error {
	updated := make(map[string]interface{})

	// Create a copy of the subscription's data.
	for k, v := range subscription.Data {
		updated[k] = v
	}

	// Then overlay the partial update data onto existing data.
	for k, v := range data {
		updated[k] = v
	}

	return Update(db, subscription, updated)
}

func normalize(subscription *Subscription) error {
	subscription.EmailAddress = strings.TrimSpace(subscription.EmailAddress)

	return nil
}

func validate(subscription *Subscription) error {
	return validation.ValidateStruct(
		subscription,
		validation.Field(&subscription.PK, validation.Required),
		validation.Field(&subscription.SubscriberPK, validation.Required),
		validation.Field(&subscription.ListPK, validation.Required),
		validation.Field(&subscription.EmailAddress, validation.Required, is.Email),
		validation.Field(&subscription.Version, validation.Required),
	)
}
