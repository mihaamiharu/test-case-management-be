# Database Migrations

This directory contains the database migration files for the Test Case Management System.

## Migration Order

The migrations should be applied in the following order:

1. `001_create_users_table.sql` - Creates the users table for user management
2. `002_create_projects_table.sql` - Creates the projects table for project management
3. `003_create_project_access.sql` - Creates the project_access table for managing user access to projects
4. `004_create_test_cases_tables.sql` - Creates tables for test suites, test cases, steps, notes, attachments, and tags
5. `005_create_test_execution_tables.sql` - Creates tables for test runs, executions, step results, and defects
6. `006_create_test_environments.sql` - Creates tables for test environments and environment variables
7. `007_create_test_plans.sql` - Creates tables for test plans and test plan items

## Database Schema

### User Management
- `users` - Stores user information including username, email, password hash, and role

### Project Management
- `projects` - Stores project information
- `project_access` - Manages user access levels to projects

### Test Case Management
- `test_suites` - Organizes test cases into logical groups
- `tags` - Provides categorization for test cases
- `test_cases` - Stores test case details
- `test_steps` - Stores Gherkin-style steps for test cases
- `step_notes` - Stores notes attached to test steps
- `step_attachments` - Stores files attached to test steps
- `test_case_tags` - Junction table for test case and tag relationships
- `test_case_history` - Tracks version history of test cases

### Test Execution
- `test_runs` - Tracks test execution sessions
- `test_executions` - Stores individual test case execution results
- `step_results` - Tracks individual step execution results
- `defects` - Tracks issues found during testing
- `defect_attachments` - Stores files attached to defects

### Test Environments
- `environments` - Stores information about test environments
- `environment_variables` - Stores environment-specific variables

### Test Planning
- `test_plans` - Organizes test cases for execution
- `test_plan_items` - Associates test cases with test plans

## Entity Relationships

- A user can own multiple projects
- A project can have multiple test suites
- A test suite can have multiple test cases
- A test case can have multiple steps
- A step can have multiple notes and attachments
- A test case can have multiple tags
- A test run can include multiple test executions
- A test execution is for a single test case
- A test execution can have multiple step results
- A test execution can have multiple defects
- A defect can have multiple attachments
- A project can have multiple environments
- An environment can have multiple variables
- A project can have multiple test plans
- A test plan can include multiple test cases 