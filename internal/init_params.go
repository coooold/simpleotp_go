package internal

import (
	"flag"
	"os"
	"strconv"
)

var (
	Secret       string
	CookieName   string
	Port		 int
)

func InitParams() {
	flag.StringVar(&Secret, "secret", "", "TOTP secret key")
	flag.StringVar(&CookieName, "cookie", "", "cookie name")
	flag.IntVar(&Port, "port", 0, "listen port default 8000")

	flag.Parse()

	if Secret == "" {
		Secret = os.Getenv("SECRET_KEY")
	}

	if CookieName == "" {
		CookieName = os.Getenv("COOKIE_NAME")
	}

	if Port == 0 {
		Port, _ = strconv.Atoi(os.Getenv("PORT"))
	}

	/// override with default value

	if Secret == "" {
		Secret = "simple_otp_go"
	}

	if CookieName == "" {
		CookieName = "simple_otp_go"
	}

	if Port == 0 {
		Port = 8000
	}
}