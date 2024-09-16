package nosqlutil

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// go aws/dynamodb doesn't support predefined error variables. And predefined errors for guregu/dynamo isn't enough.
// So we need to define our own errors to wrap the errors from dynamodb.

var (
	ErrConditionalCheckFailed = errors.New("conditional check failed")
)

const (
	CodeConditionalCheckFailed = "ConditionalCheckFailed"
)

func WrapError(err error) error {
	var dynamoErr *types.TransactionCanceledException
	if errors.As(err, &dynamoErr) {
		for _, reason := range dynamoErr.CancellationReasons {
			if reason.Code != nil && *reason.Code == CodeConditionalCheckFailed {
				return errors.Join(ErrConditionalCheckFailed, err)
			}
		}
	}
	return err
}

func IsConditionalCheckFailed(err error) bool {
	return errors.Is(WrapError(err), ErrConditionalCheckFailed)
}
