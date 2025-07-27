# Database Migrations with Atlas

This project uses [Atlas](https://atlasgo.io/) for database schema migration management.

## Setup

Atlas is already installed in this project (`atlas.exe` binary) and configured via `atlas.hcl`.

## Configuration

The Atlas configuration supports three environments:
- `dev`: Development environment (localhost PostgreSQL)
- `test`: Test environment (localhost PostgreSQL test database)
- `production`: Production environment (uses DATABASE_URL environment variable)

## Common Commands

### Check Migration Status
```bash
make migrate-status
```

### Apply Migrations
```bash
make migrate-apply
```

### Create New Migration
```bash
make migrate-diff
# You'll be prompted to enter a migration name
```

### Validate Migrations
```bash
make migrate-validate
```

### Rollback Migrations
```bash
make migrate-down
# You'll be prompted to enter the number of migrations to rollback
```

## Direct Atlas Commands

You can also use Atlas directly:

```bash
# Check status
./atlas.exe migrate status --env dev

# Apply migrations
./atlas.exe migrate apply --env dev

# Create new migration
./atlas.exe migrate diff migration_name --env dev

# Validate migrations
./atlas.exe migrate validate --env dev
```

## Migration Files

- Migration files are stored in the `migrations/` directory
- Each migration has a timestamp prefix and descriptive name
- The `atlas.sum` file contains checksums for migration integrity

## Schema Definition

The current database schema is defined in `schema.sql` and includes:
- Users table with soft delete support
- Proper indexing for performance

## Environment Variables

For production, set the `DATABASE_URL` environment variable:
```
DATABASE_URL=postgres://user:password@host:port/dbname?sslmode=require
```