-- Add release_date column to balls table
ALTER TABLE balls
    ADD COLUMN release_date TIMESTAMP;

-- Add comment to the column
COMMENT ON COLUMN balls.release_date IS 'Release date of the ball (optional)';