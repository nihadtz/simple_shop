package services

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	Access      *AccessCtrl
	sqlDriver   = "mysql"
	sqlProtocol = "tcp"
	sqlPort     = "3306"

	sqlUser     = "api"
	sqlPassword = "test"
	sqlAddress  = "mysql-db"
	dbName      = "shop"
)

type AccessCtrl struct {
	ShopSQLDB *sqlx.DB
}

func NewAccess() {
	a := new(AccessCtrl)

	var err error

	a.ShopSQLDB, err = sqlx.Connect(sqlDriver, a.getDNS())

	if err != nil {
		LogError("Error connecting to database", err)
	} else {
		fmt.Println("Connected to database")
	}

	Access = a
}

func (a AccessCtrl) getDNS() string {
	return sqlUser + ":" + sqlPassword + "@" + sqlProtocol + "(" + sqlAddress + ":" + sqlPort + ")" + "/" + dbName
}
