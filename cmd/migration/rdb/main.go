package main

import (
	"github.com/buzzryan/zenbu/internal/config"
	"github.com/buzzryan/zenbu/internal/rdbutil"
	"github.com/buzzryan/zenbu/internal/user/infra"
)

func main() {
	cfg := config.LoadConfigFromEnv()
	db := rdbutil.MustConnectMySQL(cfg.MySQLConfig)

	err := db.AutoMigrate(
		&infra.User{},
	)
	if err != nil {
		panic(err)
	}
}
