/*  viewmodel.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 10, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 10/09/18 12:42
 */

package exchangeRate

// Define ViewModel of this service

// ViewModelListExchangeRate of exchange rate
type ViewModelListExchangeRate struct {
	From          string      `json:"from"`
	To            string      `json:"to"`
	Rate          interface{} `json:"rate"`
	SevenDayRates interface{} `json:"seven_day_rates"`
}

// ViewModelRate of exchange rate
type ViewModelRate struct {
	Date string      `json:"date"`
	Rate interface{} `json:"rate"`
}

// ViewModelTrendExchangeRate of exchange rate
type ViewModelTrendExchangeRate struct {
	From     string          `json:"from"`
	To       string          `json:"to"`
	Average  interface{}     `json:"average"`
	Variance interface{}     `json:"variance"`
	Rates    []ViewModelRate `json:"rates"`
}
