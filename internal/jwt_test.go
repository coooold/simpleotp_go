package internal

import (
	"testing"
	"time"
)

func Test_JWT(t *testing.T) {
	jm := NewJWTManager([]byte("a test secret key"), 1)
	token := jm.GenerateToken("0")
	_, err := jm.Valid(token, "1")

	if err.Error() != "uid not match" {
		t.Errorf("must valid uid")
	}

	time.Sleep(time.Second * 2)

	_, err2 := jm.Valid(token, "0")

	if err2.Error() != "Token is expired" {
		t.Errorf("must valid maxAge")
	}

	// generate a new token
	token = jm.GenerateToken("0")
	_, err3 := jm.Valid(token, "0")

	if err3 != nil {
		t.Errorf("must pass validation")
	}
}
