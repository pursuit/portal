package rest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pursuit/portal/internal/rest"
	"github.com/pursuit/portal/internal/service/user/mock"

	"github.com/golang/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	for _, testcase := range []struct {
		tName string

		body     string
		username string
		password []byte

		validInput bool

		httpCode int
		output   string
	}{
		{
			tName:    "no body",
			httpCode: 422,
			output:   "invalid input",
		},
		{
			tName:    "invalid username is not a string",
			body:     `{"username":77, "password":"password123"}`,
			httpCode: 422,
			output:   "invalid input",
		},
		{
			tName:    "invalid password",
			body:     `{"username":"Bambang", "password":""}`,
			username: "Bambang",
			password: []byte("password123"),
			httpCode: 422,
			output:   "invalid input",
		},
		{
			tName:      "success",
			body:       `{"username":"Bambang", "password":"password123"}`,
			validInput: true,
			username:   "Bambang",
			password:   []byte("password123"),
			httpCode:   201,
		},
		{
			tName:      "success password is not a string",
			body:       `{"username":"Bambang", "password":777}`,
			validInput: true,
			username:   "Bambang",
			password:   []byte("7"),
			httpCode:   201,
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			mocker := gomock.NewController(t)
			defer mocker.Finish()

			svc := mock_user.NewMockService(mocker)
			if testcase.validInput {
				svc.EXPECT().Create(gomock.Any(), testcase.username, testcase.password).Return(nil)
			}

			req, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(testcase.body))
			resp := httptest.NewRecorder()

			h := rest.Handler{svc}
			h.CreateUser(resp, req)

			if testcase.httpCode != resp.Code {
				t.Errorf("Test %s httpcode is %d, should be %d", testcase.tName, resp.Code, testcase.httpCode)
			}

			if testcase.output != resp.Body.String() {
				t.Errorf("Test %s body is %s, should be %s", testcase.tName, resp.Body.String(), testcase.output)
			}
		})
	}
}
