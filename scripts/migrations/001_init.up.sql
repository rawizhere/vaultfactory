-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create user_sessions table
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create data_items table
CREATE TABLE data_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('login_password', 'text_data', 'binary_data', 'bank_card')),
    name VARCHAR(255) NOT NULL,
    metadata TEXT,
    encrypted_data BYTEA NOT NULL,
    encryption_key BYTEA NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version BIGINT DEFAULT 1
);

-- Create data_versions table
CREATE TABLE data_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data_id UUID NOT NULL REFERENCES data_items(id) ON DELETE CASCADE,
    version BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for users table
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Create indexes for user_sessions table
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_refresh_token ON user_sessions(refresh_token);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

-- Create indexes for data_items table
CREATE INDEX idx_data_items_user_id ON data_items(user_id);
CREATE INDEX idx_data_items_type ON data_items(type);
CREATE INDEX idx_data_items_updated_at ON data_items(updated_at);
CREATE INDEX idx_data_items_user_id_updated_at ON data_items(user_id, updated_at);

-- Create indexes for data_versions table
CREATE INDEX idx_data_versions_data_id ON data_versions(data_id);
CREATE INDEX idx_data_versions_version ON data_versions(version);
CREATE INDEX idx_data_versions_data_id_version ON data_versions(data_id, version);


