package user

import (
	"context"
	"errors"
	"time"
	"unicode"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/proto/out/event"
	"github.com/pursuit/portal/internal/repo"

	"github.com/pursuit/event-go/pkg"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"golang.org/x/crypto/bcrypt"
)

func (this Svc) Create(ctx context.Context, username string, password []byte) (int, *internal.E) {
	defer func() {
		for i := 0; i < len(password); i++ {
			password[i] = 0
		}
	}()

	if !validInputCreate(username, password) {
		return 0, &internal.E{errors.New("invalid input"), internal.EInvalidInput}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(password, 12)
	if err != nil {
		return 0, &internal.E{err, internal.EFailedHashPassword}
	}

	return this.process(ctx, username, hashedPassword)
}

func (this Svc) process(ctx context.Context, username string, hashedPassword []byte) (int, *internal.E) {
	var userID int
	err := repo.Transaction(ctx, this.DB, func(db repo.DB) *internal.E {
		now := time.Now().UTC()
		id, err := this.UserRepo.Create(ctx, db, username, hashedPassword, now)
		if err != nil {
			return err
		}

		createdAtProto, _ := ptypes.TimestampProto(now)
		payload := event_proto.UserCreated{
			Id:        uint64(id),
			Username:  username,
			CreatedAt: createdAtProto,
		}

		protodata, _ := proto.Marshal(&payload)

		if err := pkg.StoreEvent(ctx, db, pkg.EventData{
			Topic:   "portal.user.created.x2",
			Payload: protodata,
		}); err != nil {
			return &internal.E{err, internal.EFailedStoreEvent}
		}

		userID = id
		return nil
	})

	return userID, err
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
