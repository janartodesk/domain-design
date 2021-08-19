package domain

import "github.com/gofrs/uuid"

// Subscriber is an entity subscribed to a list.
type Subscriber struct {
	PK           uuid.UUID
	EmailAddress EmailAddress
	Version      uint32
}

// Validate the subscriber.
func (s *Subscriber) Validate() error {
	return nil
}

// CreateSubscriber creates a subscriber.
func CreateSubscriber(addr EmailAddress) (*Subscriber, error) {
	s := &Subscriber{
		PK:           uuid.Must(uuid.NewV4()),
		EmailAddress: addr,
	}

	if err := s.Validate(); err != nil {
		return nil, err
	}

	return s, nil
}

// ForgetSubscriber forgets a subscriber.
func ForgetSubscriber(s Subscriber) (*Subscriber, error) {
	s.EmailAddress = "forgotten@smaily.email"

	return &s, nil
}
