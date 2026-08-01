package auth

import (
	"context"
	"errors"
	"time"

	"github.com/ikhwan11/auth-be/internal/employee"
	refreshtoken "github.com/ikhwan11/auth-be/internal/refresh_token"
	"github.com/ikhwan11/auth-be/internal/user"
)

type Service struct {
	employeeRepo employee.Repository
	userRepo     user.Repository
	tokenRepo    refreshtoken.Repository
	jwt          *JWT
}

func NewService(
	employeeRepo employee.Repository,
	userRepo user.Repository,
	tokenRepo refreshtoken.Repository,
	jwt *JWT,
) *Service {
	return &Service{
		employeeRepo: employeeRepo,
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		jwt:          jwt,
	}
}

func (s *Service) CheckEmployee(
	ctx context.Context,
	req CheckEmployeeRequest,
) (*CheckEmployeeResponse, error) {
	employeeData, err := s.employeeRepo.FindByEmployeeNo(
		ctx,
		req.EmployeeNo,
	)
	if err != nil {
		if errors.Is(err, employee.ErrEmployeeNotFound) {
			return &CheckEmployeeResponse{
				Status:   CheckEmployeeStatusNotFound,
				Employee: nil,
			}, nil
		}

		return nil, err
	}

	employeeInfo := &EmployeeInfo{
		EmployeeNo: employeeData.EmployeeNo,
		Name:       employeeData.Name,
		Position:   employeeData.Position,
	}

	_, err = s.userRepo.FindByEmployeeNo(
		ctx,
		req.EmployeeNo,
	)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return &CheckEmployeeResponse{
				Status:   CheckEmployeeStatusRegister,
				Employee: employeeInfo,
			}, nil
		}

		return nil, err
	}

	return &CheckEmployeeResponse{
		Status:   CheckEmployeeStatusLogin,
		Employee: employeeInfo,
	}, nil
}

func (s *Service) Register(
	ctx context.Context,
	req RegisterRequest,
) (*TokenResponse, error) {
	// ==========================
	// Validasi Employee
	// ==========================

	_, err := s.employeeRepo.FindByEmployeeNo(ctx, req.EmployeeNo)
	if err != nil {
		if errors.Is(err, employee.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}

		return nil, err
	}

	// ==========================
	// Validasi User
	// ==========================

	_, err = s.userRepo.FindByEmployeeNo(ctx, req.EmployeeNo)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	if !errors.Is(err, user.ErrUserNotFound) {
		return nil, err
	}

	// ==========================
	// Validasi Password
	// ==========================

	if req.Password != req.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}

	// ==========================
	// Hash Password
	// ==========================

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// ==========================
	// Create User
	// ==========================

	newUser := &user.User{
		EmployeeNo: req.EmployeeNo,
		Password:   hashedPassword,
		RoleID:     nil,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// ==========================
	// Issue Token
	// ==========================

	identity := UserIdentity{
		UserID:     newUser.ID,
		EmployeeNo: newUser.EmployeeNo,
		Role:       "",
	}

	return s.issueToken(ctx, identity)
}

func (s *Service) issueToken(
	ctx context.Context,
	identity UserIdentity,
) (*TokenResponse, error) {
	accessToken, err := s.jwt.GenerateAccessToken(identity)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(identity)
	if err != nil {
		return nil, err
	}

	if err := s.tokenRepo.RevokeByUserID(
		ctx,
		identity.UserID,
	); err != nil {
		return nil, err
	}

	rt := &refreshtoken.RefreshToken{
		UserID:    identity.UserID,
		Token:     refreshToken,
		ExpiresAt: s.jwt.RefreshTokenExpiry(),
		Revoked:   false,
	}

	if err := s.tokenRepo.Create(ctx, rt); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwt.AccessTokenDuration().Seconds()),
	}, nil
}

func (s *Service) Login(
	ctx context.Context,
	req LoginRequest,
) (*TokenResponse, error) {
	// ==========================
	// Cari User
	// ==========================

	existingUser, err := s.userRepo.FindByEmployeeNo(
		ctx,
		req.EmployeeNo,
	)
	if err != nil {

		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrInvalidCredential
		}

		return nil, err
	}

	// ==========================
	// Compare Password
	// ==========================

	if err := CheckPassword(
		req.Password,
		existingUser.Password,
	); err != nil {
		return nil, ErrInvalidCredential
	}

	// ==========================
	// JWT Identity
	// ==========================

	identity := UserIdentity{
		UserID:     existingUser.ID,
		EmployeeNo: existingUser.EmployeeNo,
		Role:       "",
	}

	return s.issueToken(
		ctx,
		identity,
	)
}

func (s *Service) RefreshToken(
	ctx context.Context,
	req RefreshTokenRequest,
) (*TokenResponse, error) {
	rt, err := s.tokenRepo.FindByToken(
		ctx,
		req.RefreshToken,
	)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if rt.Revoked {
		return nil, ErrInvalidRefreshToken
	}

	if rt.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidRefreshToken
	}

	userData, err := s.userRepo.FindByID(
		ctx,
		rt.UserID,
	)
	if err != nil {
		return nil, err
	}

	identity := UserIdentity{
		UserID:     userData.ID,
		EmployeeNo: userData.EmployeeNo,
		Role:       "",
	}

	// issueToken akan otomatis revoke token lama
	return s.issueToken(
		ctx,
		identity,
	)
}

func (s *Service) Logout(
	ctx context.Context,
	req LogoutRequest,
) error {
	rt, err := s.tokenRepo.FindByToken(
		ctx,
		req.RefreshToken,
	)
	if err != nil {
		return ErrLogoutFailed
	}

	if err := s.tokenRepo.Revoke(
		ctx,
		rt.ID,
	); err != nil {
		return err
	}

	return nil
}
