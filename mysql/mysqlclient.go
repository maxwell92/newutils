package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	mylog "github.com/maxwell92/newutils/log"
	"sync"
	"time"
)

var log = mylog.Log

type MysqlClient struct {
	DB *sql.DB
}

var instance *MysqlClient

var once sync.Once

func MysqlInstance() *MysqlClient {
	once.Do(func() {
		instance = new(MysqlClient)
	})
	return instance
}

type MySQLConfig struct {
	Driver    string
	Endpoints string
	MaxOpen   int
	MaxIdle   int
}

func (c *MysqlClient) Open(c *MySQLConfig) {

	db, err := sql.Open(c.Driver, c.Endpoints)

	if err != nil {
		log.Fatalf("MysqlClient Open Error: err=%s", err)
		return
	}

	// Set Connection Pool
	db.SetMaxOpenConns(c.MaxOpen)
	db.SetMaxIdleConns(c.MaxIdle)

	c.DB = db

}

func (c *MysqlClient) Close() {
	c.DB.Close()
}

func (c *MysqlClient) Conn() *sql.DB {
	return c.DB
}

// Ping the connection, keep connection alive
func (c *MysqlClient) Ping() {
	select {
	case <-time.After(time.Millisecond * time.Duration(config.DELAY_MILLISECONDS)):
		err := c.DB.Ping()
		if err != nil {
			log.Fatalf("MysqlClient Ping Error: err=%s", err)
			c.Open()
		}
	}
}
