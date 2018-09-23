/*  exchange.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 05, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 05/09/18 02:49
 */

package exchangeRate

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// ExchangeRate define model exchange rate
type ExchangeRate struct {
	id            string
	date          time.Time
	from          string
	to            string
	rate          float32
	sevenDayRates float32
}

// ExchangeTrack define model exchange rate tracks
type ExchangeTrack struct {
	id   string
	from string
	to   string
}

// NextID returns next id.
func NextID() string {
	return bson.NewObjectId().Hex()
}

// SetID of exchangeRate
func (acc *ExchangeRate) SetID(id string) {
	acc.id = id
}

// SetDate of exchangeRate
func (acc *ExchangeRate) SetDate(date time.Time) {
	acc.date = date
}

// SetFrom of exchangeRate
func (acc *ExchangeRate) SetFrom(from string) {
	acc.from = from
}

// SetTo Of exchangeRate
func (acc *ExchangeRate) SetTo(to string) {
	acc.to = to
}

// SetRate of exchangeRate
func (acc *ExchangeRate) SetRate(rate float32) {
	acc.rate = rate
}

// SetSevenDayRates of exchangeRate
func (acc *ExchangeRate) SetSevenDayRates(sevenDayRates float32) {
	acc.sevenDayRates = sevenDayRates
}

// ID of exchangeRate
func (acc ExchangeRate) ID() string {
	return acc.id
}

// Date of ExchangeRate
func (acc ExchangeRate) Date() time.Time {
	return acc.date
}

// From of ExchangeRate
func (acc ExchangeRate) From() string {
	return acc.from
}

// To of ExchangeRate
func (acc ExchangeRate) To() string {
	return acc.to
}

// Rate of ExchangeRate
func (acc ExchangeRate) Rate() float32 {
	return acc.rate
}

// SevenDayRates of ExchangeRate
func (acc ExchangeRate) SevenDayRates() float32 {
	return acc.sevenDayRates
}
