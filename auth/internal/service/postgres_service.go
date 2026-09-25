package service

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/domain"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/repository/postgres"
	"gorm.io/gorm"
)

type DBService struct {
	db *gorm.DB
}

func NewService() (Service, error) {
	db, err := postgres.NewPool(
		postgres.NewConfigMust(),
		&domain.User{},
	)
	if err != nil {
		return nil, err
	}

	return &DBService{
		db: db.DB,
	}, nil
}

func (s *DBService) Authenticate(ctx context.Context, username string, password string) (domain.User, error) {
	user, err := s.readByUsername(ctx, username)
	if err != nil {
		return domain.User{}, fmt.Errorf("service: %w", err)
	}

	if user.CheckPassword(password) {
		return domain.User{}, fmt.Errorf("not found: %w", err)
	}

	return user, nil
}

func (s *DBService) readByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("repository: failed to read user info by username: %w", err)
	}

	return user, nil
}

func (s *DBService) Create(ctx context.Context, user domain.User) (domain.User, error) {
	hashedPsw, err := user.HashPassword()
	if err != nil {
		return domain.User{}, fmt.Errorf("service: %w", err)
	}
	user.Password = hashedPsw

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("repository: failed to create user: %w", err)
	}

	return user, nil
}

func (s *DBService) Delete(ctx context.Context, id int) error {
	if err := s.db.WithContext(ctx).Delete(&domain.User{ID: id}).Error; err != nil {
		return fmt.Errorf("repository: failed to delete user: %w", err)
	}

	return nil
}

func (s *DBService) Read(ctx context.Context, id int) (domain.User, error) {
	var user domain.User
	if err := s.db.WithContext(ctx).First(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("repository: failed to read user info: %w", err)
	}

	return user, nil
}

func (s *DBService) Update(ctx context.Context, user domain.User) (domain.User, error) {
	var existingUser domain.User
	if err := s.db.WithContext(ctx).First(&existingUser, user.ID).Error; err != nil {
		return domain.User{}, fmt.Errorf("repository: user not found: %w", err)
	}

	if user.Username != "" {
		existingUser.Username = user.Username
	}

	if user.Password != "" {
		hashedPsw, err := user.HashPassword()
		if err != nil {
			return domain.User{}, fmt.Errorf("service: failed to hash new password: %w", err)
		}
		existingUser.Password = hashedPsw
	}

	if err := s.db.WithContext(ctx).Save(&existingUser).Error; err != nil {
		return domain.User{}, fmt.Errorf("repository: failed to save user: %w", err)
	}

	return existingUser, nil
}

func (s *DBService) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return db.Close()
}
