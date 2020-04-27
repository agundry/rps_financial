package main

import (
	"github.com/agundry/rps_financial/app"
	"github.com/agundry/rps_financial/config"
	log "github.com/agundry/rps_financial/util"
)

func main() {
	cfg := config.Configure()

	var standardLogger = log.NewLogger()

	application := app.App{}
	application.WithLogging(standardLogger)
	application.Initialize(
		cfg.Db.Username,
		cfg.Db.Password,
		cfg.Db.Server,
		cfg.Db.DbName)

	application.Run(":" + cfg.Server.Port)
	standardLogger.Infof("Application started: Listening on port %s", cfg.Server.Port)
}
