package gormtypes

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type UUID [16]byte

func (u *UUID) String() string {
	return fmt.Sprintf("%x", *u)
}

func (u *UUID) GormDataType() string {
	return "binary(16)"
}

// Scan tells GORM how to receive from the database
func (u *UUID) Scan(value interface{}) error {
	bytes, err := value.([]uint8)
	if !err {
		return errors.New("failed to convert to [16]byte")
	}
	if len(bytes) != 16 {
		return errors.New("failed to convert to [16]byte")
	}
	*u = [16]byte(bytes)
	return nil
}

// Value tells GORM how to create into the database
func (u UUID) Value() (driver.Value, error) {
	return u[:], nil
}
