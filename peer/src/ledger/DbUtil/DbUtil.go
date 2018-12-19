// DbUtil project DbUtil.go
package DbUtil

import (
	"database/sql"

	_ "ledger/github.com/go-sql-driver/mysql"
)

func DbInit() (s *sql.DB, err error) {
	db, err := sql.Open("mysql", "root:wahaha@tcp(127.0.0.1:3306)/?charset=utf8")
	return db, err
}
func DbOpen() (s *sql.DB, err error) {
	db, err := sql.Open("mysql", "root:wahaha@tcp(127.0.0.1:3306)/blockchain2?charset=utf8")
	return db, err
}
