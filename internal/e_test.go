package internal_test

import (
	"errors"
	"testing"

	"github.com/pursuit/portal/internal"
)

func TestE(t *testing.T) {
	err := "this is a random error"
	e := internal.E{Err: errors.New(err)}
	if err != e.Error() {
		t.Fatalf("error is %s, should be %s", e.Error(), err)
	}
}
