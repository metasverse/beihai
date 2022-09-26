package g

import (
	"database/sql"
	"github.com/eatmoreapple/sqlbuilder"
	"github.com/go-redis/redis/v9"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DB    *sql.DB
	Redis *redis.Client
)

// Session is an interface for sql.DB and sql.Tx
type Session interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	sqlbuilder.Session
}
