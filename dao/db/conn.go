package db

import (
	"database/sql"
	"filestore-server/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"sync"
)

type MySQLConn struct {
	conn *sql.DB
}

var (
	once   sync.Once
	dbConn *MySQLConn
)

func GetDBInstance() *MySQLConn {
	once.Do(func() {
		dbConn = &MySQLConn{}
	})
	return dbConn
}

func (m *MySQLConn) GetConn() *sql.DB {
	return dbConn.conn
}

func (m *MySQLConn) Init(conf config.MySQLConf) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.DBPwd,
		conf.Host,
		conf.Port,
		conf.DBName)
	var err error
	m.conn, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Error(err.Error())
	}
	m.conn.SetMaxIdleConns(conf.MaxIdleConns)
	m.conn.SetMaxOpenConns(conf.MaxOpenConns)

	err = m.conn.Ping()
	if err != nil {
		log.Error(err.Error())
	}
}
