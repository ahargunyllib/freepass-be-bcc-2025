package dto

import (
	"time"

	"github.com/google/uuid"
)

type SessionResponse struct {
	ID              uuid.UUID                `json:"id"`
	Title           string                   `json:"title"`
	Description     *string                  `json:"description"`
	Type            int16                    `json:"type"`
	Tags            []string                 `json:"tags"`
	StartAt         time.Time                `json:"start_at"`
	EndAt           time.Time                `json:"end_at"`
	Room            *string                  `json:"room"`
	Status          int16                    `json:"status"`
	MeetingURL      *string                  `json:"meeting_url"`
	Capacity        int                      `json:"capacity"`
	ImageURI        *string                  `json:"image_uri"`
	Proposer        UserResponse             `json:"proposer"`
	SessionAtendees []SessionAtendeeResponse `json:"session_attendees"`
}

type SessionAtendeeResponse struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	Review    *string   `json:"review"`
	Reason    *string   `json:"reason"`
	User      UserResponse
	Session   SessionResponse
}

type GetSessionsQuery struct {
	Search     string    `query:"search" validate:"omitempty,max=255"`
	Type       int16     `query:"type" validate:"omitempty,numeric,oneof= 1"`
	Tags       []string  `query:"tags" validate:"omitempty,dive,oneof=PM PD FE BE DS CP"`
	Limit      int       `query:"limit" validate:"omitempty,numeric,min=1,max=100"`
	Page       int       `query:"page" validate:"omitempty,numeric,min=1"`
	SortBy     string    `query:"sort_by" validate:"omitempty,oneof=id title start_at end_at room capacity"`
	SortOrder  string    `query:"sort_order" validate:"omitempty,oneof=asc desc"`
	BeforeAt   time.Time `query:"before_at" validate:"omitempty"`
	AfterAt    time.Time `query:"after_at" validate:"omitempty"`
	Status     int16     `query:"status" validate:"omitempty,numeric,oneof=1 2 3"`
	ProposerID uuid.UUID `query:"proposer_id" validate:"omitempty,uuid"`
}

type GetSessionsResponse struct {
	Sessions []SessionResponse  `json:"sessions"`
	Meta     PaginationResponse `json:"meta"`
}

type GetSessionEventQuery struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

type GetSessionEventResponse struct {
	Session SessionResponse `json:"session"`
}

type CreateSessionRequest struct {
	ProposerID  uuid.UUID
	Title       string    `json:"title" validate:"required,min=3,max=255"`
	Description string    `json:"description" validate:"omitempty,max=255"`
	Type        int16     `json:"type" validate:"required,numeric,oneof= 1"`
	Tags        []string  `json:"tags" validate:"omitempty,dive,oneof=PM PD FE BE DS CP"`
	StartAt     time.Time `json:"start_at" validate:"required"`
	EndAt       time.Time `json:"end_at" validate:"required,gtefield=StartAt"`
	Room        string    `json:"room" validate:"omitempty,max=255"`
	MeetingURL  string    `json:"meeting_url" validate:"omitempty,url"`
	Capacity    int       `json:"capacity" validate:"required,numeric,min=1,max=100"`
}

type UpdateSessionRequest struct {
	ID          uuid.UUID `param:"id" validate:"required,uuid"`
	Title       string    `json:"title" validate:"omitempty,min=3,max=255"`
	Description string    `json:"description" validate:"omitempty,max=255"`
	Type        int16     `json:"type" validate:"omitempty,numeric,oneof= 1"`
	Tags        []string  `json:"tags" validate:"omitempty,dive,oneof=PM PD FE BE DS CP"`
	StartAt     time.Time `json:"start_at" validate:"omitempty"`
	EndAt       time.Time `json:"end_at" validate:"omitempty,gtefield=StartAt"`
	Room        string    `json:"room" validate:"omitempty,max=255"`
	MeetingURL  string    `json:"meeting_url" validate:"omitempty,url"`
	Capacity    int       `json:"capacity" validate:"omitempty,numeric,min=1,max=100"`
}

type DeleteSessionQuery struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

type AcceptSessionQuery struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

type AcceptSessionRequest struct {
	Title       string    `json:"title" validate:"omitempty,min=3,max=255"`
	Description string    `json:"description" validate:"omitempty,max=255"`
	Type        int16     `json:"type" validate:"omitempty,numeric,oneof= 1"`
	Tags        []string  `json:"tags" validate:"omitempty,dive,oneof=PM PD FE BE DS CP"`
	StartAt     time.Time `json:"start_at" validate:"omitempty"`
	EndAt       time.Time `json:"end_at" validate:"omitempty,gtefield=StartAt"`
	Room        string    `json:"room" validate:"omitempty,max=255"`
	MeetingURL  string    `json:"meeting_url" validate:"omitempty,url"`
	Capacity    int       `json:"capacity" validate:"omitempty,numeric,min=1,max=100"`
}

type RejectSessionQuery struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

type RejectSessionRequest struct {
	Reason string `json:"reason" validate:"required,min=3,max=255"`
}

type RegisterSessionQuery struct {
	SessionID uuid.UUID `param:"sessionID" validate:"required,uuid"`
}

type RegisterSessionRequest struct {
	UserID uuid.UUID // from context
}

type UnregisterSessionQuery struct {
	SessionID uuid.UUID `param:"sessionID" validate:"required,uuid"`
}

type UnregisterSessionRequest struct {
	UserID uuid.UUID // from context
	Reason string `json:"reason" validate:"required,min=3,max=255"`
}

type ReviewSessionQuery struct {
	SessionID uuid.UUID `param:"sessionID" validate:"required,uuid"`
}

type ReviewSessionRequest struct {
	UserID uuid.UUID // from context
	Review string `json:"review" validate:"required,min=3,max=255"`
}

type DeleteReviewSessionQuery struct {
	SessionID uuid.UUID `param:"sessionID" validate:"required,uuid"`
	UserID    uuid.UUID `param:"userID" validate:"required,uuid"`
}

type DeleteReviewSessionRequest struct {
	Reason string `json:"reason" validate:"required,min=3,max=255"`
}

type CancelSessionQuery struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}
