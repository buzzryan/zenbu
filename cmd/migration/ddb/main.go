package main

import (
	"context"
	"log"

	awscfg "github.com/aws/aws-sdk-go-v2/config"

	"github.com/buzzryan/zenbu/internal/config"
	"github.com/buzzryan/zenbu/internal/nosqlutil"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfigFromEnv()
	awsCfg, err := awscfg.LoadDefaultConfig(ctx)
	if err != nil {
		log.Panicf("failed to load AWS config: %v", err)
	}
	ddb := nosqlutil.ConnectDDB(awsCfg, cfg.DynamoConfig)

	err = ddb.CreateTable(cfg.TableName, &nosqlutil.CommonSchema{}).Run(ctx)
	if err != nil {
		log.Panicf("failed to create table: %v", err)
	}

	tables, err := ddb.ListTables().All(ctx)
	if err != nil {
		log.Panicf("failed to list tables: %v", err)
	}
	log.Printf("tables: %v\n", tables)
}
