/*  handler.go.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 05, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 05/09/18 03:11
 */

package http

import (
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/vaksi/foreign_currency/helpers"
	"github.com/vaksi/foreign_currency/internal/exchangeRate"
	"github.com/vaksi/foreign_currency/internal/exchangeRateTrack"
)

type handler struct {
	exchangeRateService      *exchangeRate.Service
	exchangeRateTrackService *exchangeRateTrack.Service
}

type exchangeRatePayload struct {
	Date string  `json:"date"`
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float32 `json:"rate"`
}

type exchangeRateToTrackPayload struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (h *handler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, helpers.APIOK)
	return
}

func (h *handler) createDailyExchangeRate(c *gin.Context) {
	var (
		payload exchangeRatePayload
		resp    helpers.APIResponse
	)
	if err := c.Bind(&payload); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, helpers.APIErrorInvalidData)
	}

	// check parallel validation
	err := validation.Errors{
		"date": validation.Validate(payload.Date, validation.Required, validation.Date("2006-01-02")),
		"from": validation.Validate(payload.From, validation.Required, validation.By(helpers.ValidateInValueCurrency)),
		"to":   validation.Validate(payload.To, validation.Required, validation.By(helpers.ValidateInValueCurrency), validation.NotIn(payload.From).Error("parameter 'to' don't same as 'from' ")),
		"rate": validation.Validate(fmt.Sprintf("%f", payload.Rate), validation.Required, is.Float),
	}.Filter()

	if err != nil {
		log.Error(err)
		resp = helpers.APIErrorValidation
		resp.Errors = err
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Parse date
	date, err := helpers.ParseDate(payload.Date)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, helpers.APIErrorUnknown)
		return
	}

	if err := h.exchangeRateService.CreateNewDailyExchangeRate(date, payload.From, payload.To, payload.Rate); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIErrorUnknown)
		return
	}

	c.JSON(http.StatusCreated, helpers.APICreated)
	return
}

func (h *handler) getListExchangeRate(c *gin.Context) {
	var resp helpers.APIResponse

	// get param date
	date := c.Request.URL.Query().Get("date")

	// validate date
	if err := validation.Validate(date, validation.Required, validation.Date("2006-01-02")); err != nil {
		resp = helpers.APIErrorValidation
		resp.Errors = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Parse String Date to time.Time
	d, err := helpers.ParseDate(date)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, helpers.APIErrorUnknown)
		return
	}

	// retrieve data from service
	data, err := h.exchangeRateService.GetListDailyExchangeRate(d)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIErrorUnknown)
		return
	}

	// response success
	resp = helpers.APIOK
	resp.Data = data
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) getTrendExchangeRate(c *gin.Context) {
	var resp helpers.APIResponse

	// get param
	params := c.Request.URL.Query()
	from := params.Get("from")
	to := params.Get("to")

	// validate
	err := validation.Errors{
		"from": validation.Validate(from, validation.Required, validation.By(helpers.ValidateInValueCurrency)),
		"to":   validation.Validate(to, validation.Required, validation.By(helpers.ValidateInValueCurrency), validation.NotIn(from).Error("parameter 'to' don't same as 'from' ")),
	}.Filter()
	if err != nil {
		log.Error(err)
		resp = helpers.APIErrorValidation
		resp.Errors = err
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// retrieve data from service
	data, err := h.exchangeRateService.GetTrendExchangeRate(from, to)
	if err != nil {
		if errNotFound, ok := err.(exchangeRate.ConstraintErrorNotFound); ok {
			resp = helpers.APIDataNotFound
			resp.Data = errNotFound.Error()
			c.JSON(http.StatusOK, resp)
			return
		}
		c.JSON(http.StatusInternalServerError, helpers.APIErrorUnknown)
		return
	}

	// response success
	resp = helpers.APIOK
	resp.Data = data
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) addExchangeRateToTrack(c *gin.Context) {
	var (
		resp    helpers.APIResponse
		payload exchangeRateToTrackPayload
	)

	if err := c.Bind(&payload); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, helpers.APIErrorInvalidData)
	}

	// validate param
	err := validation.Errors{
		"from": validation.Validate(payload.From, validation.Required, validation.By(helpers.ValidateInValueCurrency)),
		"to":   validation.Validate(payload.To, validation.Required, validation.By(helpers.ValidateInValueCurrency), validation.NotIn(payload.From).Error("parameter 'to' don't same as 'from' ")),
	}.Filter()
	if err != nil {
		log.Error(err)
		resp = helpers.APIErrorValidation
		resp.Errors = err
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// store data
	if err := h.exchangeRateTrackService.AddExchangeRateToTrack(payload.From, payload.To); err != nil {
		me, _ := err.(*mysql.MySQLError)
		if me.Number == 1062 {
			resp = helpers.APIOKDuplicatedExpression
			resp.Message = me.Message
			c.JSON(http.StatusOK, resp)
			return
		}
		resp = helpers.APIErrorUnknown
		resp.Errors = err
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusCreated, helpers.APICreated)
	return
}

func (h *handler) deleteExchangeRateToTrack(c *gin.Context) {
	var (
		resp helpers.APIResponse
	)

	// get param
	params := c.Request.URL.Query()
	from := params.Get("from")
	to := params.Get("to")

	// validate
	err := validation.Errors{
		"from": validation.Validate(from, validation.Required, validation.By(helpers.ValidateInValueCurrency)),
		"to":   validation.Validate(to, validation.Required, validation.By(helpers.ValidateInValueCurrency), validation.NotIn(from).Error("parameter 'to' don't same as 'from' ")),
	}.Filter()
	if err != nil {
		log.Error(err)
		resp = helpers.APIErrorValidation
		resp.Errors = err
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// store data
	if err := h.exchangeRateTrackService.DeleteExchangeRateToTrack(from, to); err != nil {
		resp = helpers.APIErrorUnknown
		resp.Errors = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusAccepted, helpers.APICreated)
	return
}

// NewHandler constructs new http.Handler.
func NewHandler(exchangeRateSvc *exchangeRate.Service, exchangeRateTrackSvc *exchangeRateTrack.Service) http.Handler {
	h := &handler{exchangeRateService: exchangeRateSvc, exchangeRateTrackService: exchangeRateTrackSvc}

	// Initialize router
	router := gin.Default()
	route := router.Group("/foreign-currency/v1")
	{
		route.GET("/ping", h.ping)
		route.POST("/exchange-rates", h.createDailyExchangeRate)
		route.GET("/exchange-rates", h.getListExchangeRate)
		route.GET("/exchange-rates/trend", h.getTrendExchangeRate)
		route.POST("/exchange-rates/tracks", h.addExchangeRateToTrack)
		route.DELETE("/exchange-rates/tracks", h.deleteExchangeRateToTrack)
	}
	//
	// router.Run(":8081")
	return router
}
