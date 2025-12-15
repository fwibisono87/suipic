# Authentication Implementation

## Overview

This authentication system provides JWT-based authentication with role-based access control (RBAC) for the Suipic photo manager application.

## Features

- **JWT Token Generation and Validation**: Secure token-based authentication
- **Password Hashing**: bcrypt for secure password storage
- **Role-Based Access Control**: Three user roles - admin, photographer, and user
- **Album Permissions**: Fine-grained access control for albums
- **Admin User Seeding**: Automatic creation of initial admin user on startup

## User Roles

### Admin (`admin`)
- Full system access
- Can access all albums and photos
- Can delete any photos
- Cannot be registered via API (only via environment variables)

### Photographer (`photographer`)
- Can upload photos
- Can manage their own albums
- Can view albums they have permission to access

### User (`user`)
- Default role for new registrations
- Can view public albums
- Can view albums they have explicit permission to access

## API Endpoints

### Authentication

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "role": "user"  // Optional: "user" or "photographer", defaults to "user"
}
```

Response:
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Get Current User
```http
GET /api/v1/auth/me
Authorization: Bearer <token>
```

Response:
```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "username",
  "role": "user"
}
```

## Middleware

### AuthRequired
Validates JWT token and extracts user information.

Usage in routes:
```go
app.Get("/protected", middleware.AuthRequired(authService), handler)
```

Sets the following context locals:
- `user_id` (int64)
- `user_email` (string)
- `user_username` (string)
- `user_role` (models.UserRole)

### AdminOnly
Requires authenticated user to have admin role.

Usage:
```go
app.Delete("/admin-only", middleware.AdminOnly(authService), handler)
```

### PhotographerOnly
Requires authenticated user to have photographer or admin role.

Usage:
```go
app.Post("/photographer-only", middleware.PhotographerOnly(authService), handler)
```

### AlbumAccess
Validates user access to specific album with optional edit permission requirement.

Usage:
```go
// View access
app.Get("/albums/:album_id", middleware.AlbumAccess(authService, false), handler)

// Edit access
app.Put("/albums/:album_id", middleware.AlbumAccess(authService, true), handler)
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Albums Table
```sql
CREATE TABLE albums (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Album Permissions Table
```sql
CREATE TABLE album_permissions (
    id SERIAL PRIMARY KEY,
    album_id INTEGER NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    can_view BOOLEAN NOT NULL DEFAULT true,
    can_edit BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(album_id, user_id)
);
```

## Environment Variables

Add the following to your `.env` file:

```env
# JWT Configuration
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRY=24h

# Initial Admin User
ADMIN_EMAIL=admin@suipic.local
ADMIN_PASSWORD=admin123
ADMIN_USERNAME=admin
```

## Admin User Seeding

On application startup, if the admin credentials are provided in environment variables, the system will automatically create an admin user if one doesn't already exist with the specified email.

To disable admin seeding, simply leave the admin environment variables empty.

## Security Considerations

1. **JWT Secret**: Always use a strong, random secret in production
2. **Password Requirements**: Implement password strength requirements in production
3. **Token Expiry**: Adjust token expiry based on security requirements
4. **HTTPS**: Always use HTTPS in production to protect tokens in transit
5. **Rate Limiting**: Consider implementing rate limiting on auth endpoints

## Example Usage

### Protected Route Example
```go
photos := v1.Group("/photos")
photos.Post("/", middleware.PhotographerOnly(authService), photoHandler.UploadPhoto)
photos.Get("/:id", photoHandler.DownloadPhoto)
photos.Delete("/:id", middleware.AdminOnly(authService), photoHandler.DeletePhoto)
```

### Accessing User Info in Handler
```go
func (h *Handler) MyHandler(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int64)
    userRole := c.Locals("user_role").(models.UserRole)
    
    // Use user information
    return c.JSON(fiber.Map{
        "user_id": userID,
        "role": userRole,
    })
}
```

## Album Access Control Logic

Album access is determined by the following rules:

1. **Admin users**: Full access to all albums
2. **Album owner**: Full access to their own albums
3. **Public albums**: All authenticated users can view (but not edit)
4. **Explicit permissions**: Users with album_permissions records have access based on their `can_view` and `can_edit` flags

The `CanAccessAlbum` method in AuthService handles this logic automatically.
