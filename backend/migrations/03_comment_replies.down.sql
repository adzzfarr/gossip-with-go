DROP INDEX IF EXISTS idx_comments_parent_comment_id;

ALTER TABLE comments
DROP CONSTRAINT IF EXISTS fk_parent_comment;

ALTER TABLE comments
DROP COLUMN IF EXISTS parent_comment_id;