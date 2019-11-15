package lists

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-pg/pg/orm"
)

// Create a subscriber list.
func Create(db orm.DB, list *List) error {
	list.Version = 1

	if err := normalize(list); err != nil {
		return err
	} else if err := validate(list); err != nil {
		return err
	}

	if _, err := db.Model(list).Insert(); err != nil {
		return err
	}

	return nil
}

// Rename a subscriber list.
func Rename(db orm.DB, list *List, title string) error {
	// Stash current title for change comparison.
	prevTitle := list.Title

	list.Title = title
	list.Version++

	if err := normalize(list); err != nil {
		return err
	} else if err := validate(list); err != nil {
		return err
	}

	if list.Title == prevTitle {
		return errors.New("title did not change")
	}

	if r, err := db.Model(list).Where("workspace_pk = ?workspace_pk AND list_pk = ?list_pk AND version = ?version").Update(); err != nil {
		return err
	} else if r.RowsAffected() == 0 {
		return errors.New("list was not found")
	}

	return nil
}

// Delete a subscriber list.
func Delete(db orm.DB, list *List) error {
	if r, err := db.Model(list).Where("workspace_pk = ?workspace_pk AND list_pk = ?list_pk AND version = ?version").Delete(); err != nil {
		return err
	} else if r.RowsAffected() == 0 {
		return errors.New("list was not found")
	}

	return nil
}

func normalize(list *List) error {
	list.Title = strings.TrimSpace(list.Title)

	return nil
}

func validate(list *List) error {
	return validation.ValidateStruct(
		list,
		validation.Field(&list.PK, validation.Required),
		validation.Field(&list.WorkspacePK, validation.Required),
		validation.Field(&list.Title, validation.Required),
		validation.Field(&list.Version, validation.Required),
	)
}
