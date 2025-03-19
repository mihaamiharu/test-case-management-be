CREATE TABLE IF NOT EXISTS project_access (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    level ENUM('view', 'edit') NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE KEY idx_project_user (project_id, user_id),
    INDEX idx_user (user_id),
    CONSTRAINT fk_project_access_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_project_access_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
); 