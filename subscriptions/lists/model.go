package lists

import (
	"github.com/gofrs/uuid"
)

// List represents a subscriber list.
type List struct {
	PK          uuid.UUID `sql:"pk,pk"`
	WorkspacePK uuid.UUID `sql:"workspace_pk"`
	Title       string    `sql:"title"`
	Version     uint32    `sql:"version"`

	tableName struct{} `sql:"lists"`
}
