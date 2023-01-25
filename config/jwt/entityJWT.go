package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

var JWT_KEY = []byte("rahasia")

type JWTclaim struct {
	ID      int64
	UserID  int64
	StoreID int64
	Email   string
	Name    string
	jwt.StandardClaims
}
