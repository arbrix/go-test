package main

import (
	"fmt"
	"github.com/arbrix/go-test/app"
	"log"
	"path/filepath"
)

func main() {
	dir, err := filepath.Abs("./config/")
	if err != nil {
		log.Fatal(err)
	}
	a := app.NewApp(&app.AppConfig{}, &app.AppOrm{})
	a.GetConfig().SetBasePath(dir)

	err = a.Run()
	if err != nil {
		fmt.Printf("%v", err)
	}
}
