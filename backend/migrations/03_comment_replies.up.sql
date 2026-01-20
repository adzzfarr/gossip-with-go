-- Add parent_comment_id to comments table for replies
ALTER TABLE comments
ADD COLUMN parent_comment_id INTEGER DEFAULT NULL;

-- Add foreign key constraint for parent_comment_id 
-- Self-referential; can reply to another comment in the same table
ALTER TABLE comments
ADD CONSTRAINT fk_parent_comment
FOREIGN KEY (parent_comment_id)
REFERENCES comments(comment_id)
ON DELETE CASCADE; -- If a parent comment is deleted, its replies are also deleted

CREATE INDEX idx_comments_parent_comment_id ON comments(parent_comment_id);

COMMENT ON COLUMN comments.parent_comment_id IS 'References the parent comment for replies; NULL if top-level comment';