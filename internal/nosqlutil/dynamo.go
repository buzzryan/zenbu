package nosqlutil

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/guregu/dynamo/v2"
)

func MustConnectDDB() *dynamo.DB {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Panicf("unable to load aws sdk config, %v", err)
	}

	return dynamo.New(cfg)
}
