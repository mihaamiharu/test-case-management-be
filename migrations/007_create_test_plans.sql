-- Create test_plans table for organizing test cases for execution
CREATE TABLE IF NOT EXISTS test_plans (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status ENUM('draft', 'active', 'archived') NOT NULL DEFAULT 'draft',
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id),
    UNIQUE KEY unique_plan_name_per_project (project_id, name)
);

-- Create test_plan_items table for associating test cases with test plans
CREATE TABLE IF NOT EXISTS test_plan_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    test_plan_id BIGINT NOT NULL,
    test_case_id BIGINT NOT NULL,
    priority INT NOT NULL DEFAULT 0 COMMENT 'Execution order priority',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (test_plan_id) REFERENCES test_plans(id) ON DELETE CASCADE,
    FOREIGN KEY (test_case_id) REFERENCES test_cases(id) ON DELETE CASCADE,
    UNIQUE KEY unique_test_case_per_plan (test_plan_id, test_case_id)
);

-- Add test_plan_id to test_runs table
ALTER TABLE test_runs
ADD COLUMN test_plan_id BIGINT NULL AFTER project_id,
ADD CONSTRAINT fk_test_run_plan FOREIGN KEY (test_plan_id) REFERENCES test_plans(id) ON DELETE SET NULL; 