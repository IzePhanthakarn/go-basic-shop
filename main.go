package main

import (
	"os"

	"github.com/IzePhanthakarn/go-basic-shop/config"
	"github.com/IzePhanthakarn/go-basic-shop/modules/servers"
	"github.com/IzePhanthakarn/go-basic-shop/pkg/databases"

	_ "github.com/IzePhanthakarn/go-basic-shop/docs"
)

// @title Swagger Basic Shop API 1.0
// @version 1.0.0
// @description This is a sample swagger for Basic Shop
// @host localhost:3000
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token.

func envPath() string {
	if len(os.Args) == 1 {
		return ".env.dev"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())
	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.NewServer(cfg, db).Start()
}
