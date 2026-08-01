package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID     uuid.UUID `json:"uid"`
	EmployeeNo string    `json:"emp"`
	Role       string    `json:"role"`

	jwt.RegisteredClaims
}

type UserIdentity struct {
	UserID     uuid.UUID
	EmployeeNo string
	Role       string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
}
