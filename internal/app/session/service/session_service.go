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
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/uuid"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/validator"
)

type sessionService struct {
	repo      contracts.SessionRepository
	validator validator.ValidatorInterface
	uuid      uuid.CustomUUIDInterface
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

	id, err := s.uuid.NewV7()
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

	sessionResponse := dto.SessionResponse{
		ID:          session.ID,
		Title:       session.Title,
		Description: &session.Description.String,
		Type:        session.Type,
		Tags:        session.TagsArray(),
		StartAt:     session.StartAt,
		EndAt:       session.EndAt,
		Room:        &session.Room.String,
		MeetingURL:  &session.MeetingURL.String,
		Capacity:    session.Capacity,
		ImageURI:    &session.ImageURI.String,
		Status:      session.Status,
		Proposer: dto.UserResponse{
			ID:    session.Proposer.ID,
			Name:  session.Proposer.Name,
			Email: session.Proposer.Email,
		},
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
			Description: &session.Description.String,
			Type:        session.Type,
			Tags:        session.TagsArray(),
			StartAt:     session.StartAt,
			EndAt:       session.EndAt,
			Room:        &session.Room.String,
			MeetingURL:  &session.MeetingURL.String,
			Capacity:    session.Capacity,
			ImageURI:    &session.ImageURI.String,
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

func NewSessionService(
	repo contracts.SessionRepository,
	validator validator.ValidatorInterface,
	uuid uuid.CustomUUIDInterface,
) contracts.SessionService {
	return &sessionService{
		repo:      repo,
		validator: validator,
		uuid:      uuid,
	}
}
