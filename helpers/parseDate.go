/*  parseDate.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 10, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 10/09/18 10:06
 */

package helpers

import "time"

// ParseDate of exchangeRate
func ParseDate(date string) (d time.Time, err error) {
	// parse string to date
	d, err = time.Parse("2006-01-02", date)
	if err != nil {
		return
	}
	return
}
