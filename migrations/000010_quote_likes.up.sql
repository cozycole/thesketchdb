CREATE TABLE IF NOT EXISTS quote_likes (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quote_id INT NOT NULL REFERENCES quote(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, quote_id)
);

CREATE INDEX IF NOT EXISTS idx_quote_likes_quote_id ON quote_likes(quote_id);
CREATE INDEX IF NOT EXISTS idx_quote_likes_user_id ON quote_likes(user_id);
