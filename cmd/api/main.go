package main

import (
	"github.com/SatorNetwork/sator-api/cmd/api/app"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq" // init pg driver
	"log"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}
