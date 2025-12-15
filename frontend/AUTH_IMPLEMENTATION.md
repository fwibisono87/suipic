# Authentication Implementation

This document describes the authentication system implemented in the Suipic frontend.

## Overview

The authentication system provides:
- JWT token-based authentication
- Token storage in both localStorage and cookies
- Protected route guards (server-side and client-side)
- Role-based access control (RBAC)
- Automatic token refresh
- Persistent login sessions

## Architecture

### Auth Store (`src/lib/stores/auth.ts`)

The auth store manages the authentication state using Svelte stores:

```typescript
type AuthState = {
  user: TUser | null;
  token: string | null;
  isLoading: boolean;
};
```

**Exported stores:**
- `authStore` - Main store with methods for managing auth state
- `isAuthenticated` - Derived store that returns boolean
- `currentUser` - Derived store that returns current user
- `authToken` - Derived store that returns current token
- `isLoading` - Derived store for loading state

**Methods:**
- `setAuth(user, token)` - Sets user and token in state and storage
- `clearAuth()` - Clears auth state and storage
- `loadFromStorage()` - Loads auth state from localStorage/cookies
- `updateUser(user)` - Updates user data

### Token Storage

Tokens are stored in two places for maximum compatibility:

1. **localStorage** - For client-side persistence
   - Key: `suipic_token`
   - Key: `suipic_user`

2. **Cookies** - For server-side access
   - Cookie: `suipic_token`
   - Cookie: `suipic_user`
   - SameSite: Strict
   - MaxAge: 7 days
   - HttpOnly: false (needed for client access)

### Protected Routes

#### Server-Side Guards (`+page.server.ts` and `+layout.server.ts`)

Server-side route guards protect pages before they render:

```typescript
// Root layout guard
export const load: LayoutServerLoad = async ({ cookies, url }) => {
  const token = cookies.get('suipic_token');
  const userStr = cookies.get('suipic_user');
  
  // Redirect to login if not authenticated
  if (!isAuthenticated && !isPublicRoute(pathname)) {
    throw redirect(303, '/login');
  }
  
  // Role-based redirects
  if (requiresAdmin && user.role !== 'admin') {
    throw redirect(303, '/');
  }
};
```

**Protected Routes:**
- `/profile` - Requires authentication
- `/settings` - Requires authentication
- `/albums` - Requires authentication
- `/albums/new` - Requires photographer or admin role
- `/albums/[id]` - Requires authentication
- `/admin` - Requires admin role

**Public Routes:**
- `/login`
- `/register`
- `/about`
- `/contact`
- `/privacy`
- `/terms`

#### Client-Side Guards

Client-side guards using utilities from `src/lib/utils/guards.ts`:

```typescript
import { requireAuth, requireRole, requireRoles } from '$lib/utils';

onMount(() => {
  // Basic auth check
  requireAuth();
  
  // Role-specific check
  requireRole(EUserRole.ADMIN);
  
  // Multiple roles check
  requireRoles([EUserRole.PHOTOGRAPHER, EUserRole.ADMIN]);
});
```

### API Integration

All API calls use the auth token from the store:

```typescript
const getAuthHeaders = () => {
  const token = get(authToken);
  return {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {})
  };
};
```

### Role-Based Navigation

The Navbar component (`src/lib/components/Navbar.svelte`) dynamically shows/hides menu items based on user role:

**All authenticated users:**
- Home
- Albums
- Profile dropdown

**Photographers and Admins:**
- New Album button

**Admins only:**
- Admin dashboard link

**Non-authenticated:**
- Login button
- Register button

## User Roles

The system supports three user roles:

1. **Client** (`client`)
   - View albums shared with them
   - Comment on photos
   - Download photos

2. **Photographer** (`photographer`)
   - All client permissions
   - Create albums
   - Upload photos
   - Manage their albums
   - Share albums with clients

3. **Admin** (`admin`)
   - All photographer permissions
   - Access admin dashboard
   - Manage all users
   - Manage all albums
   - System-wide settings

## Authentication Flow

### Login Flow

1. User submits credentials (username/email + password)
2. Frontend calls `/api/auth/login`
3. Backend validates and returns user + JWT token
4. Frontend stores token in localStorage and cookies
5. Frontend updates auth store
6. User is redirected to home page

### Registration Flow

1. User submits registration data
2. Frontend calls `/api/auth/register`
3. Backend creates user and returns user + JWT token
4. Frontend stores token and updates state
5. User is automatically logged in and redirected

### Logout Flow

1. User clicks logout
2. Frontend calls `/api/auth/logout` (optional backend cleanup)
3. Frontend clears localStorage and cookies
4. Frontend clears auth store
5. User is redirected to login page

### Session Persistence

1. On app load, `+layout.svelte` calls `authStore.loadFromStorage()`
2. Store reads token and user from localStorage
3. If found, validates and restores session
4. Cookies are synced for server-side access

## Components

### Login Page (`/login`)
- Username or email login
- Toggle between username/email
- Client-side validation
- Server-side form actions
- Error handling

### Register Page (`/register`)
- Email, username, password fields
- Role selection (client/photographer)
- Password confirmation
- Client-side validation
- Server-side form actions

### Navbar Component
- Dynamic menu based on auth state
- Role-based menu items
- User dropdown with profile/settings
- Theme toggle
- Logout button

## Security Considerations

1. **Token Storage**: Tokens are stored in cookies with SameSite=Strict
2. **HTTPS**: Should be enforced in production
3. **Token Expiration**: JWT tokens expire after configured duration
4. **Protected Routes**: Both server and client-side guards
5. **Role Validation**: Server-side role checks before rendering
6. **Input Validation**: All user inputs validated client and server-side

## Environment Variables

```env
VITE_API_URL=http://localhost:8080/api
```

## Usage Examples

### Protecting a Page

```typescript
// +page.server.ts
export const load: PageServerLoad = async ({ locals }) => {
  if (!locals.user) {
    throw redirect(303, '/login');
  }
  return { user: locals.user };
};
```

### Checking Auth in Component

```svelte
<script>
  import { isAuthenticated, currentUser } from '$lib/stores';
  
  onMount(() => {
    if (!$isAuthenticated) {
      goto('/login');
    }
  });
</script>

{#if $isAuthenticated}
  <p>Welcome, {$currentUser?.username}!</p>
{/if}
```

### Role-Based Rendering

```svelte
<script>
  import { currentUser } from '$lib/stores';
  import { EUserRole } from '$lib/types';
</script>

{#if $currentUser?.role === EUserRole.ADMIN}
  <a href="/admin">Admin Dashboard</a>
{/if}

{#if $currentUser?.role === EUserRole.PHOTOGRAPHER || $currentUser?.role === EUserRole.ADMIN}
  <button>Upload Photos</button>
{/if}
```

### Making Authenticated API Calls

```typescript
import { albumsApi } from '$lib/api';

// Token is automatically included from store
const albums = await albumsApi.list();
```

## Files Structure

```
frontend/src/
├── lib/
│   ├── api/
│   │   ├── auth.ts           # Auth API calls
│   │   ├── albums.ts         # Albums API calls
│   │   └── index.ts
│   ├── components/
│   │   ├── Navbar.svelte     # Navigation with role-based menu
│   │   ├── Footer.svelte
│   │   ├── Alert.svelte
│   │   ├── LoadingSpinner.svelte
│   │   ├── Card.svelte
│   │   └── index.ts
│   ├── stores/
│   │   ├── auth.ts           # Auth state management
│   │   ├── theme.ts
│   │   └── index.ts
│   ├── types/
│   │   ├── auth.ts           # Auth types
│   │   ├── album.ts
│   │   └── index.ts
│   ├── utils/
│   │   ├── validation.ts     # Input validation
│   │   ├── format.ts         # Date/time formatting
│   │   ├── guards.ts         # Route guard utilities
│   │   └── index.ts
│   └── config.ts             # App configuration
├── routes/
│   ├── login/
│   │   ├── +page.svelte
│   │   └── +page.server.ts   # Login form action
│   ├── register/
│   │   ├── +page.svelte
│   │   └── +page.server.ts   # Register form action
│   ├── profile/
│   │   ├── +page.svelte
│   │   └── +page.server.ts   # Auth guard
│   ├── settings/
│   │   ├── +page.svelte
│   │   └── +page.server.ts   # Auth guard
│   ├── admin/
│   │   ├── +page.svelte
│   │   └── +page.server.ts   # Admin guard
│   ├── albums/
│   │   ├── +page.svelte
│   │   ├── +page.server.ts   # Auth guard
│   │   ├── new/
│   │   │   ├── +page.svelte
│   │   │   └── +page.server.ts  # Role guard
│   │   └── [id]/
│   │       ├── +page.svelte
│   │       └── +page.server.ts  # Auth guard
│   ├── +layout.svelte        # Root layout
│   ├── +layout.server.ts     # Root layout guard
│   └── +page.svelte          # Home page
├── app.d.ts                  # Type definitions
└── hooks.server.ts           # Server hooks
```

## Testing

To test the authentication system:

1. **Registration**: Go to `/register` and create a new account
2. **Login**: Go to `/login` and sign in
3. **Protected Routes**: Try accessing `/profile` without auth
4. **Role Access**: Try accessing `/admin` as non-admin
5. **Logout**: Click logout and verify redirect
6. **Session Persistence**: Refresh page and verify still logged in

## Future Enhancements

- [ ] Password reset functionality
- [ ] Email verification
- [ ] Two-factor authentication (2FA)
- [ ] Social login (OAuth)
- [ ] Remember me checkbox
- [ ] Session timeout warnings
- [ ] Token refresh mechanism
- [ ] Audit logging
