#!/bin/bash

# Load environment variables from .env file
if [ ! -f .env ]; then
  echo "❌ .env file not found. Please create one with DB_NAME, DB_USER, DB_PASSWORD, and DB_SUPERUSER."
  exit 1
fi

export $(grep -v '^#' .env | xargs)

# Validate required variables
: "${DB_NAME:?Missing DB_NAME in .env}"
: "${DB_USER:?Missing DB_USER in .env}"
: "${DB_PASSWORD:?Missing DB_PASSWORD in .env}"

DB_SUPERUSER=${DB_SUPERUSER:-postgres}

echo "🔧 Setting up PostgreSQL database '${DB_NAME}' and user '${DB_USER}'..."

# Run SQL commands as superuser
PGPASSWORD=${PGPASSWORD:-""} psql -U "$DB_SUPERUSER" -d postgres <<EOF
-- Create user with CREATEDB privilege if it doesn't exist
DO \$\$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles WHERE rolname = '${DB_USER}'
   ) THEN
      CREATE ROLE ${DB_USER} WITH LOGIN PASSWORD '${DB_PASS}' CREATEDB;
   END IF;
END
\$\$;
EOF

# Now create the database in a separate command
PGPASSWORD=${PGPASSWORD:-""} psql -U "$DB_SUPERUSER" -d postgres -tAc \
"SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}'" | grep -q 1 || \
PGPASSWORD=${PGPASSWORD:-""} createdb -U "$DB_SUPERUSER" -O "$DB_USER" "$DB_NAME"

# Connect as new user and grant privileges
PGPASSWORD=$DB_PASS psql -U "$DB_USER" -d "$DB_NAME" <<EOF
GRANT ALL ON SCHEMA public TO ${DB_USER};
GRANT ALL ON ALL TABLES IN SCHEMA public TO ${DB_USER};
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO ${DB_USER};
GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO ${DB_USER};

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ${DB_USER};
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO ${DB_USER};
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO ${DB_USER};
EOF

echo "✅ Setup complete! '${DB_NAME}' is ready for use with user '${DB_USER}'."
