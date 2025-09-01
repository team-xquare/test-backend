CREATE TABLE IF NOT EXISTS applications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    project_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    tier VARCHAR(50) NOT NULL,
    
    -- GitHub Configuration
    github_owner VARCHAR(255),
    github_repo VARCHAR(255),
    github_branch VARCHAR(100) DEFAULT 'main',
    github_installation_id VARCHAR(50),
    github_trigger_paths JSON,
    
    -- Build Configuration
    build_type VARCHAR(50), -- gradle, nodejs, react, vite, vue, nextjs, go, rust, maven, django, flask, docker
    build_config JSON,
    
    -- Endpoints
    endpoints JSON,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
    INDEX idx_project_id (project_id),
    INDEX idx_github_repo (github_owner, github_repo)
);