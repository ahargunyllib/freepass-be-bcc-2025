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
}

type SessionService interface {
	GetSessions(ctx context.Context, query dto.GetSessionsQuery) (dto.GetSessionsResponse, error)
	GetSession(ctx context.Context, query dto.GetSessionEventQuery) (dto.GetSessionEventResponse, error)
	CreateSession(ctx context.Context, req dto.CreateSessionRequest) error
	UpdateSession(ctx context.Context, req dto.UpdateSessionRequest) error
	DeleteSession(ctx context.Context, query dto.DeleteSessionQuery) error
	AcceptSession(ctx context.Context, query dto.AcceptSessionQuery, req dto.AcceptSessionRequest) error
	RejectSession(ctx context.Context, query dto.RejectSessionQuery, req dto.RejectSessionRequest) error
}
