package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User is used by pop to map your users database table to your go code.
type User struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"-" db:"password"`
	LoginAt    time.Time `json:"login_at" db:"login_at"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	Provider   string    `json:"provider" db:"provider"`
	ProviderID string    `json:"provider_id" db:"provider_id"`
	GroupID    nulls.Int `json:"-" db:"group_id"`
	Group      *Group    `json:"group,omitempty" belongs_to:"group"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Create is to create user
func (u *User) Create(tx *pop.Connection) error {
	u.Email = strings.ToLower(u.Email)
	pw, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(pw)
	return tx.Create(u)
}
