CREATE TABLE sessions (
  id VARCHAR(255) PRIMARY KEY,
  proposer_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  type SMALLINT NOT NULL DEFAULT 1, -- 1: workshop,
  tags SMALLINT NOT NULL DEFAULT 1,
  status SMALLINT NOT NULL DEFAULT 1, -- 1: pending, 2: approved, 3: rejected
  start_at TIMESTAMP NOT NULL,
  end_at TIMESTAMP NOT NULL,
  room VARCHAR(255) NULL,
  meeting_url VARCHAR(255) NULL,
  capacity INT NOT NULL,
  image_uri VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  deleted_reason VARCHAR(255) NULL,

  CONSTRAINT room_or_meeting_url CHECK (
    (room IS NOT NULL AND meeting_url IS NULL) OR
    (room IS NULL AND meeting_url IS NOT NULL)
  )
);

CREATE TRIGGER update_sessions_timestamp
BEFORE UPDATE ON sessions
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
