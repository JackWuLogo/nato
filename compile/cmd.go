package compile

import (
	"github.com/micro/cli/v2"
	"os"
)

var Cmd = []*cli.Command{
	{
		Name:    "version",
		Aliases: []string{"ver"},
		Usage:   "show command version info",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "show command version info format json",
			},
		},
		Action: func(c *cli.Context) error {
			if c.String("json") == "true" {
				EchoVersionJson()
			} else {
				EchoVersion(nil)
			}
			os.Exit(0)
			return nil
		},
	},
}
