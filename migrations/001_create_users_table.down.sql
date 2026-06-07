-- Rollback de la migración 001

DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
