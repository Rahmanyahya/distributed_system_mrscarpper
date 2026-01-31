-- Drop indexes
DROP INDEX IF EXISTS idx_admin_email;
DROP INDEX IF EXISTS idx_config_created_at;
-- Drop tables
DROP TABLE IF EXISTS admin;
DROP TABLE IF EXISTS config;
DROP TABLE IF EXISTS agents;
