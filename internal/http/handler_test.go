/*  handler.go.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 05, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 05/09/18 03:11
 */

package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	gohttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"github.com/vaksi/foreign_currency/helpers"
	"github.com/vaksi/foreign_currency/internal/exchangeRate"
	exchangeRateMocks "github.com/vaksi/foreign_currency/internal/exchangeRate/mocks"
	"github.com/vaksi/foreign_currency/internal/exchangeRateTrack"
	exchangeRateTrackMocks "github.com/vaksi/foreign_currency/internal/exchangeRateTrack/mocks"
)

func request(method, path string, body interface{}) *gohttp.Request {
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(b)

	req := httptest.NewRequest(method, path, buf)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func TestHandler_ping(t *testing.T) {
	exchangeRateSvc := &exchangeRate.Service{}
	exchangeRateTrackSvc := &exchangeRateTrack.Service{}
	Convey("When ping status not valid routes endpoint", t, func() {
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)
		// endpoint
		req := httptest.NewRequest("GET", "/pingfalse", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusNotFound)
	})

	Convey("When Ping status ok", t, func() {
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)
		// endpoint
		req := httptest.NewRequest("GET", "/foreign-currency/v1/ping", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// Expectation
		exp, _ := json.Marshal(helpers.APIOK)

		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusOK)
		So(rec.Body.String(), ShouldEqual, string(exp))
	})
}

func Test_delegate_CreateDailyExchangeRate(t *testing.T) {
	exchangeRateTrackSvc := &exchangeRateTrack.Service{}
	Convey("When CreateDailyExchangeRate status not valid routes endpoint", t, func() {
		exchangeRateRepo := new(exchangeRateMocks.RepositoryFactory)
		exchangeRateSvc := &exchangeRate.Service{ExchangeRateRepo: exchangeRateRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		// endpoint
		req := httptest.NewRequest("POST", "/foreign-currency/exchange-rates-false", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusNotFound)
	})

	Convey("When CreateDailyExchangeRate error internal server", t, func() {
		// mocking
		exchangeRateRepo := new(exchangeRateMocks.RepositoryFactory)
		exchangeRateRepo.On("Store", mock.Anything).Return(errors.New("Something Error"))

		exchangeRateSvc := &exchangeRate.Service{ExchangeRateRepo: exchangeRateRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		date := "2018-07-01"
		from := "USD"
		to := "IDR"
		rate := 14000.0
		payload := map[string]interface{}{
			"date": date,
			"from": from,
			"to":   to,
			"rate": rate,
		}

		// endpoint
		req := request("POST", "/foreign-currency/v1/exchange-rates", payload)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		exp, _ := json.Marshal(helpers.APIErrorUnknown)

		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusInternalServerError)
		So(rec.Body.String(), ShouldEqual, string(exp))
	})

	Convey("When CreateDailyExchangeRate status success", t, func() {
		// mocking
		exchangeRateRepo := new(exchangeRateMocks.RepositoryFactory)
		exchangeRateRepo.On("Store", mock.Anything).Return(nil)

		exchangeRateSvc := &exchangeRate.Service{ExchangeRateRepo: exchangeRateRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		date := "2018-07-01"
		from := "USD"
		to := "IDR"
		rate := 14000.0
		payload := map[string]interface{}{
			"date": date,
			"from": from,
			"to":   to,
			"rate": rate,
		}

		// endpoint
		req := request("POST", "/foreign-currency/v1/exchange-rates", payload)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		exp, _ := json.Marshal(helpers.APICreated)
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusCreated)
		So(rec.Body.String(), ShouldEqual, string(exp))
	})

	Convey("When CreateDailyExchangeRate error invalid value", t, func() {
		// mocking
		exchangeRateRepo := new(exchangeRateMocks.RepositoryFactory)
		exchangeRateRepo.On("Store", mock.Anything).Return(nil)

		exchangeRateSvc := &exchangeRate.Service{ExchangeRateRepo: exchangeRateRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		date := "2018-07-01 ef"
		from := "USDtest"
		to := "IDRtest"
		rate := 0
		payload := map[string]interface{}{
			"date": date,
			"from": from,
			"to":   to,
			"rate": rate,
		}

		// endpoint
		req := request("POST", "/foreign-currency/v1/exchange-rates", payload)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		r := helpers.APIErrorValidation
		r.Errors = map[string]interface{}{
			"date": "must be a valid date",
			"from": "Value USDtest is not valid in Currency",
			"to":   "Value IDRtest is not valid in Currency",
		}
		exp, _ := json.Marshal(r)

		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusBadRequest)
		So(rec.Body.String(), ShouldEqual, string(exp))
	})
}

// parse string to date
func parseDate(date string) (d time.Time) {
	d, _ = time.Parse("2006-01-02", date)
	return
}

func Test_handler_getListExchangeRate(t *testing.T) {
	exchangeRateTrackSvc := &exchangeRateTrack.Service{}
	// Set Value ExchageRates
	var list []*exchangeRate.ExchangeRate
	// exc1
	exc1 := &exchangeRate.ExchangeRate{}
	exc1.SetFrom("USD")
	exc1.SetTo("IDR")
	exc1.SetRate(14347)
	exc1.SetSevenDayRates(14289)
	list = append(list, exc1)

	Convey("When getListExchangeRate Success", t, func() {
		// mocking
		exchangeRateRepo := new(exchangeRateMocks.RepositoryFactory)
		exchangeRateRepo.On("WhereByDate", mock.Anything).Return(list, nil)

		exchangeRateSvc := &exchangeRate.Service{ExchangeRateRepo: exchangeRateRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		// endpoint
		req := request("GET", "/foreign-currency/v1/exchange-rates?date=2018-07-02", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusOK)
	})

	Convey("When getListExchangeRate error validation param", t, func() {
		// mocking
		exchangeRateRepo := new(exchangeRateMocks.RepositoryFactory)
		exchangeRateRepo.On("WhereByDate", mock.Anything).Return(list, nil)

		exchangeRateSvc := &exchangeRate.Service{ExchangeRateRepo: exchangeRateRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		// endpoint
		req := request("GET", "/foreign-currency/v1/exchange-rates?date=2018-07-02test", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusBadRequest)
	})

	Convey("When getListExchangeRate error internals", t, func() {
		// mocking
		exchangeRateRepo := new(exchangeRateMocks.RepositoryFactory)
		exchangeRateRepo.On("WhereByDate", mock.Anything).Return(nil, errors.New("any errors"))

		exchangeRateSvc := &exchangeRate.Service{ExchangeRateRepo: exchangeRateRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		// endpoint
		req := request("GET", "/foreign-currency/v1/exchange-rates?date=2018-07-02", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusInternalServerError)
	})
}

func Test_handler_getTrendExchangeRate(t *testing.T) {
	exchangeRateTrackSvc := &exchangeRateTrack.Service{}
	// Set Value ExchageRates
	var list []*exchangeRate.ExchangeRate

	// exc1
	exc1 := &exchangeRate.ExchangeRate{}
	exc1.SetFrom("GBP")
	exc1.SetTo("USD")
	exc1.SetDate(parseDate("2018-07-08"))
	exc1.SetRate(1.417)
	list = append(list, exc1)

	// exc2
	exc2 := &exchangeRate.ExchangeRate{}
	exc2.SetFrom("GBP")
	exc2.SetTo("USD")
	exc2.SetDate(parseDate("2018-07-07"))
	exc2.SetRate(1.295)
	list = append(list, exc2)

	// exc3
	exc3 := &exchangeRate.ExchangeRate{}
	exc3.SetFrom("GBP")
	exc3.SetTo("USD")
	exc3.SetDate(parseDate("2018-07-06"))
	exc3.SetRate(1.199)
	list = append(list, exc3)

	Convey("When getTrendExchangeRate Success", t, func() {
		// mocking
		exchangeRateRepo := new(exchangeRateMocks.RepositoryFactory)
		exchangeRateRepo.On("WhereByTrend", "USD", "GBP").Return(list, nil)

		exchangeRateSvc := &exchangeRate.Service{ExchangeRateRepo: exchangeRateRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		// endpoint
		req := request("GET", "/foreign-currency/v1/exchange-rates/trend?from=USD&to=GBP", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusOK)
	})

	Convey("When getTrendExchangeRate Error Validation", t, func() {
		// mocking
		exchangeRateRepo := new(exchangeRateMocks.RepositoryFactory)

		exchangeRateSvc := &exchangeRate.Service{ExchangeRateRepo: exchangeRateRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		// endpoint
		req := request("GET", "/foreign-currency/v1/exchange-rates/trend?from=USDtest&to=GBPtest", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusBadRequest)
	})

	Convey("When getTrendExchangeRate Data Not Found", t, func() {
		// mocking
		exchangeRateRepo := new(exchangeRateMocks.RepositoryFactory)
		exchangeRateRepo.On("WhereByTrend", "GBP", "USD").Return(nil, nil)

		exchangeRateSvc := &exchangeRate.Service{ExchangeRateRepo: exchangeRateRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		// endpoint
		req := request("GET", "/foreign-currency/v1/exchange-rates/trend?from=GBP&to=USD", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		resp := helpers.APIDataNotFound
		resp.Data = "Trend Not Found"
		exp, _ := json.Marshal(resp)
		buf := bytes.NewBuffer(exp)
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusOK)
		So(rec.Body.String(), ShouldEqual, buf.String())
	})
}

func Test_handler_addExchangeRateToTrack(t *testing.T) {
	exchangeRateSvc := &exchangeRate.Service{}
	Convey("When addExchangeRateToTrack status success", t, func() {
		// mocking
		exchangeRateTrackRepo := new(exchangeRateTrackMocks.RepositoryFactory)
		exchangeRateTrackRepo.On("Store", mock.Anything).Return(nil)

		exchangeRateTrackSvc := &exchangeRateTrack.Service{ExchangeRateTrackRepo: exchangeRateTrackRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		from := "USD"
		to := "IDR"

		payload := map[string]interface{}{
			"from": from,
			"to":   to,
		}

		// endpoint
		req := request("POST", "/foreign-currency/v1/exchange-rates/tracks", payload)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		exp, _ := json.Marshal(helpers.APICreated)
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusCreated)
		So(rec.Body.String(), ShouldEqual, string(exp))
	})

	Convey("When addExchangeRateToTrack status invalid data", t, func() {
		// mocking
		exchangeRateTrackRepo := new(exchangeRateTrackMocks.RepositoryFactory)

		exchangeRateTrackSvc := &exchangeRateTrack.Service{ExchangeRateTrackRepo: exchangeRateTrackRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		from := "USDTest"
		to := "IDRTest"

		payload := map[string]interface{}{
			"from": from,
			"to":   to,
		}

		// endpoint
		req := request("POST", "/foreign-currency/v1/exchange-rates/tracks", payload)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusBadRequest)
	})

	Convey("When addExchangeRateToTrack status Error Internal", t, func() {
		// mocking
		exchangeRateTrackRepo := new(exchangeRateTrackMocks.RepositoryFactory)
		exchangeRateTrackRepo.On("Store", mock.Anything).Return(errors.New("any errors"))

		exchangeRateTrackSvc := &exchangeRateTrack.Service{ExchangeRateTrackRepo: exchangeRateTrackRepo}

		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		from := "USD"
		to := "IDR"

		payload := map[string]interface{}{
			"from": from,
			"to":   to,
		}

		// endpoint
		req := request("POST", "/foreign-currency/v1/exchange-rates/tracks", payload)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusInternalServerError)
	})

	Convey("When addExchangeRateToTrack status Duplicate Expression", t, func() {
		// mocking
		exchangeRateTrackRepo := new(exchangeRateTrackMocks.RepositoryFactory)
		exchangeRateTrackRepo.On("Store", mock.Anything).Return(&mysql.MySQLError{
			Number:  uint16(1062),
			Message: "duplicate data",
		})

		exchangeRateTrackSvc := &exchangeRateTrack.Service{ExchangeRateTrackRepo: exchangeRateTrackRepo}

		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		from := "USD"
		to := "IDR"

		payload := map[string]interface{}{
			"from": from,
			"to":   to,
		}

		// endpoint
		req := request("POST", "/foreign-currency/v1/exchange-rates/tracks", payload)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusOK)
	})
}

func Test_handler_deleteExchangeRateToTrack(t *testing.T) {
	exchangeRateSvc := &exchangeRate.Service{}
	Convey("When deleteExchangeRateToTrack status success", t, func() {
		// mocking
		exchangeRateTrackRepo := new(exchangeRateTrackMocks.RepositoryFactory)
		exchangeRateTrackRepo.On("Delete", "gbpusd").Return(nil)

		exchangeRateTrackSvc := &exchangeRateTrack.Service{ExchangeRateTrackRepo: exchangeRateTrackRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		// endpoint
		req := request("DELETE", "/foreign-currency/v1/exchange-rates/tracks?from=GBP&to=USD", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		exp, _ := json.Marshal(helpers.APIAccepted)
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusAccepted)
		So(rec.Body.String(), ShouldEqual, string(exp))
	})

	Convey("When deleteExchangeRateToTrack status invalid data", t, func() {
		// mocking
		exchangeRateTrackRepo := new(exchangeRateTrackMocks.RepositoryFactory)

		exchangeRateTrackSvc := &exchangeRateTrack.Service{ExchangeRateTrackRepo: exchangeRateTrackRepo}
		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		// endpoint
		req := request("DELETE", "/foreign-currency/v1/exchange-rates/tracks?from=GBPtes&to=USDtes", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusBadRequest)
	})

	Convey("When deleteExchangeRateToTrack status Error Internal", t, func() {
		// mocking
		exchangeRateTrackRepo := new(exchangeRateTrackMocks.RepositoryFactory)
		exchangeRateTrackRepo.On("Delete", mock.Anything).Return(errors.New("any errors"))

		exchangeRateTrackSvc := &exchangeRateTrack.Service{ExchangeRateTrackRepo: exchangeRateTrackRepo}

		h := NewHandler(exchangeRateSvc, exchangeRateTrackSvc)

		// endpoint
		req := request("DELETE", "/foreign-currency/v1/exchange-rates/tracks?from=GBP&to=USD", nil)
		rec := httptest.NewRecorder()

		// exec
		h.ServeHTTP(rec, req)

		// expectation
		fmt.Println(rec.Body.String())
		// assert/check result
		So(rec.Code, ShouldEqual, gohttp.StatusInternalServerError)
	})

}
