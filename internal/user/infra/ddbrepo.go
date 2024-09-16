package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/dynamo/v2"

	"github.com/buzzryan/zenbu/internal/nosqlutil"
	"github.com/buzzryan/zenbu/internal/user/domain"
)

const (
	userPartitionKeyPrefix = "USER"
	usernamePartitionKey   = "USERNAME"
	userProfileSortKey     = "PROFILE"
)

type dynamoUserRepo struct {
	ddb       *dynamo.DB
	tableName string
}

func NewDynamoUserRepo(ddb *dynamo.DB, tableName string) domain.UserRepo {
	return &dynamoUserRepo{ddb: ddb, tableName: tableName}
}

type UserProfile struct {
	nosqlutil.CommonSchema

	Username  string    `dynamo:"un"`
	Password  string    `dynamo:"pw"`
	CreatedAt time.Time `dynamo:"ca"`
	UpdatedAt time.Time `dynamo:"ua"`
}

func (un *UserProfile) toDomainEntity() *domain.User {
	return &domain.User{
		ID:        uuid.MustParse(un.PartitionKey[len(userPartitionKeyPrefix)+1:]),
		Username:  un.Username,
		Password:  domain.Password(un.Password),
		CreatedAt: un.CreatedAt,
		UpdatedAt: un.UpdatedAt,
	}
}

func buildUserProfile(u *domain.User) *UserProfile {
	return &UserProfile{
		CommonSchema: nosqlutil.CommonSchema{
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
	nosqlutil.CommonSchema
	UserID string `dynamo:"uid"`
}

func buildUsername(u *domain.User) *Username {
	return &Username{
		CommonSchema: nosqlutil.CommonSchema{
			PartitionKey: usernamePartitionKey,
			SortKey:      u.Username,
		},
		UserID: u.ID.String(),
	}
}

func (ur *dynamoUserRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	createUsername := ur.ddb.Table(ur.tableName).
		Put(buildUsername(u)).IncludeItemInCondCheckFail(true).If("attribute_not_exists(pk)")
	createUserProfile := ur.ddb.Table(ur.tableName).
		Put(buildUserProfile(u)).If("attribute_not_exists(pk)")

	err := ur.ddb.WriteTx().Put(createUsername).Put(createUserProfile).Run(ctx)

	if nosqlutil.IsConditionalCheckFailed(err) {
		return nil, domain.ErrUsernameAlreadyExists
	}
	if err != nil {
		return nil, fmt.Errorf("dynamoUserRepo.Create failed: %w", err)
	}

	return u, nil
}

func (ur *dynamoUserRepo) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	userProfile := &UserProfile{}
	err := ur.ddb.Table(ur.tableName).
		Get("pk", userPartitionKeyPrefix+"#"+id.String()).
		Range("sk", dynamo.Equal, userProfileSortKey).
		One(ctx, &userProfile)

	if errors.Is(err, dynamo.ErrNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("dynamoUserRepo.Get failed: %w", err)
	}

	return userProfile.toDomainEntity(), nil
}
