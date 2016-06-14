package main

import ()
import (
	"github.com/schnaidar/stationcrawler"
	"github.com/schnaidar/radiowatch"
	"github.com/codegangsta/cli"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

type config struct {
	Database string `json:"database"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     string `json:"port"`
	User     string `json:"user"`
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}

	app.Action = func(c *cli.Context) {
		var cfg config
		file, err := ioutil.ReadFile("radiowatch.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when reading config file: %v", err.Error())
			os.Exit(1)
		}
		err = json.Unmarshal(file, &cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
		
		writer := radiowatch.NewMysqlWriter(cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		watcher := radiowatch.NewWatcher(writer)
		watcher.SetInterval("20s")

		watcher.AddCrawlers([]radiowatch.Crawler{
			stationcrawler.NewNjoy(),
			stationcrawler.NewNdr2(),
			stationcrawler.NewDasDing(),
			stationcrawler.NewHr3(),
			stationcrawler.NewYouFm(),
			stationcrawler.NewFfn(),
		})

		watcher.StartCrawling()
		channel := make(chan bool)
		<-channel
	}

	app.Run(os.Args)
}
