package container

import (
	"github.com/jmoiron/sqlx"

	"github.com/ikhwan11/auth-be/internal/auth"
	"github.com/ikhwan11/auth-be/internal/config"
	"github.com/ikhwan11/auth-be/internal/employee"
	refreshtoken "github.com/ikhwan11/auth-be/internal/refresh_token"
	"github.com/ikhwan11/auth-be/internal/user"
)

type Container struct {
	AuthHandler *auth.Handler
}

func New(
	employeeDB *sqlx.DB,
	authDB *sqlx.DB,
	cfg config.AppConfig,
) (*Container, error) {
	employeeRepo := employee.NewRepository(employeeDB)
	userRepo := user.NewRepository(authDB)
	tokenRepo := refreshtoken.NewRepository(authDB)

	jwtManager, err := auth.NewJWT(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpired,
		cfg.JWT.RefreshExpired,
		cfg.JWT.Issuer,
	)
	if err != nil {
		return nil, err
	}

	authService := auth.NewService(
		employeeRepo,
		userRepo,
		tokenRepo,
		jwtManager,
	)

	authHandler := auth.NewHandler(authService)

	return &Container{
		AuthHandler: authHandler,
	}, nil
}
