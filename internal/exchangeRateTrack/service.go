/*  service.go.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 12, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 12/09/18 02:52
 */

package exchangeRateTrack

import (
	log "github.com/sirupsen/logrus"
)

// Service this struct for service Exchange Rate Tracks
type Service struct {
	ExchangeRateTrackRepo RepositoryFactory
}

// AddExchangeRateToTrack of exchange rate track service
func (s *Service) AddExchangeRateToTrack(from, to string) (err error) {
	id := IDGenerator(from, to)
	exchangeRateTrack := new(ExchangeRateTrack)
	exchangeRateTrack.id = id
	exchangeRateTrack.from = from
	exchangeRateTrack.to = to
	if err = s.ExchangeRateTrackRepo.Store(exchangeRateTrack); err != nil {
		log.Error(err)
		return err
	}
	return
}

// DeleteExchangeRateToTrack of exchange rate track service
func (s *Service) DeleteExchangeRateToTrack(from, to string) (err error) {
	id := IDGenerator(from, to)
	if err = s.ExchangeRateTrackRepo.Delete(id); err != nil {
		log.Error(err)
		return err
	}
	return
}
