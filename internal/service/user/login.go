package user

import (
	"context"
	"errors"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/model"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

func (this Svc) Login(ctx context.Context, username string, password []byte) (string, *internal.E) {
	defer func() {
		for i := 0; i < len(password); i++ {
			password[i] = 0
		}
	}()

	user, err := this.UserRepo.GetByUsername(ctx, this.DB, username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, password); err != nil {
		return "", &internal.E{
			Err:    errors.New("invalid password"),
			Status: internal.EInvalidPassword,
		}
	}

	jwtBody := model.Jwt{
		ID: user.ID,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtBody)
	token, _ := at.SignedString([]byte("zxcwqe"))

	return token, nil
}
