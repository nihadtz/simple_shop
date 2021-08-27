package models

type User struct {
	ID           int    `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	Email        string `db:"email" json:"email"`
	Password     string `db:"-" json:"password"`
	PasswordHash string `db:"password" json:"-"`
	Type         string `db:"type" json:"type"`
	Token        string `db:"token" json:"token"`
	Algorithm    string `db:"alg" json:"-"`
	Issued       string `db:"issued" json:"issued"`
}
