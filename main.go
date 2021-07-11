package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "proxychecker",
		Usage: "fight the loneliness!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "url",
				Usage: "check proxy list from url with extension file is text",
			},
			&cli.StringFlag{
				Name:  "file",
				Usage: "check proxy from file",
			},
			&cli.StringFlag{
				Name:  "ip",
				Usage: "checker ip proxy",
			},
			&cli.StringFlag{
				Name:  "out",
				Value: "active_proxy.txt",
				Usage: "output file",
			},
		},
		Action: func(c *cli.Context) error {
			outfile := c.String("out")
			if c.String("ip") != "" {
				checkProxyFromIP(c.String("ip"))
			}

			if c.String("url") != "" {
				checkProxyFromURL(c.String("url"), outfile)
			}

			if c.String("file") != "" {
				checkProxyFromFile(c.String("file"), outfile)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
