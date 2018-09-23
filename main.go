/*  main.go.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 08, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 08/09/18 03:33
 */

package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/vaksi/foreign_currency/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "API Foreign Currency"
	app.Usage = "Used as service for Foreign Currency"
	app.UsageText = "[global options] command [arguments]"
	app.Version = "1.0"
	app.Commands = []cli.Command{
		cmd.Serve,
	}

	if err := app.Run(os.Args); err != nil {
		cli.OsExiter(1)
	}
}
