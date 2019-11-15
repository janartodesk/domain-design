package lists

import (
	"github.com/go-pg/pg/orm"
	"github.com/gofrs/uuid"
)

// Get a subscriber list.
func Get(db orm.DB, workspacePK, listPK uuid.UUID) (*List, error) {
	list := &List{}

	if err := db.Model(list).Where("pk = ? AND workspace_pk = ?", listPK, workspacePK).Select(); err != nil {
		return nil, err
	}

	return list, nil
}
