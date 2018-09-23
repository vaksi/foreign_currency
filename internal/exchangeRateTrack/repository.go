/*  repository.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 12, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 12/09/18 02:52
 */

package exchangeRateTrack

import (
	"fmt"

	"github.com/vaksi/foreign_currency/infrastructures"
)

// Repository of struct exchange rate tracks
type Repository struct {
	DB infrastructures.MySQLFactory
}

// RepositoryFactory of interface repository tracks
type RepositoryFactory interface {
	Store(*ExchangeRateTrack) error
	Delete(id string) error
}

// NewRepositoryFactory of new factory repository
func NewRepositoryFactory(db infrastructures.MySQLFactory) RepositoryFactory {
	return &Repository{
		DB: db,
	}
}

// Store is function to store exchange rate tracker
func (r *Repository) Store(exchangeRateTrack *ExchangeRateTrack) (err error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}

	// prepare transaction
	stmt, err := tx.Prepare("INSERT INTO tracks (uid,`from`,`to`) VALUES (?,?,?)")
	if err != nil {
		return
	}
	_, err = stmt.Exec(exchangeRateTrack.ID(), exchangeRateTrack.From(), exchangeRateTrack.To())
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
	fmt.Println("success")
	return
}

// Delete this function to delete exchange rate tracker
func (r *Repository) Delete(id string) (err error) {
	db, err := r.DB.GetDB()
	if err != nil {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}

	// prepare transaction
	stmt, err := tx.Prepare("DELETE FROM tracks WHERE uid=?")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		// roll back db if error
		if errTx := tx.Rollback(); errTx != nil {
			return errTx
		}
		return err
	}
	aft, err := res.RowsAffected()
	if err != nil {
		return
	}
	if aft < 1 {
		return fmt.Errorf("data not available")
	}

	// commit transaction
	if errTx := tx.Commit(); errTx != nil {
		return errTx
	}

	return
}
