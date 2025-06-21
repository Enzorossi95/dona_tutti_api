package donor

import (
	"github.com/google/uuid"
)

type Donor struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	IsVerified bool      `json:"is_verified"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
}
