package admin

import (
	"context"
	"distributed_system/internal/config"
	"distributed_system/internal/domain/admin"
	"distributed_system/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase struct {
	repository admin.Repostory
	cfg        *config.Config
}

func NewAdminUsecase(repository admin.Repostory, cfg *config.Config) admin.Usecase {
	return &AdminUsecase{repository: repository, cfg: cfg}
}

func (u *AdminUsecase) Login(ctx context.Context, input *admin.InputLogin) (string, error) {
	account, err := u.repository.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.IsNotFound(err) {
			return "", errors.NotFound("admin")
		}

		return "", errors.Wrap(err, "admin", "failed to get admin")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(input.Password)); err != nil {
		return "", errors.Wrap(err, "admin", "invalid password")
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &admin.Claims{Role: "admin"}).SignedString([]byte(u.cfg.Security.JWTSecret))
	if err != nil {
		return "", errors.Wrap(err, "admin", "failed to create token")
	}

	return token, nil
}