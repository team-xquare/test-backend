-- Create junction table for N:M relationship between users and github installations
CREATE TABLE IF NOT EXISTS user_github_installations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    installation_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_installation (user_id, installation_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_installation_id (installation_id)
);

-- Remove user_id from github_installations table since we now use junction table
ALTER TABLE github_installations DROP FOREIGN KEY github_installations_ibfk_1;
ALTER TABLE github_installations DROP COLUMN user_id;

-- Remove unique constraint on installation_id to allow same installation for multiple users
ALTER TABLE github_installations DROP INDEX installation_id;