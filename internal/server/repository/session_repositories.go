package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tempizhere/vaultfactory/internal/shared/interfaces"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
	"github.com/uptrace/bun"
)

type sessionRepository struct {
	db *bun.DB
}

func NewSessionRepository(db *bun.DB) interfaces.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *models.UserSession) error {
	_, err := r.db.NewInsert().Model(session).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *sessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.UserSession, error) {
	session := new(models.UserSession)
	err := r.db.NewSelect().
		Model(session).
		Relation("User").
		Where("refresh_token = ?", refreshToken).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by refresh token: %w", err)
	}
	return session, nil
}

func (r *sessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserSession, error) {
	var sessions []*models.UserSession
	err := r.db.NewSelect().
		Model(&sessions).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions by user id: %w", err)
	}
	return sessions, nil
}

func (r *sessionRepository) Update(ctx context.Context, session *models.UserSession) error {
	_, err := r.db.NewUpdate().Model(session).Where("id = ?", session.ID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

func (r *sessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*models.UserSession)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*models.UserSession)(nil)).Where("user_id = ?", userID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete sessions by user id: %w", err)
	}
	return nil
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*models.UserSession)(nil)).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}
	return nil
}
