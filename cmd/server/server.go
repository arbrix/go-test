package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/arbrix/go-test/service"
	"github.com/codegangsta/cli"
)

func getConfig(c *cli.Context) (service.Config, error) {
	jsonPath := c.GlobalString("config")
	config := service.Config{}

	if _, err := os.Stat(jsonPath); err != nil {
		return config, errors.New("config path not valid")
	}

	file, err := os.Open(jsonPath)
	if err != nil {
		return config, err
	}
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)
	return config, err
}

func main() {

	app := cli.NewApp()
	app.Name = "task"
	app.Usage = "work with the `task` microservice"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{"config, c", "config.json", "config file to use", "APP_CONFIG"},
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run the http server",
			Action: func(c *cli.Context) {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return
				}

				svc := service.TaskService{}

				if err = svc.Run(cfg); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "migratedb",
			Usage: "Perform database migrations",
			Action: func(c *cli.Context) {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return
				}

				svc := service.TaskService{}

				if err = svc.Migrate(cfg); err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)

}
