package organizer

import (
	"time"

	"github.com/google/uuid"
)

// Organizer represents the domain entity for organizers
type Organizer struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Verified  bool      `json:"verified"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Website   string    `json:"website"`
	Address   string    `json:"address"`
}
