DROP TRIGGER IF EXISTS update_users_timestamp ON users;

DROP INDEX IF EXISTS users_email_index;

DROP TABLE IF EXISTS users;
