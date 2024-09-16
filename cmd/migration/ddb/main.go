package main

import (
	"context"
	"log"

	"github.com/buzzryan/zenbu/internal/config"
	"github.com/buzzryan/zenbu/internal/nosqlutil"
)

func main() {
	cfg := config.LoadConfigFromEnv()
	ddb := nosqlutil.MustConnectDDB(cfg.DynamoConfig)

	ctx := context.Background()

	err := ddb.CreateTable(cfg.TableName, &nosqlutil.CommonSchema{}).Run(ctx)
	if err != nil {
		log.Panicf("failed to create table: %v", err)
	}

	tables, err := ddb.ListTables().All(ctx)
	if err != nil {
		log.Panicf("failed to list tables: %v", err)
	}
	log.Printf("tables: %v\n", tables)
}
