DROP TRIGGER IF EXISTS update_vote_count_trigger ON votes;
DROP FUNCTION IF EXISTS update_vote_count();

ALTER TABLE posts DROP COLUMN IF EXISTS vote_count;
ALTER TABLE comments DROP COLUMN IF EXISTS vote_count;

DROP TABLE IF EXISTS votes;