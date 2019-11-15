package subscriptionss

import (
	"github.com/go-pg/pg/orm"
	"github.com/gofrs/uuid"

	"domain/subscriptions/lists"
	"domain/subscriptions/subscribers"
	"domain/subscriptions/subscriptions"
)

// EventPublisher specifies interface for publishing domain events.
type EventPublisher interface {
	Publish(topic string, payload interface{})
}

// Usecase specifies list use cases.
type Usecase struct {
	db     orm.DB
	events EventPublisher
}

// Create a subscriber list.
func (u *Usecase) Create(workspacePK uuid.UUID, title string) (*lists.List, error) {
	list := &lists.List{
		PK:          uuid.Must(uuid.NewV4()),
		WorkspacePK: workspacePK,
		Title:       title,
		Version:     1,
	}

	if err := lists.Create(u.db, list); err != nil {
		return nil, err
	}

	u.events.Publish("subscriptions", "list created")

	return list, nil
}

// Rename a subscriber list.
func (u *Usecase) Rename(list *lists.List, title string) error {
	if err := lists.Rename(u.db, list, title); err != nil {
		return err
	}

	u.events.Publish("subscriptions", "list renamed")

	return nil
}

// Delete a subscriber list.
func (u *Usecase) Delete(list *lists.List) error {
	if err := lists.Delete(u.db, list); err != nil {
		return err
	}

	u.events.Publish("subscriptions", "list deleted")

	return nil
}

// Subscribe an email address into a list.
func (u *Usecase) Subscribe(subscriber *subscribers.Subscriber, list *lists.List, data map[string]interface{}) (*subscriptions.Subscription, error) {
	subscription := &subscriptions.Subscription{
		PK:           uuid.Must(uuid.NewV4()),
		SubscriberPK: subscriber.PK,
		ListPK:       list.PK,
		EmailAddress: subscriber.EmailAddress,
		Data:         data,
		Version:      1,
	}

	if err := subscriptions.Create(u.db, subscription); err != nil {
		return nil, err
	}

	u.events.Publish("subscriptions", "subscription created")

	return subscription, nil
}

// Cancel a subscription.
func (u *Usecase) Cancel(subscription *subscriptions.Subscription) error {
	if err := subscriptions.Cancel(u.db, subscription); err != nil {
		return err
	}

	u.events.Publish("subscriptions", "subscription canceled")

	return nil
}

// Update a subscription with new data.
func (u *Usecase) Update(subscription *subscriptions.Subscription, data map[string]interface{}) error {
	if err := subscriptions.Update(u.db, subscription, data); err != nil {
		return err
	}

	u.events.Publish("subscriptions", "subscription updated")

	return nil
}

// Forget a subscriber.
func (u *Usecase) Forget(subscriber *subscribers.Subscriber) error {
	if err := subscribers.Forget(u.db, subscriber); err != nil {
		return err
	}

	u.events.Publish("subscriptions", "subscriber forgotten")

	return nil
}
