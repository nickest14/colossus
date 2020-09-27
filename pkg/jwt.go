package pkg

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(envy.Get("JWT_KEY", ""))
}

// Claims defines the custom jwt standard claims
type Claims struct {
	Email    string `json:"email"`
	Provider string `json:"prodvider"`
	Group    string `json:"group"`
	jwt.StandardClaims
}

// CreateJWTToken function is used to create user JWT token
func CreateJWTToken(email string, provider string, providerID string, group string) (string, error) {
	claims := Claims{
		Email:    email,
		Provider: provider,
		Group:    group,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Id:        providerID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// CheckJWTToken function is used to check whether the token is correct or not
func CheckJWTToken(auth string) (*Claims, error) {
	if !strings.HasPrefix(auth, "JWT ") {
		return &Claims{}, errors.New("tokenstring should contains 'JWT'")
	}
	token := strings.Split(auth, "JWT ")[1]
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtSecret, nil
	})
	if err != nil {
		message := ""
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				message = "token is malformed"
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				message = "token could not be verified because of signing problems"
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				message = "signature validation failed"
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				message = "token is expired"
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				message = "token is not yet valid before sometime"
			} else {
				message = "can not handle this token"
			}
		}
		return &Claims{}, errors.New(message)
	}
	if claim, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claim, nil
	}
	return &Claims{}, errors.New("token is not valid")
}
