package rest

import (
	"encoding/json"
	"net/http"
)

type createUserPayload struct {
	Username string          `json:"username"`
	Password json.RawMessage `json:"password"`
}

func (this Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload createUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(422)
		w.Write([]byte("invalid input"))
		return
	}

	if len(payload.Password) < 3 {
		w.WriteHeader(422)
		w.Write([]byte("invalid input"))
		return
	}

	if err := this.UserService.Create(r.Context(), payload.Username, payload.Password[1:len(payload.Password)-1]); err != nil {
		httpStatus := err.Status / 1_000
		w.WriteHeader(httpStatus)
		if httpStatus >= 400 && httpStatus < 500 {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("please try again in a few moment"))
		return
	}

	w.WriteHeader(201)
}
