package internal

import (
	"fmt"
	"github.com/pquerna/otp/totp"
	"log"
	"net/http"
	"net/url"
	"time"
)

const loginFormTpl = `
<html>
<head>
<title>Please Log In</title>
</head>
<body>
<form action="/totp/login" method="POST">
<input type="text" name="code">
<input type="hidden" name="original_uri" value="%s">
<input type="submit" value="Submit">
</form>
</body>
</html>`

var lastLoginAttemptTime = time.Now()
var jwtManager = NewJWTManager([]byte(Secret), maxAge)

const loginStopInterval = 2 * time.Second // stop time before next login attemt
const maxAge = 3600

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		originalURI := r.Header.Get("X-Original-URI")

		w.Header().Set("Content-Type", "text/html")
		html := fmt.Sprintf(loginFormTpl, url.QueryEscape(originalURI))
		_, _ = w.Write([]byte(html))
		return
	}

	if r.Method == http.MethodPost {
		if time.Since(lastLoginAttemptTime) < loginStopInterval {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte("Slow down. Hold your horses"))
			return
		}

		_ = r.ParseForm()
		code := r.PostForm.Get("code")
		originalURI := r.PostForm.Get("original_uri")
		if originalURI == "" {
			originalURI ="/"
		}

		if totp.Validate(code, Secret) {
			jwtToken := jwtManager.GenerateToken("0")

			http.SetCookie(w, &http.Cookie{
				Name:       CookieName,
				Value:      jwtToken,
				Path:       "/",
				Domain:     "",
				Expires:    time.Now().Add(time.Second * maxAge),
				RawExpires: "",
				MaxAge:     0,
				Secure:     false,
				HttpOnly:   false,
				SameSite:   0,
				Raw:        "",
				Unparsed:   nil,
			})

			log.Print("totp valid passed")


			http.Redirect(w, r, originalURI, http.StatusFound)
			return
		}

		log.Print("totp valid failed")

		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	jwtCookie, err := r.Cookie(CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwtToken := jwtCookie.Value
	if ok, _ := jwtManager.Valid(jwtToken, "0"); ok {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}