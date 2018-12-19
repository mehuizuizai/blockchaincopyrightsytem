// DbDao project DbDao.go
package DbDao

import (
	"database/sql"
	"fmt"
	"ledger/DbUtil"
	"log"

	_ "ledger/github.com/go-sql-driver/mysql"
)

func QueryForMap(sql string, args ...interface{}) (map[string]string, error) { // how to []map -> map
	//	results := []map[string]string{}
	results, err := QueryForMapSlice(sql, args...)
	if err != nil {
		log.Fatal("Get result failed:", err.Error())

	}
	if len(results) > 0 {
		result := results[0]
		fmt.Println(result)
		return result, err
	}
	return nil, err
}
func QueryForMapSlice(sql string, args ...interface{}) ([]map[string]string, error) {
	db, err := DbUtil.DbOpen()
	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
	}
	defer db.Close()
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatal("prepare failed:", err.Error())
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()
	result, err := HandleResultToMapSlice(rows)
	return result, err

}

// InserToDb
// tx ->transaction
func InsertToDb_PreExe(sql string, args ...interface{}) map[string]string {
	db, err := DbUtil.DbOpen()
	tx, err := db.Begin()
	fmt.Println(sql, args)
	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
		return nil
	}
	defer db.Close()
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Fatal("prepare failed:", err.Error())
		return nil
	}
	defer stmt.Close()
	result, err := stmt.Exec(args...)
	fmt.Println("here")
	if err != nil {
		//		log.Fatal("Insert failed:", err.Error())
		return nil
	}
	affect, err := result.RowsAffected()
	if err != nil {
		log.Fatal("GetRowsAffected failed:", err.Error())
		return nil
	}
	workList, _ := QueryForMap("SELECT * FROM tb_Work where WorkID =?", args[0]) //WorkID
	tx.Rollback()
	fmt.Println("Success InserWork !affect is :", affect)
	return workList
}

func InsertToDb(sql string, args ...interface{}) error {
	db, err := DbUtil.DbOpen()
	fmt.Println(sql, args)
	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatal("prepare failed:", err.Error())
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(args...)
	if err != nil {
		log.Fatal("Insert failed:", err.Error())
		return err
	}
	affect, err := result.RowsAffected()
	if err != nil {
		log.Fatal("GetRowsAffected failed:", err.Error())
		return err
	}
	fmt.Println("Success InserWork !affect is :", affect)
	return nil
}

//UpdateDb
func UpdateDb(sql string, args ...interface{}) error {
	db, err := DbUtil.DbOpen()

	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatal("prepare failed:", err.Error())
		return err
	}
	defer stmt.Close()
	result, err := stmt.Query(args...)
	if err != nil {
		log.Fatal("Update failed:", err.Error())
		return err
	}

	fmt.Println("Success Update !", result)
	return nil

}
func UpdateDb_PreExe(sql string, args ...interface{}) (error, map[string]string) {
	db, err := DbUtil.DbOpen()
	tx, err := db.Begin()
	fmt.Println(sql, args)
	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
		return err, nil
	}
	defer db.Close()
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Fatal("prepare failed:", err.Error())
		return err, nil
	}
	defer stmt.Close()
	result, err := stmt.Query(args...)
	print(result)
	if err != nil {
		log.Fatal("Update failed:", err.Error())
		return err, nil
	}
	workList, _ := QueryForMap("SELECT * FROM tb_Work where WorkID =?", args[2]) //WorkID
	tx.Rollback()
	return nil, workList

}
func HandleResultToMapSlice(rows *sql.Rows) ([]map[string]string, error) {
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal("Get rows failed:", err.Error())
	}
	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols)) //_,val
	for i := range values {
		scans[i] = &values[i] //scans -> values
	}
	results := []map[string]string{}
	i := 0
	for rows.Next() { // scan all rows
		if err := rows.Scan(scans...); err != nil {
			log.Fatal("Scan result err:", err.Error())
		}
		row := make(map[string]string)
		for k, v := range values {
			key := cols[k]
			row[key] = string(v)
		}
		results = append(results, row)
		i++
	}
	return results, nil
}
