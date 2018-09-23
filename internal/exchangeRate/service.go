/*  services.go.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 05, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 05/09/18 03:04
 */

package exchangeRate

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vaksi/foreign_currency/constants"
)

// Service this struct for service Exchange Rate
type Service struct {
	ExchangeRateRepo RepositoryFactory
}

// ConstraintErrorNotFound of this service
type ConstraintErrorNotFound string

func (c ConstraintErrorNotFound) Error() string {
	return string(c)
}

// CreateNewDailyExchangeRate of Create New exchange rate service
func (s *Service) CreateNewDailyExchangeRate(date time.Time, from, to string, rate float32) (err error) {
	id := NextID()

	exchangeRate := new(ExchangeRate)
	exchangeRate.id = id
	exchangeRate.date = date
	exchangeRate.from = from
	exchangeRate.to = to
	exchangeRate.rate = rate

	if err = s.ExchangeRateRepo.Store(exchangeRate); err != nil {
		log.Error(err)
		return
	}

	return
}

// GetListDailyExchangeRate of exchange rates service
func (s *Service) GetListDailyExchangeRate(date time.Time) (exchangeRateList []ViewModelListExchangeRate, err error) {
	// GetData
	exchangeRates, err := s.ExchangeRateRepo.WhereByDate(date)
	if err != nil {
		log.Error(err)
		return
	}

	// presenter data
	for _, exc := range exchangeRates {
		var rate, sevenDayRates interface{}
		rate = exc.Rate()
		sevenDayRates = exc.SevenDayRates()
		if exc.Rate() == 0 {
			rate = constants.MissingData
		}
		if exc.SevenDayRates() == 0 {
			sevenDayRates = constants.MissingData
		}
		exchangeRateList = append(exchangeRateList, ViewModelListExchangeRate{
			From:          exc.From(),
			To:            exc.To(),
			Rate:          rate,
			SevenDayRates: sevenDayRates,
		})
	}
	return
}

// GetTrendExchangeRate of exchange rates service
func (s *Service) GetTrendExchangeRate(from, to string) (trend ViewModelTrendExchangeRate, err error) {
	// get data
	exchangeRates, err := s.ExchangeRateRepo.WhereByTrend(from, to)
	if err != nil {
		log.Error(err)
		return
	}

	if len(exchangeRates) <= 0 {
		err = ConstraintErrorNotFound("Trend Not Found")
		return
	}
	var (
		total float32
		max   = exchangeRates[0].rate
		min   = exchangeRates[0].rate
	)
	for _, rate := range exchangeRates {
		trend.From = rate.from
		trend.To = rate.to
		trend.Rates = append(trend.Rates, ViewModelRate{
			Date: rate.date.Format("2006-01-02"),
			Rate: rate.rate,
		})
		total += rate.rate
		if max < rate.rate {
			max = rate.rate
		}
		if min > rate.rate {
			min = rate.rate
		}
	}
	trend.Average = total / float32(len(exchangeRates))
	trend.Variance = max - min
	return
}
