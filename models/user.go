package models

import (
	"errors"

	"github.com/nihadtz/simple_shop/services"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	Email        string `db:"email" json:"email"`
	Password     string `db:"-" json:"password"`
	PasswordHash []byte `db:"password" json:"-"`
	Type         string `db:"type" json:"type"`
	Token        string `db:"token" json:"token"`
	Algorithm    string `db:"alg" json:"-"`
	Issued       int64  `db:"issued" json:"issued"`
}

var createUser = `CREATE TABLE User (
    id int(11) NOT NULL AUTO_INCREMENT,
    name varchar(128) NOT NULL DEFAULT '',
    email varchar(128) NOT NULL DEFAULT '',
    password varchar(128) NOT NULL DEFAULT '',
    type ENUM('customer','admin') NOT NULL DEFAULT 'customer',
    token varchar(255) NOT NULL DEFAULT '',
    alg varchar(10) NOT NULL DEFAULT '',
    issued int(11) NOT NULL DEFAULT 0,
    PRIMARY KEY(id),
    UNIQUE KEY(email)
    ) ENGINE=INNODB`

func (u *User) Get(id int) error {
	db := services.Access.GetDB()

	query := `SELECT 
				id, name, email, password, type, token, alg, issued
			FROM
			 	User
			WHERE
				id = ?`

	err := db.Get(u, query, id)

	return err
}

func (u *User) GetByEmail(email string) error {
	db := services.Access.GetDB()

	query := `SELECT 
				id, name, email, password, type, token, alg, issued
			FROM
			 	User
			WHERE
				email = ?`

	err := db.Get(u, query, email)

	return err
}

func (u *User) GetByEmailPassword(email, password string) error {
	err := u.GetByEmail(email)

	if err != nil {
		err = errors.New("This User doesn't exists")
		services.LogError("Error no account with specified email", err)
		return err
	}

	err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))

	if err != nil {
		services.LogError("Error Matcing password", err)
		return err
	}

	return nil
}

func (u *User) Create() (err error) {
	db := services.Access.GetDB()
	u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), 12)

	if err != nil {
		services.LogError("Error bcrypt password", err)
		return err
	}

	query := `INSERT INTO User
						(name, email, password, type) 
					VALUES
						(:name, :email, :password, :type)`

	_, err = db.NamedExec(query, u)

	return err
}

func (u User) UpdateTokenInfo() error {
	db := services.Access.GetDB()

	query := `UPDATE User SET
				token=?, alg=?, issued=?
			WHERE 
				id=?`

	_, err := db.Exec(query, u.Token, u.Algorithm, u.Issued, u.ID)

	if err != nil {
		services.LogError("Error updating Token", err)
	}

	return err
}

func (u *User) ClearToken() error {
	u.Token = ""
	u.Algorithm = ""
	u.Issued = 0

	err := u.UpdateTokenInfo()

	if err != nil {
		services.LogError("Error clearing Token", err)
	}

	return err
}
