/*  modek.go.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 12, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 12/09/18 02:46
 */

package exchangeRateTrack

import (
	"strings"
)

// ExchangeRateTrack define model exchange rate
type ExchangeRateTrack struct {
	id   string
	from string
	to   string
}

// IDGenerator returns next id.
func IDGenerator(from, to string) string {
	return strings.ToLower(from + to)
}

// SetID of ExchangeRateTrack
func (acc *ExchangeRateTrack) SetID(id string) {
	acc.id = id
}

// SetFrom of ExchangeRateTrack
func (acc *ExchangeRateTrack) SetFrom(from string) {
	acc.from = from
}

// SetTo Of ExchangeRateTrack
func (acc *ExchangeRateTrack) SetTo(to string) {
	acc.to = to
}

// ID of ExchangeRateTrack
func (acc ExchangeRateTrack) ID() string {
	return acc.id
}

// From of ExchangeRateTrack
func (acc ExchangeRateTrack) From() string {
	return acc.from
}

// To of ExchangeRateTrack
func (acc ExchangeRateTrack) To() string {
	return acc.to
}
