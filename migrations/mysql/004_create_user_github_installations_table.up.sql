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