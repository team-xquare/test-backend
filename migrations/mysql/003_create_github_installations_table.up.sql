CREATE TABLE IF NOT EXISTS github_installations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    installation_id VARCHAR(50) UNIQUE NOT NULL,
    account_login VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) NOT NULL,
    permissions JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_installation_id (installation_id)
);