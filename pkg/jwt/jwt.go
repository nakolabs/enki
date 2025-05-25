package jwt

import (
	"context"
	commonError "enuma-elish/pkg/error"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const ContextKey = "JWT"

type Payload struct {
	Exp  int64  `json:"exp"`
	Iat  int64  `json:"iat"`
	Sub  string `json:"sub"`
	Iss  string `json:"iss"`
	Nbf  int64  `json:"nbf"`
	Aud  string `json:"aud"`
	User User   `json:"user"`
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	SchoolID uuid.UUID `json:"school_id"`
	RoleID   uuid.UUID `json:"role_id"`
}

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(p.Exp, 0)}, nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(p.Iat, 0)}, nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(p.Nbf, 0)}, nil
}

func (p *Payload) GetIssuer() (string, error) {
	return p.Iss, nil
}

func (p *Payload) GetSubject() (string, error) {
	return p.Sub, nil
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{p.Aud}, nil
}

func GenerateToken(payload Payload, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Verify(token, secret string) (*jwt.Token, error) {
	payload := new(Payload)
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), jwt.WithIssuedAt(), jwt.WithExpirationRequired())
	tokenValidation, err := parser.ParseWithClaims(token, payload, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !tokenValidation.Valid {
		return nil, commonError.New("token is invalid", http.StatusUnauthorized)
	}

	return tokenValidation, nil
}

func ExtractToken(token *jwt.Token) (*Payload, error) {
	claims, ok := token.Claims.(*Payload)
	if !ok {
		return nil, commonError.New("failed to extract token claims", http.StatusUnauthorized)
	}
	return claims, nil
}

func ExtractContext(c context.Context) (*Payload, error) {
	claims, ok := c.Value(ContextKey).(*Payload)
	if !ok {
		return nil, commonError.New("failed to extract claims from context", http.StatusUnauthorized)
	}
	return claims, nil
}
