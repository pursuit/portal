package rest_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestUserLoginInvalid(t *testing.T) {
	resp, err := http.Post("http://localhost:5002/sessions", "application/json", strings.NewReader(`{"username":"t","password":"thisisaord"}`))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 422 {
		t.Errorf("Test login code is %d, should be 422", resp.StatusCode)
	}

	var buf strings.Builder
	io.Copy(&buf, resp.Body)
	output := buf.String()
	if output != "" {
		t.Errorf("Test login output si %s, should be empty", output)
	}
}

func TestUser(t *testing.T) {
	for _, testcase := range []struct {
		tName string

		body io.Reader

		httpCode int
		output   string
	}{
		{
			tName:    "nil body",
			httpCode: 422,
			output:   "invalid input",
		},
		{
			tName:    "username size is too small",
			body:     strings.NewReader(`{"username":"t","password":"thisisaord"}`),
			httpCode: 422,
			output:   "invalid input",
		},
		{
			tName:    "username size is too big",
			body:     strings.NewReader(`{"username":"tqweqweqweqwdsacasqwadsaqwe","password":"thisisaord"}`),
			httpCode: 422,
			output:   "invalid input",
		},
		{
			tName:    "no username",
			body:     strings.NewReader(`{"password":"thisisaord"}`),
			httpCode: 422,
			output:   "invalid input",
		},
		{
			tName:    "password size is too small",
			body:     strings.NewReader(`{"username":"test1234","password":""}`),
			httpCode: 422,
			output:   "invalid input",
		},
		{
			tName:    "password size is too big",
			body:     strings.NewReader(`{"username":"test1234","password":"qweqwedascassdeqweadszcawqe"}`),
			httpCode: 422,
			output:   "invalid input",
		},
		{
			tName:    "username contain special char",
			body:     strings.NewReader(`{"username":"tes?t1234","password":"thisisaord"}`),
			httpCode: 422,
			output:   "invalid input",
		},
		{
			tName:    "password contain special char",
			body:     strings.NewReader(`{"username":"test1234","password":"thisi?saord"}`),
			httpCode: 422,
			output:   "invalid input",
		},
		{
			tName:    "valid",
			body:     strings.NewReader(`{"username":"test1234","password":"thisisaord"}`),
			httpCode: 201,
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			resp, err := http.Post("http://localhost:5002/users", "application/json", testcase.body)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			var buf strings.Builder
			io.Copy(&buf, resp.Body)
			output := buf.String()
			if testcase.output != output {
				t.Errorf("Test %s output si %s, should be %s", testcase.tName, output, testcase.output)
			}

			if testcase.httpCode != resp.StatusCode {
				t.Errorf("Test %s code is %d, should be %d", testcase.tName, resp.StatusCode, testcase.httpCode)
			}
		})
	}

	respL, err := http.Post("http://localhost:5002/sessions", "application/json", strings.NewReader(`{"username":"test1234","password":"thisisaord"}`))
	if err != nil {
		panic(err)
	}
	defer respL.Body.Close()
	if respL.StatusCode != 201 {
		t.Errorf("Test login after created code is %d, should be 201", respL.StatusCode)
	}

	var bufL strings.Builder
	io.Copy(&bufL, respL.Body)
	outputL := bufL.String()
	if outputL == "" {
		t.Error("Test login after created code output should not be empty")
	}
}
