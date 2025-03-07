package main

import (
	"go-easy/cmd/api"
	"go-easy/config"
	"go-easy/infrastructure/cache"
	"go-easy/infrastructure/database"
	"go-easy/infrastructure/messaging"
	httpclient "go-easy/internal/delivery/httpclient"
	"time"
)

func main() {

	config.LoadConfig()

	database.InitDb()

	httpClient := httpclient.NewHTTPClient(5, 3, 1*time.Second)
	defer httpClient.Close()

	hybridCache := cache.NewHybridCache()
	if hybridCache != nil {

	}

	mq := messaging.NewMessageQueue()
	mq.StartMessageQueue()
	defer mq.Close()

	api.StartApiServer()
}
