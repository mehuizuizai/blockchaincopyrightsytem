// DbUtil project DbUtil.go
package DbUtil

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DbOpen() (s *sql.DB, err error) {
	db, err := sql.Open("mysql", "root:wahaha@tcp(127.0.0.1:3306)/block_chain?charset=utf8")
	return db, err
}
