package jwt

import (
	"context"
	commonError "enuma-elish/pkg/error"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"time"
)

const ContextKey = "JWT"

func GenerateToken(duration time.Duration, payload map[string]interface{}, secret string) (string, error) {
	expirationTime := time.Now().Add(duration).Unix()
	claims := jwt.MapClaims{}
	claims["exp"] = expirationTime
	for key, v := range payload {
		claims[key] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	return tokenString, nil
}

func Verify(token, secret string) (*jwt.Token, error) {
	tokenValidation, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !tokenValidation.Valid {
		return nil, commonError.New("token is invalid", http.StatusUnauthorized)
	}

	if claims, ok := tokenValidation.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
			return nil, commonError.New("token is expired", http.StatusUnauthorized)
		}
	}

	return tokenValidation, nil
}

func ExtractToken(token jwt.Token) (jwt.MapClaims, error) {
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}
	return claim, nil
}

func ExtractContext(c context.Context) (jwt.MapClaims, error) {
	claims, ok := c.Value(ContextKey).(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to extract claims from context")
	}
	return claims, nil
}
