CREATE TABLE teams (
                       team_name VARCHAR(255) PRIMARY KEY,
                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
                       user_id VARCHAR(255) PRIMARY KEY,
                       username VARCHAR(255) NOT NULL,
                       team_name VARCHAR(255) REFERENCES teams(team_name) ON DELETE CASCADE,
                       is_active BOOLEAN DEFAULT TRUE,
                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE pull_requests (
                               pull_request_id VARCHAR(255) PRIMARY KEY,
                               pull_request_name VARCHAR(255) NOT NULL,
                               author_id VARCHAR(255) REFERENCES users(user_id),
                               status VARCHAR(50) DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
                               assigned_reviewers JSONB DEFAULT '[]',
                               created_at TIMESTAMP DEFAULT NOW(),
                               merged_at TIMESTAMP NULL
);

CREATE INDEX idx_users_team_active ON users(team_name, is_active);
CREATE INDEX idx_pull_requests_status ON pull_requests(status);
CREATE INDEX idx_pull_requests_author ON pull_requests(author_id);
CREATE INDEX idx_users_active ON users(is_active) WHERE is_active = true;