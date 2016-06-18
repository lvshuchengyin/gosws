// dbpool
package db

import (
	"database/sql"
	"errors"

	"github.com/lvshuchengyin/gosws/logger"
)

var (
	dbPool *DBPool
)

type DBPool struct {
	source  string
	maxIdle int
	maxOpen int

	db *sql.DB
}

func InitDBPool(source string, maxIdle, maxOpen int) *DBPool {
	pool := &DBPool{
		source:  source,
		maxIdle: maxIdle,
		maxOpen: maxOpen,
	}

	pool.connect()
	dbPool = pool

	return pool
}

func (self *DBPool) connect() (err error) {
	//sourceName := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s", self.username, self.password, self.protocol, hp[0], hp[1], self.dbname, self.charset)
	var db *sql.DB
	db, err = sql.Open("mysql", self.source)
	if err != nil {
		logger.Error("dbconn connect fail: %s", err.Error())
		return
	}

	db.SetMaxIdleConns(self.maxIdle)
	db.SetMaxOpenConns(self.maxOpen)

	self.db = db
	return
}

func (self *DBPool) getDB() (db *sql.DB, err error) {
	if self.db == nil {
		return nil, errors.New("dbpool not init")
	}

	return self.db, nil
}

// add, del, update
func DBExec(query string, args ...interface{}) (result sql.Result, err error) {
	db, err := dbPool.getDB()
	if err != nil {
		logger.Error("get b err:%v", err)
		return
	}

	result, err = db.Exec(query, args...)
	return
}

// remember rows must call Close()
func DBQuery(query string, args ...interface{}) (rows *sql.Rows, err error) {
	db, err := dbPool.getDB()
	if err != nil {
		return
	}

	rows, err = db.Query(query, args...)
	return
}

func DBQueryRow(query string, args ...interface{}) (row *sql.Row) {
	db, err := dbPool.getDB()
	if err != nil {
		return
	}

	row = db.QueryRow(query, args...)
	return
}
