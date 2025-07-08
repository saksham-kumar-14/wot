ALTER TABLE posts

ADD CONSTRAINT foreign_key_user FOREIGN KEY (user_id) REFERENCES users(id);
