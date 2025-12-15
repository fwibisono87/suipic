# Authentication Implementation Summary

## Files Created

### Models
- **`models/user.go`**: User, Album, and AlbumPermission models with role definitions

### Services
- **`services/database.go`**: Database service with PostgreSQL integration, schema initialization, and user/album CRUD operations
- **`services/auth.go`**: Authentication service with JWT token generation/validation, bcrypt password hashing, user registration/login, and album access control

### Handlers
- **`handlers/auth.go`**: HTTP handlers for register, login, and get current user endpoints

### Middleware
- **`middleware/auth.go`**: Role-based middleware including:
  - `AuthRequired`: JWT token validation
  - `AdminOnly`: Admin-only access
  - `PhotographerOnly`: Photographer/admin access
  - `AlbumAccess`: Album-specific access control with view/edit permissions

### Documentation
- **`AUTH.md`**: Complete authentication system documentation

## Files Modified

### Configuration
- **`config/config.go`**: Added `AdminConfig` struct with admin user settings
- **`.env.example`**: Added admin user environment variables (ADMIN_EMAIL, ADMIN_PASSWORD, ADMIN_USERNAME)

### Main Application
- **`main.go`**: 
  - Initialized database service
  - Initialized auth service with admin seeding
  - Added auth routes (register, login, me)
  - Applied middleware to existing photo routes (PhotographerOnly for upload, AdminOnly for delete)

### Dependencies
- **`go.mod`**: Added required dependencies:
  - `github.com/golang-jwt/jwt/v5` (JWT tokens)
  - `golang.org/x/crypto` (bcrypt password hashing)
  - `github.com/lib/pq` (PostgreSQL driver)

## Key Features Implemented

### 1. JWT Token Authentication
- Token generation with configurable expiry
- Token validation with claims extraction
- Secure signing using HMAC-SHA256

### 2. Password Security
- bcrypt hashing with default cost
- Secure password comparison
- Password hash never exposed in API responses

### 3. Role-Based Access Control (RBAC)
Three user roles:
- **Admin**: Full system access
- **Photographer**: Can upload photos and manage albums
- **User**: Basic access to public/permitted albums

### 4. Album Permissions
- Owner-based access
- Public album support
- Fine-grained permissions (view/edit) per user
- Admin override for all albums

### 5. Admin User Seeding
- Automatic creation of initial admin user on startup
- Configured via environment variables
- Idempotent (checks if user exists before creation)

## Database Schema

### Tables Created (auto-initialized)
1. **users**: User accounts with roles
2. **albums**: Photo albums with ownership
3. **album_permissions**: User-specific album access rights

All tables include proper indexes, foreign keys, and cascading deletes.

## API Endpoints Added

### Authentication Routes
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/me` - Get current user (authenticated)

### Protected Routes (Modified)
- `POST /api/v1/photos/` - Now requires photographer role
- `DELETE /api/v1/photos/:id` - Now requires admin role

## Environment Variables Required

```env
# JWT Configuration
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRY=24h

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=suipic
DB_PASSWORD=password
DB_NAME=suipic
DB_SSLMODE=disable

# Initial Admin User (optional)
ADMIN_EMAIL=admin@suipic.local
ADMIN_PASSWORD=admin123
ADMIN_USERNAME=admin
```

## Usage Example

### Register a New User
```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "photographer@example.com",
    "username": "photographer1",
    "password": "securepass123",
    "role": "photographer"
  }'
```

### Login
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "photographer@example.com",
    "password": "securepass123"
  }'
```

### Access Protected Route
```bash
curl -X POST http://localhost:3000/api/v1/photos/ \
  -H "Authorization: Bearer <your-jwt-token>" \
  -F "file=@photo.jpg"
```

## Security Features

1. **JWT Secret**: Configurable via environment variable
2. **Password Hashing**: bcrypt with appropriate cost factor
3. **Token Expiry**: Configurable duration
4. **SQL Injection Protection**: Parameterized queries
5. **Authorization Headers**: Bearer token standard
6. **Role Validation**: Middleware-based access control

## Integration Notes

- Database schema is automatically created on first run
- Admin user is seeded on startup if credentials provided
- All existing routes continue to work
- Photo upload/delete now have role restrictions
- No breaking changes to existing API contracts
