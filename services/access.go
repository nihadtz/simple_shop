package services

import (
	"database/sql"
	"fmt"
	"os"
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
	sqlPort     = os.Getenv("MYSQL_PORT")

	sqlUser     = os.Getenv("MYSQL_USER")
	sqlPassword = os.Getenv("MYSQL_PASSWORD")
	sqlAddress  = os.Getenv("MYSQL_ADDRESS")
	dbName      = os.Getenv("MYSQL_DATABASE")
)

type AccessCtrl struct {
	ShopSQLDB *sqlx.DB
}

func NewAccess(runas string) {
	a := new(AccessCtrl)

	var err error
	var connString string

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

func MigrateDB(runas string) {
	a := new(AccessCtrl)
	var db *sql.DB
	var err error

	if runas == "prod" {
		sqlProtocol = "unix"

		db, err = sql.Open("mysql", a.getDNSProd()+"&multiStatements=true")

		if err != nil {
			fmt.Println(err.Error())
			LogError("Migration Error connecting to database", err)
			return
		}
	} else {
		db, err = sql.Open("mysql", a.getDNS()+"?multiStatements=true")

		if err != nil {
			fmt.Println(err.Error())
			LogError("Migration Error connecting to database", err)
			return
		}
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		fmt.Println(err.Error())
		LogError("Migration Error getting database driver", err)
		return
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(err.Error())
		LogError("Error creating migration instance", err)
		return
	}

	migration.Steps(Configuration.DB.MigrationStep)
}
