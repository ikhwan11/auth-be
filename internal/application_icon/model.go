package application_icon

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationIcon struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	FileData  []byte    `db:"file_data" json:"-"`
	MimeType  string    `db:"mime_type" json:"mime_type"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
