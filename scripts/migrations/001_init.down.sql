-- Drop indexes for data_versions table
DROP INDEX IF EXISTS idx_data_versions_data_id_version;
DROP INDEX IF EXISTS idx_data_versions_version;
DROP INDEX IF EXISTS idx_data_versions_data_id;

-- Drop indexes for data_items table
DROP INDEX IF EXISTS idx_data_items_user_id_updated_at;
DROP INDEX IF EXISTS idx_data_items_updated_at;
DROP INDEX IF EXISTS idx_data_items_type;
DROP INDEX IF EXISTS idx_data_items_user_id;

-- Drop indexes for user_sessions table
DROP INDEX IF EXISTS idx_user_sessions_expires_at;
DROP INDEX IF EXISTS idx_user_sessions_refresh_token;
DROP INDEX IF EXISTS idx_user_sessions_user_id;

-- Drop indexes for users table
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS data_versions;
DROP TABLE IF EXISTS data_items;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS users;

-- Drop UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";


