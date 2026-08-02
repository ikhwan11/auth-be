package application

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Code      string     `db:"code" json:"code"`
	URL       string     `db:"url" json:"url"`
	IconID    *uuid.UUID `db:"icon_id" json:"icon_id"`
	IsDefault bool       `db:"is_default" json:"is_default"`
	IsActive  bool       `db:"is_active" json:"is_active"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}
