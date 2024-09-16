package domain

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"

	"github.com/google/uuid"
)

// Password is a type for user password.
// It is formatted as "<algorithm>$<iterations>$<salt>$<hash>".
type Password string

func (p Password) String() string {
	return string(p)
}

const (
	keyLen           = 32
	hashIterations   = 120000
	encryptAlgorithm = "pbkdf2_sha256"
	saltCharset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	saltLen          = 12
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewPassword(plain string) Password {
	salt := make([]byte, saltLen)
	for i := range salt {
		salt[i] = saltCharset[seededRand.Intn(len(saltCharset))]
	}
	encrypted := base64.StdEncoding.EncodeToString(
		pbkdf2.Key([]byte(plain), salt, hashIterations, keyLen, sha256.New),
	)
	return Password(fmt.Sprintf("%s$%d$%s$%s", encryptAlgorithm, hashIterations, salt, encrypted))
}

func (p Password) Compare(plain string) bool {
	pwd := strings.Split(string(p), "$")
	if len(pwd) != 4 {
		slog.Error("invalid password format")
		return false
	}

	encrypted := base64.StdEncoding.EncodeToString(
		pbkdf2.Key([]byte(plain), []byte(pwd[2]), hashIterations, keyLen, sha256.New),
	)
	return pwd[3] == encrypted
}

type User struct {
	ID uuid.UUID

	Username  string
	Password  Password
	CreatedAt time.Time
	UpdatedAt time.Time
}
