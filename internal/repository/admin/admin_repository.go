package admin

import (
	"context"
	"distributed_system/internal/domain/admin"
	"distributed_system/pkg/errors"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) admin.Repostory {
	return &repository{
		db: db,
	}
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*admin.Admin, error) {
	var admin admin.Admin

	if err := r.db.First(&admin, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("admin")
		}
		return nil, errors.Database(err)
	}

	return &admin, nil
}