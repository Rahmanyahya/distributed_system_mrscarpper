package admin

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type Admin struct {
	UUID      string `json:"uuid" gorm:"column:uuid;type:text;primaryKey"`
	Email     string `json:"email" gorm:"column:email;type:text"`
	Password  string `json:"password" gorm:"column:password;type:text"`
	CreatedAt string `json:"created_at" gorm:"column:created_at;type:text"`
}

func (a *Admin) TableName() string { return "admin" }

type Repostory interface {
	GetByEmail(ctx context.Context, email string) (*Admin, error)
}

type InputLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Usecase interface {
	Login(ctx context.Context, input *InputLogin) (string, error)
}

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}