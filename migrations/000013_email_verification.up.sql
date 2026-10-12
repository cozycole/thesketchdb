ALTER TABLE token 
ADD COLUMN IF NOT EXISTS scope TEXT NOT NULL DEFAULT 'unknown', 
ADD COLUMN ttl timestamp(0) with time zone;

ALTER TABLE users
ALTER COLUMN activated SET DEFAULT false;
