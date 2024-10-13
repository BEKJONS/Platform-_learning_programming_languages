CREATE TYPE role AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS users (
                       user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       first_name VARCHAR(255) NOT NULL,
                       last_name VARCHAR(255) NOT NULL,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       role role DEFAULT 'user',
                       phone_number VARCHAR(255),
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       deleted_at bigint default 0
);

CREATE TABLE IF NOT EXISTS favorites (
                           user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
                           task_id UUID,
                           added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           PRIMARY KEY (user_id, task_id)
);
