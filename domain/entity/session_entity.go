package entity

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/enums"
	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID         `db:"id" json:"id"`
	ProposerID       uuid.UUID         `db:"proposer_id" json:"proposer_id"`
	Title            string            `db:"title" json:"title"`
	Description      sql.NullString    `db:"description" json:"description"`
	Type             int16             `db:"type" json:"type"`
	Tags             int16             `db:"tags" json:"tags"`
	Status           int16             `db:"status" json:"status"`
	StartAt          time.Time         `db:"start_at" json:"start_at"`
	EndAt            time.Time         `db:"end_at" json:"end_at"`
	Room             sql.NullString    `db:"room" json:"room"`
	MeetingURL       sql.NullString    `db:"meeting_url" json:"meeting_url"`
	Capacity         int               `db:"capacity" json:"capacity"`
	ImageURI         sql.NullString    `db:"image_uri" json:"image_uri"`
	CreatedAt        time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time         `db:"updated_at" json:"updated_at"`
	DeletedAt        sql.NullTime      `db:"deleted_at" json:"deleted_at"`
	DeletedReason    sql.NullString    `db:"deleted_reason" json:"deleted_reason"`
	Proposer         User              `db:"proposer" json:"proposer"`
}

type SessionAttendee struct {
	SessionID     uuid.UUID      `db:"session_id" json:"session_id"`
	UserID        uuid.UUID      `db:"user_id" json:"user_id"`
	Review        sql.NullString `db:"review" json:"review"`
	Reason        sql.NullString `db:"reason" json:"reason"`
	DeletedReason sql.NullString `db:"deleted_reason" json:"deleted_reason"`
	User          User           `db:"user" json:"user"`
	Session       Session        `db:"session" json:"session"`
}

func (s *Session) TagsArray() []string {
	binary := strconv.FormatInt(int64(s.Tags), 2)

	// convert binary to array of tags
	tags := []string{}
	for i, tag := range binary {
		if tag == '1' {
			tags = append(tags, enums.ShortSessionTag[i])
		}
	}

	return tags
}
