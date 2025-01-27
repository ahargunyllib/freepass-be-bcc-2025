CREATE TABLE session_attendees (
  session_id VARCHAR(255) NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  review VARCHAR(255) NULL,
  reason VARCHAR(255) NULL,
  deleted_reason VARCHAR(255) NULL,
  CHECK (
    (review IS NOT NULL AND reason IS NULL) OR
    (review IS NULL AND reason IS NOT NULL)
  ),
  PRIMARY KEY (session_id, user_id)
);
