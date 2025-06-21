package organizer

import (
	"time"

	"github.com/google/uuid"
)

// Organizer represents the domain entity for organizers
type Organizer struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Verified  bool      `json:"verified"`
}
