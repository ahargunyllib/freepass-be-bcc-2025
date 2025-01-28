package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/ahargunyllib/freepass-be-bcc-2025/domain"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/contracts"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/dto"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/entity"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/log"
	uuidPkg "github.com/ahargunyllib/freepass-be-bcc-2025/pkg/uuid"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/validator"
	"github.com/google/uuid"
)

type sessionService struct {
	repo      contracts.SessionRepository
	validator validator.ValidatorInterface
	uuidPkg   uuidPkg.CustomUUIDInterface
}

func (s *sessionService) AcceptSession(
	ctx context.Context,
	query dto.AcceptSessionQuery,
	req dto.AcceptSessionRequest,
) error {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return valErr
	}

	session, err := s.repo.FindByID(ctx, query.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotFound
		}

		return err
	}

	tagsBinary := "000000"
	for _, tag := range req.Tags {
		switch tag {
		case "PM":
			tagsBinary = tagsBinary[:0] + "1" + tagsBinary[1:]
		case "PD":
			tagsBinary = tagsBinary[:1] + "1" + tagsBinary[2:]
		case "FE":
			tagsBinary = tagsBinary[:2] + "1" + tagsBinary[3:]
		case "BE":
			tagsBinary = tagsBinary[:3] + "1" + tagsBinary[4:]
		case "DS":
			tagsBinary = tagsBinary[:4] + "1" + tagsBinary[5:]
		case "CP":
			tagsBinary = tagsBinary[:5] + "1"
		}
	}

	tagsNumber, err := strconv.ParseInt(tagsBinary, 2, 64)
	if err != nil {
		return err
	}

	if req.Title != "" {
		session.Title = req.Title
	}

	if req.Description != "" {
		session.Description = sql.NullString{String: req.Description, Valid: true}
	}

	if req.Type != 0 {
		session.Type = req.Type
	}

	if len(req.Tags) != 0 {
		session.Tags = int16(tagsNumber)
	}

	if !req.StartAt.IsZero() {
		session.StartAt = req.StartAt
	}

	if !req.EndAt.IsZero() {
		session.EndAt = req.EndAt
	}

	if req.Room != "" {
		session.Room = sql.NullString{String: req.Room, Valid: true}
	}

	if req.MeetingURL != "" {
		session.MeetingURL = sql.NullString{String: req.MeetingURL, Valid: true}
	}

	if req.Capacity != 0 {
		session.Capacity = req.Capacity
	}

	session.Status = 2 // Accepted

	log.Info(log.LogInfo{
		"session": session,
	}, "[SessionService] AcceptSession")

	err = s.repo.Update(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionService) CreateSession(ctx context.Context, req dto.CreateSessionRequest) error {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return valErr
	}

	countPendingSessionProposal, err := s.repo.Count(
		ctx,
		"",
		0,
		0,
		time.Time{},
		time.Time{},
		req.ProposerID,
		1, // pending
	)
	if err != nil {
		return err
	}

	if countPendingSessionProposal > 0 { // If there is a pending session proposal
		return domain.ErrSessionProposalLimit
	}

	id, err := s.uuidPkg.NewV7()
	if err != nil {
		return err
	}

	tagsBinary := "000000"
	for _, tag := range req.Tags {
		switch tag {
		case "PM":
			tagsBinary = tagsBinary[:0] + "1" + tagsBinary[1:]
		case "PD":
			tagsBinary = tagsBinary[:1] + "1" + tagsBinary[2:]
		case "FE":
			tagsBinary = tagsBinary[:2] + "1" + tagsBinary[3:]
		case "BE":
			tagsBinary = tagsBinary[:3] + "1" + tagsBinary[4:]
		case "DS":
			tagsBinary = tagsBinary[:4] + "1" + tagsBinary[5:]
		case "CP":
			tagsBinary = tagsBinary[:5] + "1"
		}
	}
	tagsNumber, err := strconv.ParseInt(tagsBinary, 2, 64)
	if err != nil {
		return err
	}

	session := entity.Session{
		ID:          id,
		ProposerID:  req.ProposerID,
		Title:       req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Type:        req.Type,
		Tags:        int16(tagsNumber),
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		Room:        sql.NullString{String: req.Room, Valid: req.Room != ""},
		MeetingURL:  sql.NullString{String: req.MeetingURL, Valid: req.MeetingURL != ""},
		Capacity:    req.Capacity,
	}

	err = s.repo.Create(ctx, &session)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionService) DeleteSession(ctx context.Context, query dto.DeleteSessionQuery) error {
	valErr := s.validator.Validate(query)
	if valErr != nil {
		return valErr
	}

	session, err := s.repo.FindByID(ctx, query.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotFound
		}

		return err
	}

	if session.Status != 1 {
		return domain.ErrSessionCannotBeDeleted
	}

	err = s.repo.Delete(ctx, query.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionService) GetSession(
	ctx context.Context,
	query dto.GetSessionEventQuery,
) (dto.GetSessionEventResponse, error) {
	valErr := s.validator.Validate(query)
	if valErr != nil {
		return dto.GetSessionEventResponse{}, valErr
	}

	session, err := s.repo.FindByID(ctx, query.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GetSessionEventResponse{}, domain.ErrSessionNotFound
		}

		return dto.GetSessionEventResponse{}, err
	}

	sessionAttendees, err := s.repo.FindSessionAttendees(ctx, query.ID, uuid.Nil)
	if err != nil {
		return dto.GetSessionEventResponse{}, err
	}

	sessionAttendeeResponses := []dto.SessionAtendeeResponse{}

	for _, sessionAttendee := range sessionAttendees {
		sessionAttendeeResponses = append(sessionAttendeeResponses, dto.SessionAtendeeResponse{
			SessionID: sessionAttendee.SessionID,
			User: dto.UserResponse{
				ID:    sessionAttendee.User.ID,
				Name:  sessionAttendee.User.Name,
				Email: sessionAttendee.User.Email,
			},
			Reason: sessionAttendee.Reason.String,
			Review: sessionAttendee.Review.String,
			UserID: sessionAttendee.User.ID,
		})
	}

	sessionResponse := dto.SessionResponse{
		ID:          session.ID,
		Title:       session.Title,
		Description: session.Description.String,
		Type:        session.Type,
		Tags:        session.TagsArray(),
		StartAt:     session.StartAt,
		EndAt:       session.EndAt,
		Room:        session.Room.String,
		MeetingURL:  session.MeetingURL.String,
		Capacity:    session.Capacity,
		ImageURI:    session.ImageURI.String,
		Status:      session.Status,
		Proposer: dto.UserResponse{
			ID:    session.Proposer.ID,
			Name:  session.Proposer.Name,
			Email: session.Proposer.Email,
		},
		SessionAtendees: sessionAttendeeResponses,
	}

	res := dto.GetSessionEventResponse{
		Session: sessionResponse,
	}

	return res, nil
}

func (s *sessionService) GetSessions(ctx context.Context, query dto.GetSessionsQuery) (dto.GetSessionsResponse, error) {
	valErr := s.validator.Validate(query)
	if valErr != nil {
		return dto.GetSessionsResponse{}, valErr
	}

	if query.Limit < 1 {
		query.Limit = 10
	}

	if query.Page < 1 {
		query.Page = 1
	}

	if query.SortBy == "" {
		query.SortBy = "start_at"
	}

	if query.SortOrder == "" {
		query.SortOrder = "ASC"
	}

	tagsBinary := "000000"
	if len(query.Tags) != 0 {
		for _, tag := range query.Tags {
			switch tag {
			case "PM":
				tagsBinary = tagsBinary[:0] + "1" + tagsBinary[1:]
			case "PD":
				tagsBinary = tagsBinary[:1] + "1" + tagsBinary[2:]
			case "FE":
				tagsBinary = tagsBinary[:2] + "1" + tagsBinary[3:]
			case "BE":
				tagsBinary = tagsBinary[:3] + "1" + tagsBinary[4:]
			case "DS":
				tagsBinary = tagsBinary[:4] + "1" + tagsBinary[5:]
			case "CP":
				tagsBinary = tagsBinary[:5] + "1"
			}
		}
	}

	tagsNumber, err := strconv.ParseInt(tagsBinary, 2, 64)
	if err != nil {
		return dto.GetSessionsResponse{}, err
	}

	sessions, err := s.repo.FindAll(
		ctx,
		query.Limit,
		query.Limit*(query.Page-1),
		query.SortBy,
		query.SortOrder,
		query.Search,
		query.Type,
		int16(tagsNumber&0xFFFF), // 0xFFFF is used to get the last 16 bits
		query.BeforeAt,
		query.AfterAt,
		query.ProposerID,
		query.Status,
	)
	if err != nil {
		return dto.GetSessionsResponse{}, err
	}

	totalData, err := s.repo.Count(
		ctx,
		query.Search,
		query.Type,
		int16(tagsNumber&0xFFFF), // 0xFFFF is used to get the last 16 bits
		query.BeforeAt,
		query.AfterAt,
		query.ProposerID,
		query.Status,
	)
	if err != nil {
		return dto.GetSessionsResponse{}, err
	}

	totalPage := int(totalData) / query.Limit
	if int(totalData)%query.Limit != 0 {
		totalPage++
	}

	meta := dto.PaginationResponse{
		TotalData: totalData,
		TotalPage: totalPage,
		Page:      query.Page,
		Limit:     query.Limit,
	}

	sessionsResponse := []dto.SessionResponse{}
	for _, session := range sessions {
		sessionsResponse = append(sessionsResponse, dto.SessionResponse{
			ID:          session.ID,
			Title:       session.Title,
			Description: session.Description.String,
			Type:        session.Type,
			Tags:        session.TagsArray(),
			StartAt:     session.StartAt,
			EndAt:       session.EndAt,
			Room:        session.Room.String,
			MeetingURL:  session.MeetingURL.String,
			Capacity:    session.Capacity,
			ImageURI:    session.ImageURI.String,
			Status:      session.Status,
			Proposer: dto.UserResponse{
				ID:    session.Proposer.ID,
				Name:  session.Proposer.Name,
				Email: session.Proposer.Email,
			},
		})
	}

	res := dto.GetSessionsResponse{
		Sessions: sessionsResponse,
		Meta:     meta,
	}

	return res, nil
}

func (s *sessionService) RejectSession(
	ctx context.Context,
	query dto.RejectSessionQuery,
	req dto.RejectSessionRequest,
) error {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return valErr
	}

	session, err := s.repo.FindByID(ctx, query.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotFound
		}

		return err
	}

	session.Status = 3 // Rejected
	session.DeletedReason = sql.NullString{String: req.Reason, Valid: true}

	err = s.repo.Update(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionService) UpdateSession(ctx context.Context, req dto.UpdateSessionRequest) error {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return valErr
	}

	session, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotFound
		}

		return err
	}

	if session.Status != 1 {
		return domain.ErrSessionCannotBeUpdated
	}

	tagsBinary := "000000"
	for _, tag := range req.Tags {
		switch tag {
		case "PM":
			tagsBinary = tagsBinary[:0] + "1" + tagsBinary[1:]
		case "PD":
			tagsBinary = tagsBinary[:1] + "1" + tagsBinary[2:]
		case "FE":
			tagsBinary = tagsBinary[:2] + "1" + tagsBinary[3:]
		case "BE":
			tagsBinary = tagsBinary[:3] + "1" + tagsBinary[4:]
		case "DS":
			tagsBinary = tagsBinary[:4] + "1" + tagsBinary[5:]
		case "CP":
			tagsBinary = tagsBinary[:5] + "1"
		}
	}

	tagsNumber, err := strconv.ParseInt(tagsBinary, 2, 64)
	if err != nil {
		return err
	}

	if req.Title != "" {
		session.Title = req.Title
	}

	if req.Description != "" {
		session.Description = sql.NullString{String: req.Description, Valid: true}
	}

	if req.Type != 0 {
		session.Type = req.Type
	}

	if len(req.Tags) != 0 {
		session.Tags = int16(tagsNumber)
	}

	if !req.StartAt.IsZero() {
		session.StartAt = req.StartAt
	}

	if !req.EndAt.IsZero() {
		session.EndAt = req.EndAt
	}

	if req.Room != "" {
		session.Room = sql.NullString{String: req.Room, Valid: true}
	}

	if req.MeetingURL != "" {
		session.MeetingURL = sql.NullString{String: req.MeetingURL, Valid: true}
	}

	if req.Capacity != 0 {
		session.Capacity = req.Capacity
	}

	err = s.repo.Update(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionService) RegisterSession(
	ctx context.Context,
	query dto.RegisterSessionQuery,
	req dto.RegisterSessionRequest,
) error {
	valErr := s.validator.Validate(query)
	if valErr != nil {
		return valErr
	}

	session, err := s.repo.FindByID(ctx, query.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotFound
		}

		return err
	}

	pastSessionAttendee, err := s.repo.FindSessionAttendee(ctx, query.SessionID, req.UserID)
	if err == nil {
		if pastSessionAttendee.Reason.Valid {
			return domain.ErrSessionCancelled
		}

		return domain.ErrSessionAlreadyRegistered
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if session.Status != 2 {
		return domain.ErrSessionNotAccepted
	}

	now := time.Now()
	if session.StartAt.Before(now) {
		return domain.ErrSessionAlreadyStarted
	}

	if session.EndAt.Before(now) {
		return domain.ErrSessionAlreadyEnded
	}

	countSessionAttendees, err := s.repo.CountAttendees(
		ctx,
		query.SessionID,
		uuid.Nil,
		session.EndAt,
		session.StartAt,
		false,
	)
	if err != nil {
		return err
	}

	if int(countSessionAttendees) >= session.Capacity {
		return domain.ErrSessionFull
	}

	countUserSessionTimeConflict, err := s.repo.CountAttendees(
		ctx,
		uuid.Nil,
		req.UserID,
		session.EndAt,
		session.StartAt,
		true,
	)
	if err != nil {
		return err
	}

	if countUserSessionTimeConflict > 0 {
		return domain.ErrSessionTimeConflict
	}

	sessionAttendee := entity.SessionAttendee{
		SessionID: query.SessionID,
		UserID:    req.UserID,
	}

	err = s.repo.CreateSessionAttendee(ctx, &sessionAttendee)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionService) UnregisterSession(
	ctx context.Context,
	query dto.UnregisterSessionQuery,
	req dto.UnregisterSessionRequest,
) error {
	valErr := s.validator.Validate(query)
	if valErr != nil {
		return valErr
	}

	session, err := s.repo.FindByID(ctx, query.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotFound
		}

		return err
	}

	if session.Status != 2 {
		return domain.ErrSessionNotAccepted
	}

	now := time.Now()
	if session.StartAt.Before(now) {
		return domain.ErrSessionAlreadyStarted
	}

	if session.EndAt.Before(now) {
		return domain.ErrSessionAlreadyEnded
	}

	sessionAttendee, err := s.repo.FindSessionAttendee(ctx, query.SessionID, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotRegistered
		}

		return err
	}

	if sessionAttendee.Reason.Valid {
		return domain.ErrSessionCancelled
	}

	sessionAttendee.Reason = sql.NullString{String: req.Reason, Valid: true}

	err = s.repo.UpdateSessionAttendee(ctx, sessionAttendee)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionService) ReviewSession(
	ctx context.Context,
	query dto.ReviewSessionQuery,
	req dto.ReviewSessionRequest,
) error {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return valErr
	}

	session, err := s.repo.FindByID(ctx, query.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotFound
		}

		return err
	}

	if session.Status != 2 {
		return domain.ErrSessionNotAccepted
	}

	now := time.Now()
	if session.StartAt.After(now) {
		return domain.ErrSessionNotStarted
	}

	if session.EndAt.After(now) {
		return domain.ErrSessionNotEnded
	}

	sessionAttendee, err := s.repo.FindSessionAttendee(ctx, query.SessionID, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotRegistered
		}

		return err
	}

	if sessionAttendee.Reason.Valid {
		return domain.ErrSessionCancelled
	}

	if sessionAttendee.DeletedReason.Valid {
		return domain.ErrReviewDeleted
	}

	if sessionAttendee.Review.Valid {
		return domain.ErrSessionAlreadyReviewed
	}

	sessionAttendee.Review = sql.NullString{String: req.Review, Valid: true}

	err = s.repo.Update(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionService) DeleteReviewSession(
	ctx context.Context,
	query dto.DeleteReviewSessionQuery,
	req dto.DeleteReviewSessionRequest,
) error {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return valErr
	}

	session, err := s.repo.FindByID(ctx, query.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotFound
		}

		return err
	}

	if session.Status != 2 {
		return domain.ErrSessionNotAccepted
	}

	sessionAttendee, err := s.repo.FindSessionAttendee(ctx, query.SessionID, query.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotRegistered
		}

		return err
	}

	if sessionAttendee.Reason.Valid {
		return domain.ErrSessionCancelled
	}

	if sessionAttendee.DeletedReason.Valid {
		return domain.ErrReviewDeleted
	}

	sessionAttendee.DeletedReason = sql.NullString{String: req.Reason, Valid: true}

	err = s.repo.Update(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionService) CancelSession(
	ctx context.Context,
	query dto.CancelSessionQuery,
) error {
	valErr := s.validator.Validate(query)
	if valErr != nil {
		return valErr
	}

	session, err := s.repo.FindByID(ctx, query.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrSessionNotFound
		}

		return err
	}

	if session.Status != 2 {
		return domain.ErrSessionNotAccepted
	}

	now := time.Now()
	if session.StartAt.Before(now) {
		return domain.ErrSessionAlreadyStarted
	}

	if session.EndAt.Before(now) {
		return domain.ErrSessionAlreadyEnded
	}

	session.DeletedAt = sql.NullTime{Time: time.Now(), Valid: true}

	err = s.repo.Update(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func NewSessionService(
	repo contracts.SessionRepository,
	validator validator.ValidatorInterface,
	uuidPkg uuidPkg.CustomUUIDInterface,
) contracts.SessionService {
	return &sessionService{
		repo:      repo,
		validator: validator,
		uuidPkg:   uuidPkg,
	}
}
