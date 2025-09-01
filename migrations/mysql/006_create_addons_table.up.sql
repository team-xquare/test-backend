CREATE TABLE IF NOT EXISTS addons (
    id INT AUTO_INCREMENT PRIMARY KEY,
    project_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- mysql, redis
    tier VARCHAR(50) NOT NULL, -- x3.micro, x3.small, x3.medium, x3.large
    storage VARCHAR(50), -- 1Gi, 5Gi, 10Gi, 20Gi, 50Gi
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
    INDEX idx_project_id (project_id),
    INDEX idx_type (type)
);