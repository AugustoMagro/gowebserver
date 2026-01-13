package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("Malformed Authorization header")
	}
	return splitAuth[1], nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	expired := time.Now().UTC().Add(expiresIn)

	signedKey := []byte(tokenSecret)

	claims := &jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		ExpiresAt: jwt.NewNumericDate(expired),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signedKey)

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return uuid.Nil, err
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid user")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 256)

	refreshToken, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	src := make([]byte, refreshToken)

	return hex.EncodeToString(src), nil
}
