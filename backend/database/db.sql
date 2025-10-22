PRAGMA foreign_keys = ON;


-- users

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


-- categories

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);


-- posts
-- (no category_id here; many-to-many via post_categories)

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL CHECK(length(title) > 0),
    content TEXT NOT NULL CHECK(length(content) > 0),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    likes_count INTEGER DEFAULT 0,
    dislikes_count INTEGER DEFAULT 0,
    comments_count INTEGER DEFAULT 0,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);



-- post_categories (many-to-many)

CREATE TABLE IF NOT EXISTS post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE
);


-- comments

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL CHECK(length(content) > 0),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    likes_count INTEGER DEFAULT 0,
    dislikes_count INTEGER DEFAULT 0,
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);




-- likes (unified for posts & comments)

CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    target_type TEXT NOT NULL CHECK(target_type IN ('post','comment')),
    target_id INTEGER NOT NULL,
    value INTEGER NOT NULL CHECK(value IN (1, -1)), -- 1 = like, -1 = dislike
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, target_type, target_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);



-- Note: likes on comments are not reflected on posts counts; if you need per-comment counts, add columns in comments and similar triggers.


-- sessions

CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY, -- UUID stored as TEXT
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK(expires_at > created_at)
);


-- Indexes

CREATE INDEX IF NOT EXISTS idx_posts_user ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_post_categories_cat ON post_categories(category_id);
CREATE INDEX IF NOT EXISTS idx_post_categories_post ON post_categories(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_likes_target ON likes(target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);


-- Seed some example categories (INSERT required by the assignment)

INSERT OR IGNORE INTO categories (name) VALUES ('General'), ('Programming'), ('Announcements'), ('Off-topic');

