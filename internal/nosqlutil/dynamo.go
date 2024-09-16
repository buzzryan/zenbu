package nosqlutil

import (
	"context"
	"log"
	"log/slog"
	"net/url"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	smithy "github.com/aws/smithy-go/endpoints"
	"github.com/guregu/dynamo/v2"

	"github.com/buzzryan/zenbu/internal/config"
)

type endpointResolver struct {
	endpoint string
}

func (e *endpointResolver) ResolveEndpoint(ctx context.Context, _ dynamodb.EndpointParameters) (smithy.Endpoint, error) {
	u, err := url.Parse(e.endpoint)
	if err != nil {
		slog.ErrorContext(ctx, "connect dynamodb: failed to parse endpoint", slog.Any("err", err))
		return smithy.Endpoint{}, err
	}
	return smithy.Endpoint{URI: *u}, nil
}

func MustConnectDDB(ddbConfig config.DynamoConfig) *dynamo.DB {
	cfg, err := awscfg.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Panicf("unable to load aws sdk config, %v", err)
	}

	if ddbConfig.Endpoint == "" {
		return dynamo.New(cfg)
	}

	return dynamo.New(cfg, dynamodb.WithEndpointResolverV2(&endpointResolver{endpoint: ddbConfig.Endpoint}))
}
