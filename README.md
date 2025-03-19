# Test Case Management System

A comprehensive test case management system built with Go/Gin, MySQL, and Svelte.

## Features

- **User Management**: Authentication, authorization, and role-based access control
- **Project Management**: Create, organize, and manage testing projects
- **Project Access Control**: Grant specific users access to view or edit projects
- **Test Case Management**: Create, read, update, delete test cases
- **Test Execution**: Run tests and record results
- **Reporting**: Generate reports on test coverage and visualize results

## Backend Technology Stack

- **Go/Gin**: Fast and lightweight web framework
- **MySQL**: Relational database for data storage
- **JWT**: JSON Web Tokens for authentication
- **Repository Pattern**: Clean separation of concerns

## Getting Started

### Prerequisites

- Go 1.16+
- MySQL 5.7+

### Database Setup

1. Create a MySQL database:

```sql
CREATE DATABASE test_case_manager;
```

2. Run the migration scripts:

```bash
mysql -u root -p test_case_manager < migrations/001_create_users_table.sql
mysql -u root -p test_case_manager < migrations/002_create_projects_table.sql
mysql -u root -p test_case_manager < migrations/003_create_project_access_table.sql
```

### Configuration

1. Copy the `.env.example` file to `.env` and update the values:

```bash
cp .env.example .env
```

2. Update the database credentials and other settings in the `.env` file.

### Running the Application

1. Build and run the application:

```bash
go build -o server ./cmd/api
./server
```

2. The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token
- `GET /api/v1/auth/me` - Get current user info (requires authentication)

### Projects

- `GET /api/v1/projects` - List all projects the user has access to
- `POST /api/v1/projects` - Create a new project
- `GET /api/v1/projects/{id}` - Get a specific project
- `PUT /api/v1/projects/{id}` - Update a project
- `DELETE /api/v1/projects/{id}` - Delete a project

### Project Access

- `GET /api/v1/projects/{id}/access` - List all users with access to a project
- `POST /api/v1/projects/{id}/access` - Grant a user access to a project
- `PUT /api/v1/projects/{id}/access/{accessId}` - Update a user's access level
- `DELETE /api/v1/projects/{id}/access/{accessId}` - Revoke a user's access

## Access Control System

The system implements a granular access control mechanism:

1. **Project Owners**: The creator of a project is automatically its owner and has full control over it.
2. **Access Levels**:
   - **View**: Users with view access can see the project but not modify it.
   - **Edit**: Users with edit access can both view and modify the project.
3. **Access Management**: Project owners can grant, update, or revoke access for other users.
4. **Admin Override**: Users with the admin role can access and modify all projects.

## License

MIT 