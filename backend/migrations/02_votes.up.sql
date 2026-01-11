CREATE TABLE votes (
    vote_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    post_id INT REFERENCES posts(post_id) ON DELETE CASCADE,
    comment_id INT REFERENCES comments(comment_id) ON DELETE CASCADE,
    vote_type INT NOT NULL CHECK (vote_type IN (-1, 1)),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, post_id), -- A user can vote only once per post
    UNIQUE(user_id, comment_id), -- A user can vote only once per comment

    CHECK (
        (post_id IS NOT NULL AND comment_id IS NULL) OR
        (post_id IS NULL AND comment_id IS NOT NULL)
    )
);

CREATE INDEX idx_votes_user_id ON votes(user_id);
CREATE INDEX idx_votes_post_id ON votes(post_id);
CREATE INDEX idx_votes_comment_id ON votes(comment_id);

ALTER TABLE posts ADD COLUMN vote_count INT NOT NULL DEFAULT 0;
ALTER TABLE comments ADD COLUMN vote_count INT NOT NULL DEFAULT 0;

CREATE OR REPLACE FUNCTION update_vote_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.post_id IS NOT NULL THEN
            UPDATE posts
            SET vote_count = vote_count + NEW.vote_type
            WHERE post_id = NEW.post_id;
        ELSIF NEW.comment_id IS NOT NULL THEN
            UPDATE comments
            SET vote_count = vote_count + NEW.vote_type
            WHERE comment_id = NEW.comment_id;
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        IF NEW.post_id IS NOT NULL THEN
            UPDATE posts
            SET vote_count = vote_count + (NEW.vote_type - OLD.vote_type)
            WHERE post_id = NEW.post_id;
        ELSIF NEW.comment_id IS NOT NULL THEN
            UPDATE comments
            SET vote_count = vote_count + (NEW.vote_type - OLD.vote_type)
            WHERE comment_id = NEW.comment_id;
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        IF OLD.post_id IS NOT NULL THEN
            UPDATE posts
            SET vote_count = vote_count - OLD.vote_type
            WHERE post_id = OLD.post_id;
        ELSIF OLD.comment_id IS NOT NULL THEN
            UPDATE comments
            SET vote_count = vote_count - OLD.vote_type
            WHERE comment_id = OLD.comment_id;
        END IF;
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_vote_count_trigger
AFTER INSERT OR UPDATE OR DELETE ON votes
FOR EACH ROW
EXECUTE FUNCTION update_vote_count();
