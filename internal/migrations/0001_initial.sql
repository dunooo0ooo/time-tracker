CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       passport_number VARCHAR(20) NOT NULL,
                       surname VARCHAR(50),
                       name VARCHAR(50),
                       patronymic VARCHAR(50),
                       address TEXT,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       user_id INTEGER NOT NULL REFERENCES users(id),
                       description TEXT,
                       start_time TIMESTAMP,
                       end_time TIMESTAMP,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
