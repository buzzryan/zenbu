package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/dynamo/v2"

	nosqlutil2 "github.com/buzzryan/zenbu/internal/commonutil/nosqlutil"
	"github.com/buzzryan/zenbu/internal/user/domain"
	"github.com/buzzryan/zenbu/internal/user/usecase"
)

const (
	userPartitionKeyPrefix = "USER"
	usernamePartitionKey   = "USERNAME"
	userProfileSortKey     = "PROFILE"
)

// dynamoUserRepo is the implementation of usecase.UserRepo interface using AWS DynamoDB. (adapter)
type dynamoUserRepo struct {
	ddb       *dynamo.DB
	tableName string
}

func NewDynamoUserRepo(ddb *dynamo.DB, tableName string) usecase.UserRepo {
	return &dynamoUserRepo{ddb: ddb, tableName: tableName}
}

type UserProfile struct {
	nosqlutil2.CommonSchema

	Username  string    `dynamo:"un"`
	Password  string    `dynamo:"pw"`
	CreatedAt time.Time `dynamo:"ca"`
	UpdatedAt time.Time `dynamo:"ua"`

	Bio string `dynamo:"bio"`
}

func (un *UserProfile) toDomainEntity() *domain.User {
	return &domain.User{
		ID:        uuid.MustParse(un.PartitionKey[len(userPartitionKeyPrefix)+1:]),
		Username:  un.Username,
		Password:  domain.Password(un.Password),
		CreatedAt: un.CreatedAt,
		UpdatedAt: un.UpdatedAt,
		Bio:       un.Bio,
	}
}

func buildUserProfile(u *domain.User) *UserProfile {
	return &UserProfile{
		CommonSchema: nosqlutil2.CommonSchema{
			PartitionKey: userPartitionKeyPrefix + "#" + u.ID.String(),
			SortKey:      userProfileSortKey,
		},
		Username:  u.Username,
		Password:  u.Password.String(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type Username struct {
	nosqlutil2.CommonSchema
	UserID string `dynamo:"uid"`
}

func buildUsername(u *domain.User) *Username {
	return &Username{
		CommonSchema: nosqlutil2.CommonSchema{
			PartitionKey: usernamePartitionKey,
			SortKey:      u.Username,
		},
		UserID: u.ID.String(),
	}
}

func (dur *dynamoUserRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	createUsername := dur.ddb.Table(dur.tableName).
		Put(buildUsername(u)).IncludeItemInCondCheckFail(true).If("attribute_not_exists(pk)")
	createUserProfile := dur.ddb.Table(dur.tableName).
		Put(buildUserProfile(u)).If("attribute_not_exists(pk)")

	err := dur.ddb.WriteTx().Put(createUsername).Put(createUserProfile).Run(ctx)

	if nosqlutil2.IsConditionalCheckFailed(err) {
		return nil, usecase.ErrUsernameAlreadyExists
	}
	if err != nil {
		return nil, fmt.Errorf("dynamoUserRepo.Create failed: %w", err)
	}

	return u, nil
}

func (dur *dynamoUserRepo) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	userProfile := &UserProfile{}
	err := dur.ddb.Table(dur.tableName).
		Get("pk", userPartitionKeyPrefix+"#"+id.String()).
		Range("sk", dynamo.Equal, userProfileSortKey).
		One(ctx, &userProfile)

	if errors.Is(err, dynamo.ErrNotFound) {
		return nil, usecase.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("dynamoUserRepo.Get failed: %w", err)
	}

	return userProfile.toDomainEntity(), nil
}

func (dur *dynamoUserRepo) GetByName(ctx context.Context, username string) (*domain.User, error) {
	un := &Username{}
	err := dur.ddb.Table(dur.tableName).
		Get("pk", usernamePartitionKey).
		Range("sk", dynamo.Equal, username).One(ctx, &un)
	if errors.Is(err, dynamo.ErrNotFound) {
		return nil, usecase.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("dynamoUserRepo.GetByName failed: %w", err)
	}

	return dur.Get(ctx, uuid.MustParse(un.UserID))
}

func (dur *dynamoUserRepo) Update(ctx context.Context, u *domain.User) (*domain.User, error) {
	var res UserProfile
	err := dur.ddb.Table(dur.tableName).
		Update("pk", userPartitionKeyPrefix+"#"+u.ID.String()).
		Range("sk", userProfileSortKey).
		Set("bio", u.Bio).
		Set("ua", u.UpdatedAt).
		Value(ctx, &res) // ALL_NEW - Returns all the attributes of the item, as they appear after the UpdateItem operation.

	if errors.Is(err, dynamo.ErrNotFound) {
		return nil, usecase.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return res.toDomainEntity(), nil
}
