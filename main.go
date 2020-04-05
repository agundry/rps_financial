package main

import (
	"github.com/agundry/rps_financial/app"
	"github.com/agundry/rps_financial/config"
)

func main() {
	cfg := config.Configure()

	application := app.App{}
	application.Initialize(
		cfg.Db.Username,
		cfg.Db.Password,
		cfg.Db.Server,
		cfg.Db.DbName)

	application.Run(":" + cfg.Server.Port)
}
