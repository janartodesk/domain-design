package domain

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofrs/uuid"
)

// List is a subscriber list.
type List struct {
	PK      uuid.UUID
	Title   string
	Version uint32
}

// Validate the subscriber list.
func (l *List) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.PK, validation.Required),
		validation.Field(&l.Title, validation.Required),
		validation.Field(&l.Version, validation.Required),
	)
}

// CreateList creates a subscriber list.
func CreateList(title string) (*List, error) {
	list := &List{
		PK:      uuid.Must(uuid.NewV4()),
		Title:   title,
		Version: 1,
	}

	if err := list.Validate(); err != nil {
		return list, err
	}

	return list, nil
}

// RenameList renames a subscriber list.
func RenameList(list List, title string) (*List, error) {
	list.Title = strings.TrimSpace(title)
	list.Version++

	if err := list.Validate(); err != nil {
		return nil, err
	}

	return &list, nil
}
