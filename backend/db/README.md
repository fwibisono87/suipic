# Database Migrations

This directory contains PostgreSQL database migrations for the Suipic application.

## Schema Overview

### Tables

1. **users** - User accounts with role-based access (admin, photographer, client)
   - Fields: username, password_hash, email, friendly_name, role
   
2. **albums** - Photo albums created by photographers
   - Fields: title, date_taken, description, location, custom_fields (JSONB), thumbnail_photo_id, photographer_id
   
3. **photos** - Individual photos within albums
   - Fields: album_id, filename, title, date_time, exif_data (JSONB), pick_reject_state (enum), stars (0-5)
   
4. **album_users** - Junction table for album access control
   - Fields: album_id, user_id
   
5. **comments** - Threaded comments on photos
   - Fields: photo_id, user_id, parent_comment_id, text, created_at

## Running Migrations

### Using the migration tool:

```bash
# Run all migrations
cd backend
go run cmd/migrate/main.go -action=up

# Rollback last migration
go run cmd/migrate/main.go -action=down
```

### Direct SQL execution:

```bash
# Connect to PostgreSQL
psql -U suipic -d suipic

# Run migrations manually
\i db/migrations/000001_create_users_table.up.sql
\i db/migrations/000002_create_albums_table.up.sql
\i db/migrations/000003_create_photos_table.up.sql
\i db/migrations/000004_add_thumbnail_photo_fk.up.sql
\i db/migrations/000005_create_album_users_table.up.sql
\i db/migrations/000006_create_comments_table.up.sql
```

## Migration Files

Migrations are numbered sequentially and come in pairs:
- `*.up.sql` - Applied when migrating forward
- `*.down.sql` - Applied when rolling back

Current migrations:
1. Create users table with role enum
2. Create albums table
3. Create photos table with pick_reject_state enum
4. Add thumbnail photo foreign key to albums
5. Create album_users junction table
6. Create comments table with threading support
