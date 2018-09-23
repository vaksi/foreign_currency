/*  CLI.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 08, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 08/09/18 03:47
 */

package helpers

import (
	"time"

	"github.com/urfave/cli"
)

// StringFlag gets string flag
func StringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

// BoolFlag gets boolean flag
func BoolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

// IntFlag gets integer flag
func IntFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

// DurationFlag gets duration flag
func DurationFlag(name string, value time.Duration, usage string) cli.DurationFlag {
	return cli.DurationFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
