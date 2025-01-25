package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	Name      string         `db:"name" json:"name"`
	Email     string         `db:"email" json:"email"`
	Password  string         `db:"password" json:"-"`
	Role      int16          `db:"role" json:"role"`
	ImageURI  sql.NullString `db:"image_uri" json:"image_uri"`
	CreatedAt string         `db:"created_at" json:"created_at"`
	UpdatedAt string         `db:"updated_at" json:"updated_at"`
}
