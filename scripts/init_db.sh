#!/bin/bash

# Load environment variables
source .env

# Create database if it doesn't exist
echo "Creating database if it doesn't exist..."
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;"

# Apply schema
echo "Applying database schema..."
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD $DB_NAME < scripts/init_db.sql

echo "Database initialization complete!" 