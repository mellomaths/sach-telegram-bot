package sacs

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	repo "github.com/mellomaths/sach-telegram-bot/internal/adapters/postgresql/sqlc"
	"go.uber.org/zap"
)

type User struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	UserName  string `json:"username,omitempty"`
}

type Service interface {
	SaveSac(ctx context.Context, u User, msg string) error
}

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{repo: repo}
}

func (s *svc) SaveSac(ctx context.Context, u User, msg string) error {
	zap.L().Info("Saving SAC", zap.Int64("user_id", u.Id), zap.String("message", msg))
	user, err := s.repo.FindUserByID(ctx, u.Id)
	// if user not found, create it
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Info("User not found, creating it", zap.Int64("user_id", u.Id))
			user, err = s.repo.CreateUser(ctx, repo.CreateUserParams{
				ID:        u.Id,
				UserName:  u.UserName,
				FirstName: u.FirstName,
				LastName:  u.LastName,
			})
			if err != nil {
				zap.L().Error("Error creating user", zap.Error(err))
				return err
			}
			zap.L().Info("User created", zap.Any("user", user))
		} else {
			zap.L().Error("Error finding user", zap.Error(err))
			return err
		}
	}
	sac, err := s.repo.CreateSAC(ctx, repo.CreateSACParams{
		UserID:  user.ID,
		Message: msg,
	})
	if err != nil {
		zap.L().Error("Error creating SAC", zap.Error(err), zap.Int64("user_id", user.ID), zap.String("message", msg))
		return err
	}
	zap.L().Info("SAC created", zap.Int64("user_id", user.ID), zap.String("message", msg), zap.Int64("sac_id", sac.ID))
	return nil
}
