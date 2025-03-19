-- Create environments table for tracking different test environments
CREATE TABLE IF NOT EXISTS environments (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    base_url VARCHAR(255) NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id),
    UNIQUE KEY unique_env_name_per_project (project_id, name)
);

-- Create environment_variables table for storing environment-specific variables
CREATE TABLE IF NOT EXISTS environment_variables (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    environment_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    is_secret BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE,
    UNIQUE KEY unique_var_name_per_env (environment_id, name)
);

-- Add environment_id to test_runs table
ALTER TABLE test_runs
ADD COLUMN environment_id BIGINT NULL AFTER project_id,
ADD CONSTRAINT fk_test_run_environment FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE SET NULL; 