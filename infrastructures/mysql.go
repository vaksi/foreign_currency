/*  mysql.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 09, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 09/09/18 16:25
 */

package infrastructures

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ahartanto/morphling"
	_ "github.com/go-sql-driver/mysql"
)

// MySQLFactory
type MySQLFactory interface {
	OpenConnection(dataSourceMaster, dataSourceSlave string)
	GetDB() (*morphling.DB, error)
	SetConnMaxLifetime(time.Duration)
	SetMaxIdleConns(int)
	SetMaxOpenConns(int)
}

//MySQL DBslave or master
type MySQL struct {
	DB *morphling.DB
}

var DataSourceFormat = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local"

func NewSQLInfrastructure(db *morphling.DB) MySQLFactory {
	return &MySQL{
		DB: db,
	}
}

// OpenConnection gets a handle for a database
func (s *MySQL) OpenConnection(dataSourceMaster, dataSourceSlave string) {

	db, err := morphling.Open(morphling.MySQLDriver, dataSourceMaster, dataSourceSlave)
	if err != nil {
		log.Panic(fmt.Errorf("database : %v", err))
	}

	err = db.Ping()
	if err != nil {
		log.Panic(fmt.Errorf("failed dial database : %v", err))
	}

	s.DB = db
}

// GetDB gets database connection
func (s *MySQL) GetDB() (*morphling.DB, error) {
	err := s.DB.Ping()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return s.DB, nil
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (s *MySQL) SetConnMaxLifetime(connMaxLifetime time.Duration) {
	s.DB.SetConnMaxLifetime(connMaxLifetime)
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
func (s *MySQL) SetMaxIdleConns(maxIdleConn int) {
	s.DB.SetMaxIdleConns(maxIdleConn)
}

// SetMaxOpenConns sets the maximum amount of time a connection may be reused.
func (s *MySQL) SetMaxOpenConns(maxOpenConn int) {
	s.DB.SetMaxOpenConns(maxOpenConn)
}
