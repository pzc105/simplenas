package db

import (
	"database/sql"
	"time"

	"pnas/log"
	"pnas/setting"

	_ "github.com/go-sql-driver/mysql"
)

var mysql_db *sql.DB

func Init() {
	db, err := sql.Open("mysql", setting.GetMysqlConnectStr())
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(setting.GS.Mysql.MaxOpenConns)
	db.SetMaxIdleConns(setting.GS.Mysql.MaxIdleConns)

	err = db.Ping()
	if err != nil {
		log.Errorf("[db] %v", err)
	}

	mysql_db = db

	initRedis()
}

func Exec(sql string, args ...any) (sql.Result, error) {
	return mysql_db.Exec(sql, args...)
}

func QueryRow(sql string, args ...any) *sql.Row {
	return mysql_db.QueryRow(sql, args...)
}

func Query(sql string, args ...any) (*sql.Rows, error) {
	return mysql_db.Query(sql, args...)
}
