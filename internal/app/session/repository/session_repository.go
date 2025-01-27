package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/contracts"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/entity"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/log"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sessionRepository struct {
	db *sqlx.DB
}

func (s *sessionRepository) Create(ctx context.Context, session *entity.Session) error {
	_, err := s.db.NamedExecContext(
		ctx,
		`
		INSERT INTO sessions
		(id, title, description, start_at, end_at, type, tags, proposer_id, room, meeting_url, capacity)
		VALUES (:id, :title, :description, :start_at, :end_at, :type, :tags,
			:proposer_id, :room, :meeting_url, :capacity)
		`,
		session,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionRepository) FindAll(
	ctx context.Context,
	limit int,
	offset int,
	sortBy string,
	sortOrder string,
	search string,
	sessionType int16,
	tags int16,
	beforeAt time.Time,
	afterAt time.Time,
	proposerID uuid.UUID,
	status int16,
) ([]entity.Session, error) {
	sessions := []entity.Session{}

	query := `SELECT
		sessions.*, proposer.id as "proposer.id", proposer.name as "proposer.name",
		proposer.email as "proposer.email", proposer.role as "proposer.role"
		FROM sessions JOIN users proposer ON proposer.id=sessions.proposer_id
		WHERE 1=1
	`

	if search != "" {
		query += fmt.Sprintf(" AND (title ILIKE %s)", "'%"+search+"%'")
	}

	if sessionType != 0 {
		query += fmt.Sprintf(" AND type = %d", sessionType)
	}

	if tags != 0 {
		query += fmt.Sprintf(" AND (tags & %d) = %d", tags, tags)
	}

	if !beforeAt.IsZero() {
		query += fmt.Sprintf(" AND start_at < %s", beforeAt)
	}

	if !afterAt.IsZero() {
		query += fmt.Sprintf(" AND end_at > %s", afterAt)
	}

	if proposerID != uuid.Nil {
		query += fmt.Sprintf(" AND proposer_id = '%s'", proposerID.String())
	}

	if status != 0 {
		query += fmt.Sprintf(" AND status = %d", status)
	}

	log.Info(log.LogInfo{
		"query": query,
	}, "[SessionRepository] FindAll")

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d", sortBy, sortOrder, limit, offset)

	err := s.db.SelectContext(ctx, &sessions, query)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (s *sessionRepository) Count(
	ctx context.Context,
	search string,
	sessionType int16,
	tags int16,
	beforeAt time.Time,
	afterAt time.Time,
	proposerID uuid.UUID,
	status int16,
) (int64, error) {
	var count int64
	query := "SELECT COUNT(*) FROM sessions WHERE 1=1"

	if search != "" {
		query += fmt.Sprintf(" AND (title ILIKE %s)", "'%"+search+"%'")
	}

	if sessionType != 0 {
		query += fmt.Sprintf(" AND type = %d", sessionType)
	}

	if tags != 0 {
		query += fmt.Sprintf(" AND (tags & %d) = %d", tags, tags)
	}

	if !beforeAt.IsZero() {
		query += fmt.Sprintf(" AND start_at < %s", beforeAt)
	}

	if !afterAt.IsZero() {
		query += fmt.Sprintf(" AND end_at > %s", afterAt)
	}

	if proposerID != uuid.Nil {
		query += fmt.Sprintf(" AND proposer_id = '%s'", proposerID.String())
	}

	if status != 0 {
		query += fmt.Sprintf(" AND status = %d", status)
	}

	log.Info(log.LogInfo{
		query: query,
	}, "[SessionRepository] Count")

	err := s.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *sessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	var session entity.Session
	err := s.db.GetContext(ctx, &session, "SELECT * FROM sessions WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *sessionRepository) Update(ctx context.Context, session *entity.Session) error {
	_, err := s.db.NamedExecContext(
		ctx,
		`
		UPDATE sessions
		SET title = :title, description = :description, type = :type, tags = :tags,
			start_at = :start_at, end_at = :end_at, room = :room, meeting_url = :meeting_url,
			capacity = :capacity, status = :status
		WHERE id = :id
		`,
		session,
	)
	if err != nil {
		return err
	}

	return nil
}

func NewSessionRepository(db *sqlx.DB) contracts.SessionRepository {
	return &sessionRepository{
		db: db,
	}
}
