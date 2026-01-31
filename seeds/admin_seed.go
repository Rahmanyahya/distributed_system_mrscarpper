package seeds

import (
	"context"
	"distributed_system/internal/domain/admin"
	"distributed_system/internal/infrastructure/database"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AdminSeedData adalah data default untuk admin
type AdminSeedData struct {
	Email    string
	Password string
}

// DefaultAdminSeed mengembalikan data default admin
func DefaultAdminSeed() *AdminSeedData {
	return &AdminSeedData{
		Email:    "admin@distributed-system.com",
		Password: "Admin123!@#", // Default password, sebaiknya diubah setelah first login
	}
}

// AdminSeed menangani seeding data admin
type AdminSeed struct {
	db      *database.Database
	context context.Context
}

// NewAdminSeed membuat instance baru AdminSeed
func NewAdminSeed(db *database.Database) *AdminSeed {
	return &AdminSeed{
		db:      db,
		context: context.Background(),
	}
}

// SeedAdmin menjalankan proses seeding admin
func (s *AdminSeed) SeedAdmin(data *AdminSeedData) error {
	// Cek apakah admin sudah ada
	var existingAdmin admin.Admin
	result := s.db.DB.Where("email = ?", data.Email).First(&existingAdmin)

	if result.Error == nil {
		log.Printf("Admin dengan email %s sudah ada, melewati seeding...\n", data.Email)
		return nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal melakukan hash password: %w", err)
	}

	// Generate UUID
	adminUUID := uuid.New().String()

	// Timestamp
	now := time.Now().Format(time.RFC3339)

	// Buat admin baru
	newAdmin := &admin.Admin{
		UUID:      adminUUID,
		Email:     data.Email,
		Password:  string(hashedPassword),
		CreatedAt: now,
	}

	// Simpan ke database
	if err := s.db.DB.Create(newAdmin).Error; err != nil {
		return fmt.Errorf("gagal menyimpan admin: %w", err)
	}

	log.Printf("✅ Admin berhasil dibuat:\n")
	log.Printf("   UUID: %s\n", adminUUID)
	log.Printf("   Email: %s\n", data.Email)
	log.Printf("   Password: %s\n", data.Password)
	log.Printf("   ⚠️  Harap ubah password setelah login pertama!\n")

	return nil
}

// SeedMultipleAdmins menjalankan proses seeding untuk multiple admins
func (s *AdminSeed) SeedMultipleAdmins(admins []*AdminSeedData) error {
	for _, adminData := range admins {
		if err := s.SeedAdmin(adminData); err != nil {
			return err
		}
	}
	return nil
}

// TruncateAdmins menghapus semua data admin (gunakan dengan hati-hati)
func (s *AdminSeed) TruncateAdmins() error {
	return s.db.DB.Exec("DELETE FROM admin").Error
}

// RunAdminSeed adalah fungsi helper untuk menjalankan seeding dari main atau migration
func RunAdminSeed(db *database.Database) error {
	seed := NewAdminSeed(db)
	defaultAdmin := DefaultAdminSeed()

	return seed.SeedAdmin(defaultAdmin)
}

// RunMultipleAdminSeeds adalah fungsi helper untuk menjalankan seeding multiple admins
func RunMultipleAdminSeeds(db *database.Database, admins []*AdminSeedData) error {
	seed := NewAdminSeed(db)
	return seed.SeedMultipleAdmins(admins)
}

