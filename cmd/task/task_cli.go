package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/arbrix/go-test/client"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "task cli"
	app.Usage = "cli to work with the `task` microservice"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{"host", "http://localhost:8080", "Task service host", "APP_HOST"},
	}

	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "(title description) create a task",
			Action: func(c *cli.Context) {
				title := c.Args().Get(0)
				desc := c.Args().Get(1)

				host := c.GlobalString("host")

				client := client.TaskClient{Host: host}

				task, err := client.CreateTask(title, desc)
				if err != nil {
					log.Fatal(err)
					return
				}
				fmt.Printf("%+v\n", task)
			},
		},
		{
			Name:  "ls",
			Usage: "list all tasks",
			Action: func(c *cli.Context) {

				host := c.GlobalString("host")

				client := client.TaskClient{Host: host}

				tasks, err := client.GetAllTasks()
				if err != nil {
					log.Fatal(err)
					return
				}
				for _, task := range tasks {
					fmt.Printf("%+v\n", task)
				}
			},
		},
		{
			Name:  "doing",
			Usage: "(id) update a task status to 'doing'",
			Action: func(c *cli.Context) {
				idStr := c.Args().Get(0)
				id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Print(err)
					return
				}

				host := c.GlobalString("host")

				client := client.TaskClient{Host: host}

				task, err := client.UpdateTaskStatus(int64(id), "doing")
				if err != nil {
					log.Fatal(err)
					return
				}
				fmt.Printf("%+v\n", task)
			},
		},
		{
			Name:  "done",
			Usage: "(id) update a task status to 'done'",
			Action: func(c *cli.Context) {
				idStr := c.Args().Get(0)
				id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Print(err)
					return
				}

				host := c.GlobalString("host")

				client := client.TaskClient{Host: host}

				task, err := client.UpdateTaskStatus(int64(id), "done")
				if err != nil {
					log.Fatal(err)
					return
				}
				fmt.Printf("%+v\n", task)
			},
		},
		{
			Name:  "save",
			Usage: "(id title description) update a task title and description",
			Action: func(c *cli.Context) {
				idStr := c.Args().Get(0)
				id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Print(err)
					return
				}
				title := c.Args().Get(1)
				desc := c.Args().Get(2)

				host := c.GlobalString("host")

				client := client.TaskClient{Host: host}

				task, err := client.GetTask(int64(id))
				if err != nil {
					log.Fatal(err)
					return
				}

				task.Title = title
				task.Description = desc

				task2, err := client.UpdateTask(task)
				if err != nil {
					log.Fatal(err)
					return
				}

				fmt.Printf("%+v\n", task2)
			},
		},
	}
	app.Run(os.Args)

}
