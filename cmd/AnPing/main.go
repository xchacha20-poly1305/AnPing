package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/urfave/cli/v2"
	anping "github.com/xchacha20-poly1305/AnPing"
)

var version string = "Unknow"

func main() {
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Println("AnPing")
		fmt.Println("Version: " + cCtx.App.Version)
	}
	
	app := &cli.App{
		
		Name:    "AnPing",
		Version: version,
		Usage:   "Ping whatever you like.",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "count",
				Aliases:     []string{"c"},
				Value:       -1,
				Destination: &count,
			},
		},
		Action: func(cCtx *cli.Context) error {
			
			a, err := anping.New(cCtx.Args().Get(0))
			if err != nil {
				log.Println("Failed to start")
				return err
			}
			a.Count = count
			
			// log.Println(a.Count)
			
			a.Start()
			return nil
		},
	}
	
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var (
	count int = -1
)
