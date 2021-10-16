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
			&cli.StringFlag{
				Name:  "target",
				Usage: "target site access with proxy",
			},
			&cli.BoolFlag{
				Name:  "sock5",
				Usage: "this option for proxy sock5",
			},
			&cli.BoolFlag{
				Name:  "sock4",
				Usage: "this option for proxy sock5",
			},
		},
		Action: func(c *cli.Context) error {
			outfile := c.String("out")

			var target string
			var sock string
			if c.Bool("sock5") {
				sock = "sock5"
			}

			if c.Bool("sock4") {
				sock = "sock4"
			}

			if c.String("target") != "" {
				target = c.String("target")
			}

			if c.String("ip") != "" {
				checkProxyFromIP(c.String("ip"), sock, target)
			}

			if c.String("url") != "" {
				checkProxyFromURL(c.String("url"), outfile, sock, target)
			}

			if c.String("file") != "" {
				checkProxyFromFile(c.String("file"), outfile, sock, target)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
