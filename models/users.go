package models

type User struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	ID         string `json:"id"`
	Membership string `json:"membership"`
	CreatedAt  int    `json:"createdAt"`
}
