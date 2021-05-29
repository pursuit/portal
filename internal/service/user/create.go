package user

import (
	"context"
	"errors"
	"unicode"

	"github.com/pursuit/portal/internal"

	"golang.org/x/crypto/bcrypt"
)

func (this Svc) Create(ctx context.Context, username string, password []byte) *internal.E {
	defer func() {
		for i := 0; i < len(password); i++ {
			password[i] = 0
		}
	}()

	if !validInputCreate(username, password) {
		return &internal.E{errors.New("invalid input"), 422}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(password, 12)
	if err != nil {
		return &internal.E{err, 503}
	}

	_, err = this.UserRepo.Create(ctx, this.DB, username, hashedPassword)
	if err != nil {
		return &internal.E{err, 503}
	}

	return nil
}

func validInputCreate(username string, password []byte) bool {
	if len(username) < 6 || len(username) > 12 {
		return false
	}

	if len(password) < 6 || len(password) > 12 {
		return false
	}

	for _, ch := range username {
		if !(unicode.IsLetter(ch) || unicode.IsDigit(ch)) {
			return false
		}
	}

	for _, ch := range password {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
			return false
		}
	}

	return true
}
