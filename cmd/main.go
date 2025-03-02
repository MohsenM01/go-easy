package main

import (
	"go-easy/cmd/api"
	"go-easy/config"
	"go-easy/infrastructure/database"
)

func main() {

	config.LoadConfig()
	database.InitDb()
	api.StartApiServer()
}
