package contracts

import (
	"context"
	"time"

	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/dto"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/entity"
	"github.com/google/uuid"
)

type SessionRepository interface {
	FindAll(
		ctx context.Context,
		limit, offset int,
		sortBy, sortOrder, search string,
		sessionType int16,
		tags int16,
		beforeAt time.Time,
		afterAt time.Time,
		proposerID uuid.UUID,
		status int16,
	) ([]entity.Session, error)
	Count(
		ctx context.Context,
		search string,
		sessionType int16,
		tags int16,
		beforeAt time.Time,
		afterAt time.Time,
		proposerID uuid.UUID,
		status int16,
	) (int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Session, error)
	Create(ctx context.Context, session *entity.Session) error
	Update(ctx context.Context, session *entity.Session) error
	Delete(ctx context.Context, id uuid.UUID) error

	CountAttendees(
		ctx context.Context,
		sessionID uuid.UUID,
		userID uuid.UUID,
		beforeAt time.Time,
		afterAt time.Time,
		canceled bool,
	) (int64, error)
	FindSessionAttende(ctx context.Context, sessionID, userID uuid.UUID) (*entity.SessionAttendee, error)
	FindSessionAttendees(
		ctx context.Context,
		sessionID uuid.UUID,
		userID uuid.UUID,
	) ([]entity.SessionAttendee, error)
	CreateSessionAttende(ctx context.Context, sessionAttende *entity.SessionAttendee) error
	UpdateSessionAttende(ctx context.Context, sessionAttende *entity.SessionAttendee) error
}

type SessionService interface {
	GetSessions(ctx context.Context, query dto.GetSessionsQuery) (dto.GetSessionsResponse, error)
	GetSession(ctx context.Context, query dto.GetSessionEventQuery) (dto.GetSessionEventResponse, error)
	CreateSession(ctx context.Context, req dto.CreateSessionRequest) error
	UpdateSession(ctx context.Context, req dto.UpdateSessionRequest) error
	DeleteSession(ctx context.Context, query dto.DeleteSessionQuery) error
	CancelSession(ctx context.Context, query dto.CancelSessionQuery) error
	AcceptSession(ctx context.Context, query dto.AcceptSessionQuery, req dto.AcceptSessionRequest) error
	RejectSession(ctx context.Context, query dto.RejectSessionQuery, req dto.RejectSessionRequest) error

	RegisterSession(ctx context.Context, query dto.RegisterSessionQuery, req dto.RegisterSessionRequest) error
	UnregisterSession(ctx context.Context, query dto.UnregisterSessionQuery, req dto.UnregisterSessionRequest) error
	ReviewSession(ctx context.Context, query dto.ReviewSessionQuery, req dto.ReviewSessionRequest) error
	DeleteReviewSession(ctx context.Context, query dto.DeleteReviewSessionQuery, req dto.DeleteReviewSessionRequest) error
}
