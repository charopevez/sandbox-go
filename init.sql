-- Sample schema for practicing Go + DB
CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tasks (
    id          SERIAL PRIMARY KEY,
    user_id     INT REFERENCES users(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    done        BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Seed data
INSERT INTO users (name, email) VALUES
    ('Alice', 'alice@example.com'),
    ('Bob', 'bob@example.com'),
    ('Charlie', 'charlie@example.com');

INSERT INTO tasks (user_id, title, done) VALUES
    (1, 'Learn Go basics', TRUE),
    (1, 'Build REST API', FALSE),
    (2, 'Study goroutines', FALSE),
    (2, 'Practice live coding', FALSE),
    (3, 'Read about AWS Glue', FALSE);
