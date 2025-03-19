-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create project_access table
CREATE TABLE IF NOT EXISTS project_access (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role ENUM('owner', 'admin', 'member', 'viewer') NOT NULL DEFAULT 'member',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_project_user (project_id, user_id)
);

-- Create test_suites table
CREATE TABLE IF NOT EXISTS test_suites (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE KEY unique_suite_name_per_project (project_id, name)
);

-- Create tags table
CREATE TABLE IF NOT EXISTS tags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_tag_name (name)
);

-- Create test_cases table
CREATE TABLE IF NOT EXISTS test_cases (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT NOT NULL,
    suite_id BIGINT,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    preconditions TEXT,
    status ENUM('draft', 'active', 'deprecated') NOT NULL DEFAULT 'draft',
    priority ENUM('low', 'medium', 'high') NOT NULL DEFAULT 'medium',
    created_by BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (suite_id) REFERENCES test_suites(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

-- Create test_steps table for Gherkin-style steps
CREATE TABLE IF NOT EXISTS test_steps (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    test_case_id BIGINT NOT NULL,
    step_number INT NOT NULL,
    step_type ENUM('given', 'when', 'then', 'and', 'but') NOT NULL,
    description TEXT NOT NULL,
    expected_result TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (test_case_id) REFERENCES test_cases(id) ON DELETE CASCADE,
    UNIQUE KEY unique_step_number_per_test (test_case_id, step_number)
);

-- Create step_notes table for notes attached to steps
CREATE TABLE IF NOT EXISTS step_notes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    step_id BIGINT NOT NULL,
    note_text TEXT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (step_id) REFERENCES test_steps(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Create step_attachments table for images and other files
CREATE TABLE IF NOT EXISTS step_attachments (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    step_id BIGINT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (step_id) REFERENCES test_steps(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Create test_case_tags junction table for many-to-many relationship
CREATE TABLE IF NOT EXISTS test_case_tags (
    test_case_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    PRIMARY KEY (test_case_id, tag_id),
    FOREIGN KEY (test_case_id) REFERENCES test_cases(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Create test_case_history table for version control
CREATE TABLE IF NOT EXISTS test_case_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    test_case_id BIGINT NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    preconditions TEXT,
    status ENUM('draft', 'active', 'deprecated') NOT NULL,
    priority ENUM('low', 'medium', 'high') NOT NULL,
    version INT NOT NULL,
    changed_by BIGINT NOT NULL,
    change_summary TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (test_case_id) REFERENCES test_cases(id) ON DELETE CASCADE,
    FOREIGN KEY (changed_by) REFERENCES users(id)
); 