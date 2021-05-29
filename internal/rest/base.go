package rest

import (
	"github.com/pursuit/portal/internal/service/user"
)

type Handler struct {
	UserService user.Service
}
