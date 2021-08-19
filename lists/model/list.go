package model

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/janartodesk/domain-design/lists/domain"
	"github.com/uptrace/bun"
)

// List is a database model for a subscriber list.
type List struct {
	PK      uuid.UUID `bun:"pk,pk"`
	Title   string    `bun:"title"`
	Version uint32    `bun:"version"`

	bun.BaseModel `bun:"lists"`
}

// CreateList creates a subscriber list.
func CreateList(db bun.IDB, list *domain.List) error {
	if _, err := db.NewInsert().Model(&List{
		PK:      list.PK,
		Title:   list.Title,
		Version: list.Version,
	}).Exec(context.Background()); err != nil {
		return err
	}

	return nil
}

// DeleteList deletes a subscriber list.
func DeleteList(db bun.IDB, pk uuid.UUID) error {
	res, err := db.NewDelete().Model(&List{
		PK: pk,
	}).WherePK().Exec(context.Background())

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

// UpdateList updates a subscriber list.
func UpdateList(db bun.IDB, list *domain.List) error {
	res, err := db.NewUpdate().Model(&List{
		PK:      list.PK,
		Title:   list.Title,
		Version: list.Version,
	}).Where(
		"pk = ? AND version = ?",
		list.PK,
		list.Version-1,
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

// GetList returns a subscriber list.
func GetList(db bun.IDB, pk uuid.UUID) (*domain.List, error) {
	model := List{
		PK: pk,
	}

	if err := db.NewSelect().Model(&model).WherePK().Scan(context.Background()); err != nil {
		return nil, err
	}

	return &domain.List{
		PK:      model.PK,
		Title:   model.Title,
		Version: model.Version,
	}, nil
}

// ListLists returns a list of subscriber lists.
func ListLists(db bun.IDB, offset, limit uint32) ([]*domain.List, error) {
	model := []List{}

	if err := db.NewSelect().Model(&model).Scan(context.Background()); err != nil {
		return nil, err
	}

	res := []*domain.List{}

	for _, list := range model {
		res = append(res, &domain.List{
			PK:      list.PK,
			Title:   list.Title,
			Version: list.Version,
		})
	}

	return res, nil
}
