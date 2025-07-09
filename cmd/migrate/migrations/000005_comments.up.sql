CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    post_id bigserial NOT NULL,
    user_id bigserial NOT NULL,
    content text NOT NULL,
    likes int NOT NULL DEFAULT 0,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
