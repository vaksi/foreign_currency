/*  repository.go.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 05, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 05/09/18 03:06
 */

package exchangeRate

import (
	"time"

	"github.com/vaksi/foreign_currency/infrastructures"
)

// Repository of struct exchange rate
type Repository struct {
	DB infrastructures.MySQLFactory
}

// RepositoryFactory of interface repository factory
type RepositoryFactory interface {
	Store(*ExchangeRate) error
	WhereByDate(date time.Time) ([]*ExchangeRate, error)
	WhereByTrend(from, to string) ([]*ExchangeRate, error)
}

// NewRepositoryFactory of new factory repository
func NewRepositoryFactory(db infrastructures.MySQLFactory) RepositoryFactory {
	return &Repository{
		DB: db,
	}
}

// Store is function store exchange rate repository
func (r *Repository) Store(exchangeRate *ExchangeRate) (err error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}

	// prepare transaction
	stmt, err := tx.Prepare("INSERT INTO exchange_rates (oid, date, `from`, `to`, rate) VALUES (?,?,?,?,?)")
	if err == nil {
		return err
	}
	_, err = stmt.Exec(exchangeRate.ID(), exchangeRate.date.Format("2006-01-02"), exchangeRate.from, exchangeRate.to, exchangeRate.rate)
	if err != nil {
		// roll back db if error
		if errTx := tx.Rollback(); errTx != nil {
			return errTx
		}
		return err
	}

	// commit transaction
	if errTx := tx.Commit(); errTx != nil {
		return errTx
	}

	return
}

// WhereByDate is function find by id of exchange rate repository
func (r *Repository) WhereByDate(date time.Time) (exchangeRates []*ExchangeRate, err error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return
	}

	// get last 7 day
	lastSevenDay := date.AddDate(0, 0, -6)

	// query define
	rows, err := db.Query(
		`SELECT ex.date, ex.from, ex.to,
			IFNULL(ex.rate, 0) rate, 
			IFNULL((SELECT avg(rate) FROM exchange_rates e WHERE e.from=ex.from and e.to = ex.to AND e.date BETWEEN ? and ? ), 0) as average
		FROM exchange_rates ex, tracks t
		WHERE t.from = ex.from AND t.to = ex.to AND ex.date = ?
		ORDER BY ex.date DESC`,
		lastSevenDay.Format("2006-01-02"), date.Format("2006-01-02"), date.Format("2006-01-02"),
	)
	if err != nil {
		return
	}

	for rows.Next() {
		var exc ExchangeRate
		err = rows.Scan(
			&exc.date,
			&exc.from,
			&exc.to,
			&exc.rate,
			&exc.sevenDayRates,
		)
		// skip when scan error
		if err != nil {
			continue
		}

		exchangeRates = append(exchangeRates, &exc)
	}
	return
}

// WhereByTrend is function to retrieve all exchange rate
func (r *Repository) WhereByTrend(from, to string) (exchangeRates []*ExchangeRate, err error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return
	}

	rows, err := db.Query(
		`SELECT ex.date, ex.from, ex.to, IFNULL(ex.rate, 0) rate
		FROM exchange_rates ex, tracks t
		WHERE t.from = ex.from AND t.to = ex.to 
			AND t.from=? AND t.to=?
		ORDER BY ex.date DESC LIMIT 7`,
		from, to,
	)

	if err != nil {
		return
	}

	for rows.Next() {
		var exc ExchangeRate
		err = rows.Scan(
			&exc.date,
			&exc.from,
			&exc.to,
			&exc.rate,
		)
		// skip when scan error
		if err != nil {
			continue
		}

		exchangeRates = append(exchangeRates, &exc)
	}
	return
}
