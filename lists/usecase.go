package lists

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/janartodesk/domain-design/lists/domain"
	"github.com/janartodesk/domain-design/lists/model"
	"github.com/janartodesk/domain-design/pkg/db"
	"github.com/uptrace/bun"
	"google.golang.org/protobuf/proto"
)

type EventPublisher interface {
	Publish(proto.Message) error
}

type Usecase struct {
	db     *bun.DB
	events EventPublisher
}

// CreateList creates a subscriber list.
func (u *Usecase) CreateList(title string) (*domain.List, error) {
	list, err := domain.CreateList(title)
	if err != nil {
		return nil, err
	}

	if err := model.CreateList(u.db, list); err != nil {
		return nil, err
	}

	return list, nil
}

// RenameList renames a subscriber list.
func (u *Usecase) RenameList(listPK uuid.UUID, title string) (*domain.List, error) {
	list, err := model.GetList(u.db, listPK)
	if err != nil {
		return nil, err
	}

	list, err = domain.RenameList(*list, title)
	if err != nil {
		return nil, err
	}

	if err := model.UpdateList(u.db, list); err != nil {
		return nil, err
	}

	return list, nil
}

// DeleteList deletes a subscriber list.
func (u *Usecase) DeleteList(listPK uuid.UUID) error {
	if err := model.DeleteList(u.db, listPK); err != nil {
		return err
	}

	return u.events.Publish(&ListDeleted{
		ListPK: listPK.Bytes(),
	})
}

// SubscribeSubscriber subscribes a subscriber into a list.
func (u *Usecase) SubscribeSubscriber(listPK uuid.UUID, emailAddr domain.EmailAddress, data domain.SubscriptionData) error {
	return db.WithTransaction(u.db, func(tx bun.Tx) error {
		if _, err := u.createSubscription(tx, listPK, emailAddr, data); err != nil {
			return err
		}

		return nil
	})
}

// UnsubscribeSubscriber unsubscribes a subscriber from a list.
func (u *Usecase) Unsubscribe(listPK, subscriptionPK uuid.UUID) error {
	return db.WithTransaction(u.db, func(tx bun.Tx) error {
		subscription, err := model.GetSubscription(tx, listPK, subscriptionPK)
		if err != nil {
			return err
		}

		if _, err := u.cancelSubscription(tx, subscription); err != nil {
			return err
		}

		return nil
	})
}

// OptInSubscriber opts a subscriber into a list.
func (u *Usecase) OptInSubscriber(listPK uuid.UUID, emailAddr domain.EmailAddress, data domain.SubscriptionData) error {
	return db.WithTransaction(u.db, func(tx bun.Tx) error {
		subscription, err := u.createSubscription(tx, listPK, emailAddr, data)
		if err != nil {
			return err
		}

		return u.events.Publish(&SubscriberOptedIn{
			SubscriberPK: subscription.SubscriberPK.Bytes(),
			ListPK:       subscription.ListPK.Bytes(),
		})
	})
}

// OptOutSubscriber opts a subscriber out from a list.
func (u *Usecase) OptOutSubscriber(listPK, subscriberPK uuid.UUID) error {
	return db.WithTransaction(u.db, func(tx bun.Tx) error {
		subscription, err := model.GetSubscriptionForSubscriber(tx, listPK, subscriberPK)
		if err != nil {
			return err
		}

		subscription, err = u.cancelSubscription(tx, subscription)
		if err != nil {
			return err
		}

		return u.events.Publish(&SubscriberOptedOut{
			SubscriberPK: subscription.SubscriberPK.Bytes(),
			ListPK:       subscription.ListPK.Bytes(),
		})
	})
}

func (u *Usecase) createSubscription(tx bun.IDB, listPK uuid.UUID, emailAddr domain.EmailAddress, data domain.SubscriptionData) (*domain.Subscription, error) {
	if err := emailAddr.Validate(); err != nil {
		return nil, err
	}

	list, err := model.GetList(tx, listPK)
	if err != nil {
		return nil, err
	}

	subscriber, err := model.GetSubscriberByEmailAddress(tx, emailAddr)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			subscriber, err := domain.CreateSubscriber(emailAddr)
			if err != nil {
				return nil, err
			}

			if err := model.CreateSubscriber(tx, subscriber); err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}

	subscription, err := domain.CreateSubscription(*subscriber, *list, data)
	if err != nil {
		return nil, err
	}

	if err := model.CreateSubscription(tx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (u *Usecase) cancelSubscription(tx bun.IDB, subscription *domain.Subscription) (*domain.Subscription, error) {
	subscription, err := domain.CancelSubscription(*subscription)
	if err != nil {
		return nil, err
	}

	if err := model.UpdateSubscription(tx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}
