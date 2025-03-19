-- Create test_runs table for tracking test execution
CREATE TABLE IF NOT EXISTS test_runs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status ENUM('planned', 'in_progress', 'completed', 'aborted') NOT NULL DEFAULT 'planned',
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Create test_executions table for individual test case executions
CREATE TABLE IF NOT EXISTS test_executions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    test_run_id BIGINT NOT NULL,
    test_case_id BIGINT NOT NULL,
    status ENUM('pending', 'passed', 'failed', 'blocked', 'skipped') NOT NULL DEFAULT 'pending',
    executed_by BIGINT NULL,
    execution_time INT NULL COMMENT 'Execution time in seconds',
    notes TEXT,
    executed_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (test_run_id) REFERENCES test_runs(id) ON DELETE CASCADE,
    FOREIGN KEY (test_case_id) REFERENCES test_cases(id) ON DELETE CASCADE,
    FOREIGN KEY (executed_by) REFERENCES users(id)
);

-- Create step_results table for tracking individual step execution results
CREATE TABLE IF NOT EXISTS step_results (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    test_execution_id BIGINT NOT NULL,
    step_id BIGINT NOT NULL,
    status ENUM('passed', 'failed', 'blocked', 'skipped') NOT NULL,
    actual_result TEXT,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (test_execution_id) REFERENCES test_executions(id) ON DELETE CASCADE,
    FOREIGN KEY (step_id) REFERENCES test_steps(id) ON DELETE CASCADE
);

-- Create defects table for tracking issues found during testing
CREATE TABLE IF NOT EXISTS defects (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    test_execution_id BIGINT NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    severity ENUM('critical', 'high', 'medium', 'low') NOT NULL DEFAULT 'medium',
    status ENUM('open', 'in_progress', 'resolved', 'closed') NOT NULL DEFAULT 'open',
    reported_by BIGINT NOT NULL,
    assigned_to BIGINT NULL,
    external_id VARCHAR(100) NULL COMMENT 'ID in external issue tracking system',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (test_execution_id) REFERENCES test_executions(id) ON DELETE CASCADE,
    FOREIGN KEY (reported_by) REFERENCES users(id),
    FOREIGN KEY (assigned_to) REFERENCES users(id)
);

-- Create defect_attachments table for files attached to defects
CREATE TABLE IF NOT EXISTS defect_attachments (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    defect_id BIGINT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (defect_id) REFERENCES defects(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
); 