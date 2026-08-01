package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID     int64  `json:"uid"`
	EmployeeNo string `json:"emp"`
	Role       string `json:"role"`

	jwt.RegisteredClaims
}

type UserIdentity struct {
	UserID     int64
	EmployeeNo string
	Role       string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
}
