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
	args := []interface{}{}

	if search != "" {
		query += fmt.Sprintf(" AND (title ILIKE %$%d%)", len(args)+1)
		args = append(args, search)
	}

	if sessionType != 0 {
		query += fmt.Sprintf(" AND type = $%d", len(args)+1)
		args = append(args, sessionType)
	}

	if tags != 0 {
		query += fmt.Sprintf(" AND (tags & $%d) = $%d", len(args)+1, len(args)+1)
		args = append(args, tags, tags)
	}

	if !beforeAt.IsZero() {
		query += fmt.Sprintf(" AND start_at < $%d", len(args)+1)
		args = append(args, beforeAt)
	}

	if !afterAt.IsZero() {
		query += fmt.Sprintf(" AND end_at > $%d", len(args)+1)
		args = append(args, afterAt)
	}

	if proposerID != uuid.Nil {
		query += fmt.Sprintf(" AND proposer_id = $%d", len(args)+1)
		args = append(args, proposerID)
	}

	if status != 0 {
		query += fmt.Sprintf(" AND status = $%d", len(args)+1)
		args = append(args, status)
	}

	log.Info(log.LogInfo{
		"query": query,
	}, "[SessionRepository] FindAll")

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sortBy, sortOrder, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	err := s.db.SelectContext(ctx, &sessions, query, args...)
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
	args := []interface{}{}

	if search != "" {
		query += fmt.Sprintf(" AND (title ILIKE %$%d%)", len(args)+1)
		args = append(args, search)
	}

	if sessionType != 0 {
		query += fmt.Sprintf(" AND type = $%d", len(args)+1)
		args = append(args, sessionType)
	}

	if tags != 0 {
		query += fmt.Sprintf(" AND (tags & $%d) = $%d", len(args)+1, len(args)+1)
		args = append(args, tags, tags)
	}

	if !beforeAt.IsZero() {
		query += fmt.Sprintf(" AND start_at < $%d", len(args)+1)
		args = append(args, beforeAt)
	}

	if !afterAt.IsZero() {
		query += fmt.Sprintf(" AND end_at > $%d", len(args)+1)
		args = append(args, afterAt)
	}

	if proposerID != uuid.Nil {
		query += fmt.Sprintf(" AND proposer_id = $%d", len(args)+1)
		args = append(args, proposerID)
	}

	if status != 0 {
		query += fmt.Sprintf(" AND status = $%d", len(args)+1)
		args = append(args, status)
	}

	query += " AND deleted_at IS NULL"

	err := s.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *sessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	query := `SELECT
		sessions.*, proposer.id as "proposer.id", proposer.name as "proposer.name",
		proposer.email as "proposer.email", proposer.role as "proposer.role"
		FROM sessions JOIN users proposer ON proposer.id=sessions.proposer_id
		WHERE sessions.id = $1 AND deleted_at IS NULL
	`

	var session entity.Session
	err := s.db.GetContext(ctx, &session, query, id)
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

func (s *sessionRepository) CreateSessionAttendee(ctx context.Context, sessionAttendee *entity.SessionAttendee) error {
	_, err := s.db.NamedExecContext(
		ctx,
		`
		INSERT INTO session_attendees
		(session_id, user_id)
		VALUES (:session_id, :user_id)
		`,
		sessionAttendee,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionRepository) UpdateSessionAttendee(ctx context.Context, sessionAttendee *entity.SessionAttendee) error {
	_, err := s.db.NamedExecContext(
		ctx,
		`
		UPDATE session_attendees
		SET review = :review, reason = :reason, deleted_reason = :deleted_reason
		WHERE session_id = :session_id AND user_id = :user_id
		`,
		sessionAttendee,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionRepository) CountAttendees(
	ctx context.Context,
	sessionID uuid.UUID,
	userID uuid.UUID,
	beforeAt time.Time,
	afterAt time.Time,
	canceled bool,
) (int64, error) {
	var count int64
	query := `SELECT COUNT(session_attendees.*),
		sessions.start_at as start_at, sessions.end_at as end_at
		FROM session_attendees JOIN sessions ON sessions.id=session_attendees.session_id
		WHERE 1=1`
		args := []interface{}{}

	if sessionID != uuid.Nil {
		query += fmt.Sprintf(" AND session_id = %d", len(args)+1)
		args = append(args, sessionID)
	}

	if userID != uuid.Nil {
		query += fmt.Sprintf(" AND user_id = %d", len(args)+1)
		args = append(args, userID)
	}

	if !beforeAt.IsZero() {
		query += fmt.Sprintf(" AND end_at <= %d", len(args)+1)
		args = append(args, beforeAt)
	}

	if !afterAt.IsZero() {
		query += fmt.Sprintf(" AND start_at >= %d", len(args)+1)
		args = append(args, afterAt)
	}

	if canceled {
		query += " AND reason IS NOT NULL AND deleted_reason IS NULL"
	}
	err := s.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *sessionRepository) FindSessionAttendee(ctx context.Context, sessionID, userID uuid.UUID) (*entity.SessionAttendee, error) {
	var sessionAttende entity.SessionAttendee
	err := s.db.GetContext(
		ctx,
		&sessionAttende,
		"SELECT * FROM session_attendees WHERE session_id = $1 AND user_id = $2",
		sessionID,
		userID,
	)
	if err != nil {
		return nil, err
	}

	return &sessionAttende, nil
}

func (s *sessionRepository) FindSessionAttendees(
	ctx context.Context,
	sessionID uuid.UUID,
	userID uuid.UUID,
) ([]entity.SessionAttendee, error) {
	query := "SELECT session_attendees.*"

	if sessionID != uuid.Nil {
		query += `, sessions.proposer_id as "session.proposer_id", sessions.title as "session.title",
			sessions.description as "session.description", sessions.type as "session.type",
			sessions.tags as "session.tags", sessions.status as "session.status",
			sessions.start_at as "session.start_at", sessions.end_at as "session.end_at",
			sessions.room as "session.room", sessions.meeting_url as "session.meeting_url",
			sessions.capacity as "session.capacity", sessions.image_uri as "session.image_uri"
		`
	}

	if userID != uuid.Nil {
		query += `, users.id as "user.id", users.name as "user.name",
			users.email as "user.email", users.role as "user.role"
		`
	}

	query += ` FROM session_attendees`

	if sessionID != uuid.Nil {
		query += ` JOIN sessions ON sessions.id=session_attendees.session_id`
	}

	if userID != uuid.Nil {
		query += ` JOIN users ON users.id=session_attendees.user_id`
	}

	query += " WHERE 1=1"
	args := []interface{}{}

	if sessionID != uuid.Nil {
		query += fmt.Sprintf(" AND session_id = $%d", len(args)+1)
		args = append(args, sessionID)
	}

	if userID != uuid.Nil {
		query += fmt.Sprintf(" AND user_id = $%d", len(args)+1)
		args = append(args, userID)
	}

	sessionAttendees := []entity.SessionAttendee{}
	err := s.db.SelectContext(
		ctx,
		&sessionAttendees,
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}

	return sessionAttendees, nil
}

func NewSessionRepository(db *sqlx.DB) contracts.SessionRepository {
	return &sessionRepository{
		db: db,
	}
}
