ALTER TABLE urls
    ADD COLUMN user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX idx_urls_user_id ON urls(user_id);