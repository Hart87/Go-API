package models

import (
	//"encoding/json"
	"time"
)

type User struct {
	Name       string        `json:"name"`
	Email      string        `json:"email"`
	Password   string        `json:"password"`
	ID         []byte        `json:"id"`
	Membership string        `json:"membership"`
	CreatedAt  time.Duration `json:"createdAt"`
}
