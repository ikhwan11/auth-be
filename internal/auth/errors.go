package auth

import "errors"

var (
	ErrInvalidCredential = errors.New("invalid credential")

	ErrPasswordMismatch = errors.New("password mismatch")

	ErrEmployeeNotFound = errors.New("employee not found")

	ErrUserAlreadyExists = errors.New("user already exists")

	ErrRefreshTokenExpired = errors.New("refresh token expired")

	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	ErrLogoutFailed = errors.New("failed to logout")
)
