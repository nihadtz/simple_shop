package services

import (
	"fmt"
	"strconv"
	"time"

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

func NewAccess(runas string) {
	a := new(AccessCtrl)

	var err error
	var connString string

	sqlPort = strconv.Itoa(Configuration.DB.Port)
	sqlUser = Configuration.DB.User
	sqlPassword = Configuration.DB.Password
	sqlAddress = Configuration.DB.Address
	dbName = Configuration.DB.DBName

	if runas == "prod" {
		sqlProtocol = "unix"

		connString = a.getDNSProd()
	} else {
		connString = a.getDNS()
	}

	a.ShopSQLDB, err = sqlx.Connect(sqlDriver, connString)

	if err != nil {
		LogError("Error connecting to database", err)
	} else {
		fmt.Println("Connected to database")

		a.ShopSQLDB.SetMaxOpenConns(Configuration.DB.MaxOpenConns)
		a.ShopSQLDB.SetMaxIdleConns(Configuration.DB.MaxIdleConns)
		a.ShopSQLDB.SetConnMaxLifetime(Configuration.DB.ConnMaxLifetime * time.Minute)
	}

	Access = a
}

func (a AccessCtrl) getDNS() string {
	return sqlUser + ":" + sqlPassword + "@" + sqlProtocol + "(" + sqlAddress + ":" + sqlPort + ")" + "/" + dbName
}

func (a AccessCtrl) getDNSProd() string {
	return sqlUser + ":" + sqlPassword + "@" + sqlProtocol + "(" + sqlAddress + ")" + "/" + dbName + "?parseTime=true"
}

func (a AccessCtrl) GetDB() *sqlx.DB {
	return a.ShopSQLDB
}
