package main

import (
	"fmt"
	"github.com/arbrix/go-test/app"
	"log"
	"path/filepath"
)

func main() {

	a := app.NewApp(&app.AppConfig{}, &app.AppOrm{})

	err = a.Run()
	if err != nil {
		fmt.Printf("%v", err.Error())
	}
}
