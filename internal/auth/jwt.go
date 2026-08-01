package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type JWT struct {
	secretKey      []byte
	accessExpired  time.Duration
	refreshExpired time.Duration
	issuer         string
}

func NewJWT(
	secret string,
	accessExpired string,
	refreshExpired string,
	issuer string,
) (*JWT, error) {
	accessDuration, err := time.ParseDuration(accessExpired)
	if err != nil {
		return nil, err
	}

	refreshDuration, err := time.ParseDuration(refreshExpired)
	if err != nil {
		return nil, err
	}

	return &JWT{
		secretKey:      []byte(secret),
		accessExpired:  accessDuration,
		refreshExpired: refreshDuration,
		issuer:         issuer,
	}, nil
}

func (j *JWT) GenerateAccessToken(user UserIdentity) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID:     user.UserID,
		EmployeeNo: user.EmployeeNo,
		Role:       user.Role,

		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.UserID, 10),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessExpired)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secretKey)
}

func (j *JWT) GenerateRefreshToken(user UserIdentity) (string, error) {
	now := time.Now()

	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.UserID, 10),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshExpired)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secretKey)
}

func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma signing adalah HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}

			return j.secretKey, nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
