/*  dependecyInjector.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 08, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 08/09/18 23:44
 */

package injector

import (
	"fmt"
	"time"

	config "github.com/spf13/viper"
	"github.com/vaksi/foreign_currency/infrastructures"
	"github.com/vaksi/foreign_currency/internal/exchangeRate"
	"github.com/vaksi/foreign_currency/internal/exchangeRateTrack"
)

// NewDependencyInjection as Aggregates DI
func NewDependencyInjection() *DependencyInjection {
	return &DependencyInjection{}
}

// DependencyInjection defines attribute for Dependency injector
type DependencyInjection struct {
	Database                 infrastructures.MySQLFactory
	ExchangeRateService      exchangeRate.Service
	ExchangeRateTrackService exchangeRateTrack.Service
}

// Inject injects dependency to route
func (injector *DependencyInjection) Inject() {
	// Repository
	exchangeRateRepository := exchangeRate.NewRepositoryFactory(injector.Database)
	exchangeRateTrackRepository := exchangeRateTrack.NewRepositoryFactory(injector.Database)
	// Service
	injector.ExchangeRateService = exchangeRate.Service{ExchangeRateRepo: exchangeRateRepository}
	injector.ExchangeRateTrackService = exchangeRateTrack.Service{ExchangeRateTrackRepo: exchangeRateTrackRepository}
}

// InitSQL set the sql values
func InitSQL() *infrastructures.MySQL {
	masterUser := config.GetString("database_sql_master.user")
	masterPassword := config.GetString("database_sql_master.password")
	masterHost := config.GetString("database_sql_master.host")
	masterPort := config.GetInt("database_sql_master.port")
	masterDBName := config.GetString("database_sql_master.db_name")
	masterCharset := config.GetString("database_sql_master.charset")
	dataSourceMaster := fmt.Sprintf(infrastructures.DataSourceFormat, masterUser, masterPassword, masterHost,
		masterPort, masterDBName, masterCharset)

	slaveUser := config.GetString("database_sql_slave.user")
	slavePassword := config.GetString("database_sql_slave.password")
	slaveHost := config.GetString("database_sql_slave.host")
	slavePort := config.GetInt("database_sql_slave.port")
	slaveDBName := config.GetString("database_sql_slave.db_name")
	slaveCharset := config.GetString("database_sql_slave.charset")
	dataSourceSlave := fmt.Sprintf(infrastructures.DataSourceFormat, slaveUser, slavePassword, slaveHost,
		slavePort, slaveDBName, slaveCharset)

	sqlInfra := new(infrastructures.MySQL)
	sqlInfra.OpenConnection(dataSourceMaster, dataSourceSlave)
	sqlInfra.SetConnMaxLifetime(config.GetDuration("database.max_life_time") * time.Second)
	sqlInfra.SetMaxIdleConns(config.GetInt("database.max_idle_connection"))
	sqlInfra.SetMaxOpenConns(config.GetInt("database.max_open_connection"))

	return sqlInfra
}
