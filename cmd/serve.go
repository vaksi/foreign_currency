/*  Serve.go.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 08, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 08/09/18 03:40
 */

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"net/http"

	config "github.com/spf13/viper"
	"github.com/urfave/cli"
	"github.com/vaksi/foreign_currency/helpers"
	"github.com/vaksi/foreign_currency/infrastructures"
	"github.com/vaksi/foreign_currency/injector"
	handler "github.com/vaksi/foreign_currency/internal/http"
)

var text = `
========================================================
Foreign Currency
========================================================
- port    : %d
--------------------------------------------------------`

// Serve serves the service
var Serve = cli.Command{
	Name:   "serve",
	Usage:  "Used to run the service",
	Action: RunServer,
	Flags: []cli.Flag{
		helpers.StringFlag("config, c", "configs/app.yaml", "Custom configuration file path"),
	},
}

// RunServer serves the service
func RunServer(c *cli.Context) {
	// Set application configuration
	if err := infrastructures.SetConfiguration(c.String("configs")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set Dependencies
	inj := injector.NewDependencyInjection()
	inj.Database = injector.InitSQL()

	inj.Inject()

	// Set router
	route := handler.NewHandler(&inj.ExchangeRateService, &inj.ExchangeRateTrackService)
	s := &http.Server{
		Addr:         ":" + config.GetString("app.port"),
		Handler:      route,
		ReadTimeout:  time.Duration(config.GetInt("app.read_timeout")) * time.Second,
		WriteTimeout: time.Duration(config.GetInt("app.write_timeout")) * time.Second,
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		fmt.Println("Shutting down server...")
		if err := s.Shutdown(context.Background()); err != nil {
			fmt.Println("could not shutdown: ", err.Error())
		}
	}()

	// Listen and serve
	fmt.Println(fmt.Sprintf(text, config.GetInt("app.port")))
	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		fmt.Println("listen: ", err.Error())
	}

	fmt.Println("Server gracefully stopped")
}
