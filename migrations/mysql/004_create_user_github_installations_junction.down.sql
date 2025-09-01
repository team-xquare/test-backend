-- Reverse the junction table migration
DROP TABLE IF EXISTS user_github_installations;

-- Add back user_id column to github_installations
ALTER TABLE github_installations ADD COLUMN user_id INT NOT NULL DEFAULT 0;

-- Add back foreign key constraint
ALTER TABLE github_installations ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

-- Add back unique constraint on installation_id
ALTER TABLE github_installations ADD UNIQUE KEY installation_id (installation_id);