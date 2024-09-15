package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/buzzryan/zenbu/internal/rdbutil/gormtypes"
	"github.com/buzzryan/zenbu/internal/user/domain"
)

type userRepo struct {
	db *gorm.DB
}

// User is a struct that represents a user DB schema.
type User struct {
	ID       gormtypes.UUID `gorm:"type:binary(16);primaryKey;not null"`
	Username string         `gorm:"type:varchar(32) CHARACTER SET ascii COLLATE ascii_bin;not null;uniqueIndex"`
	Password string         `gorm:"type:varchar(100) CHARACTER SET ascii COLLATE ascii_bin;not null"`
}

func (u *User) TableName() string {
	return "user"
}

func NewUserRepo(db *gorm.DB) domain.UserRepo {
	return &userRepo{db: db}
}

func (u *User) toDomainEntity() *domain.User {
	return &domain.User{
		ID:       uuid.UUID(u.ID),
		Username: u.Username,
		Password: domain.Password(u.Password),
	}
}

func fromDomainEntity(e *domain.User) *User {
	return &User{
		ID:       gormtypes.UUID(e.ID),
		Username: e.Username,
		Password: string(e.Password),
	}
}

func (ur *userRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	record := fromDomainEntity(u)
	res := ur.db.WithContext(ctx).Create(&record)
	if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
		return nil, domain.ErrUsernameAlreadyExists
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to create user: %w", res.Error)
	}
	return record.toDomainEntity(), nil
}
