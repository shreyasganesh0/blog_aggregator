-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT feeds_fk FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    
    CONSTRAINT unique_user_feed_pair UNIQUE(feed_id, user_id) 
);

-- +goose Down
DROP TABLE feed_follows;
