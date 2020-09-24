package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"
)

// Group is used by pop to map your groups database table to your go code.
type Group struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name"`
	Users     Users     `json:"users,omitempty" has_many:"users"`
}

// String is not required by pop and may be deleted
func (g Group) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// Groups is not required by pop and may be deleted
type Groups []Group

// String is not required by pop and may be deleted
func (g Groups) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// Create is to create group
func (g *Group) Create(tx *pop.Connection) error {
	g.Name = strings.ToLower(g.Name)
	return tx.Create(g)
}

// Update is to update group
func (g *Group) Update(tx *pop.Connection) error {
	g.Name = strings.ToLower(g.Name)
	return tx.Update(g)
}

// Delete is to update group
func (g *Group) Delete(tx *pop.Connection) error {
	return tx.Destroy(g)
}
