DROP TRIGGER IF EXISTS update_sessions_timestamp ON sessions;

DROP INDEX IF EXISTS sessions_proposer_id_index;
DROP INDEX IF EXISTS sessions_title_index;
DROP INDEX IF EXISTS sessions_type_index;
DROP INDEX IF EXISTS sessions_status_index;
DROP INDEX IF EXISTS sessions_start_at_index;
DROP INDEX IF EXISTS sessions_end_at_index;

DROP TABLE IF EXISTS sessions;
