CREATE TABLE IF NOT EXISTS addons (
    id INT AUTO_INCREMENT PRIMARY KEY,
    project_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- mysql, redis, postgresql, mongodb, etc.
    tier VARCHAR(50) NOT NULL, -- small, medium, large
    storage VARCHAR(50), -- for databases that need storage specification
    config JSON, -- additional addon-specific configuration
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
    INDEX idx_project_id (project_id),
    INDEX idx_type (type)
);