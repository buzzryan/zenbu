package nosqlutil

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
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

func ConnectDDB(awsCfg aws.Config, ddbConfig config.DynamoConfig) *dynamo.DB {
	if ddbConfig.Endpoint == "" {
		return dynamo.New(awsCfg)
	}

	return dynamo.New(awsCfg, dynamodb.WithEndpointResolverV2(&endpointResolver{endpoint: ddbConfig.Endpoint}))
}

// CommonSchema is a common schema for DynamoDB.
// Zenbu uses a single-table design.
// https://aws.amazon.com/blogs/compute/creating-a-single-table-design-with-amazon-dynamodb/
type CommonSchema struct {
	PartitionKey string `dynamo:"pk,hash"`
	SortKey      string `dynamo:"sk,range"`
}
