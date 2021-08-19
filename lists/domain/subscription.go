package domain

import (
	"errors"

	"github.com/gofrs/uuid"
)

// Subscription is a subscriber's subscription to a list.
type Subscription struct {
	PK           uuid.UUID
	SubscriberPK uuid.UUID
	ListPK       uuid.UUID
	EmailAddress EmailAddress
	Data         map[string]interface{}
	IsCancelled  bool
	Version      uint32
}

// CreateSubscription creates a subscription.
func CreateSubscription(Subscriber, List, map[string]interface{}) (*Subscription, error) {
	return nil, nil
}

// CancelSubscription cancels a subscription to a list.
func CancelSubscription(subscription Subscription) (*Subscription, error) {
	if subscription.IsCancelled {
		return nil, errors.New("invariant error")
	}

	subscription.IsCancelled = true
	subscription.Version++

	return &subscription, nil
}
