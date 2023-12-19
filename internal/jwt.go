package internal


import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type JWTManager struct {
	maxAge int
	secretKey []byte
}

func NewJWTManager(secretKey []byte, maxAge int) *JWTManager {
	return &JWTManager{
		maxAge:    maxAge,
		secretKey: secretKey,
	}
}


func (jm *JWTManager) Valid(tokenString string, uid string) (bool, error) {
	claims, err := jm.parseToken(tokenString)
	if err != nil {
		return false, err
	}

	err = claims.Valid()
	if err != nil {
		return false, err
	}

	if uid != claims["uid"] {
		return false, fmt.Errorf("uid not match")
	}

	return true, nil
}


func (jm *JWTManager) GenerateToken(uid string) string {
	claims := jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(time.Duration(jm.maxAge)*time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jm.secretKey)
	if err!=nil {
		log.Print(err)
	}

	return tokenString
}

func (jm *JWTManager) parseToken(tokenString string)(jwt.MapClaims, error)  {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jm.secretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims,nil
	} else {
		return nil,err
	}
}

