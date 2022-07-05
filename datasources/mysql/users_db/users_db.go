package users_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysql_username = "mysql_username"
	mysql_password = "mysql_password"
	mysql_host     = "mysql_host"
	mysql_schema   = "mysql_schema"
)

var (
	ClientDB *sql.DB
	username = os.Getenv(mysql_username)
	password = os.Getenv(mysql_password)
	host     = os.Getenv(mysql_host)
	schema   = os.Getenv(mysql_schema)
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		username, password, host, schema,
	)
	var err error
	ClientDB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	if err = ClientDB.Ping(); err != nil {
		panic(err)
	}
	log.Println("Database Successfully Configured")
}
