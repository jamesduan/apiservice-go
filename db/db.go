package db

import (
	"database/sql"
	"log"

	"apiservice/g"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", g.Config().Database)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	DB.SetMaxOpenConns(g.Config().MaxConns)
	DB.SetMaxIdleConns(g.Config().MaxIdle)

	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}

// insert and return id
func NewRecord(sql string) (int64, error) {
	stmt, err := DB.Prepare(sql)
	if err != nil {
		log.Printf("DB.Prepare failed, sql:%s, error:%v\n", sql, err)
		return 0, err
	}

	res, err := stmt.Exec()
	if err != nil {
		log.Printf("stmt.Exec failed, sql:%s, error:%v\n", sql, err)
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("res.LastInsertId failed, sql:%s, error:%v\n", sql, err)
		return 0, err
	}

	return lastId, nil
}

// update record
func UpdateRecord(sql string) error {
	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}

	return err
}
