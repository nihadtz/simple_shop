package services

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
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

func MigrateDB() {
	a := new(AccessCtrl)

	fmt.Println(("mig"))

	db, err := sql.Open("mysql", a.getDNS()+"?multiStatements=true")

	if err != nil {
		LogError("Migration Error connecting to database", err)
		return
	}

	fmt.Println(("mig"))

	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		LogError("Migration Error getting database driver", err)
		return
	}

	fmt.Println(("mig"))

	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)

	if err != nil {
		fmt.Println(err.Error())
		LogError("Error creating migration instance", err)
		return
	}

	fmt.Println(("mig"))

	migration.Steps(Configuration.DB.MigrationStep)
}
