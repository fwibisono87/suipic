# Album Management UI Implementation

This document describes the comprehensive album management UI implementation for Suipic.

## Features Implemented

### 1. Album List View (`/albums`)
- **Grid View**: Card-based display with thumbnails, titles, descriptions, dates, locations, and custom fields
- **Table View**: Tabular display with sortable columns
- **View Toggle**: Switch between grid and table views
- **Empty State**: Helpful message and CTA when no albums exist
- **Create Button**: Quick access to create new albums (for photographers/admins only)

### 2. Album Detail View (`/albums/[id]`)
- **Full Album Information**: Display of title, description, location, date taken
- **Custom Fields**: Badge display of all custom fields
- **Assigned Clients**: Display of all clients assigned to the album (for photographers/admins)
- **Photo Gallery**: Grid display of all photos in the album with stars and pick/reject states
- **Edit/Delete Actions**: Quick access buttons (for album owner/admin)
- **Photo Upload**: Direct upload functionality (for album owner/admin)
- **Breadcrumb Navigation**: Easy navigation back to albums list

### 3. Create Album Form (`/albums/new`)
- **Title**: Required text input
- **Description**: Optional textarea
- **Location**: Optional text input
- **Date Taken**: Optional date picker
- **Custom Fields**: Dynamic key-value pairs with add/remove functionality
- **Client Assignment**: Multi-select checkboxes to assign clients
- **Validation**: Client-side validation for required fields
- **Cancel/Submit Actions**: Navigation controls

### 4. Edit Album Form (`/albums/[id]/edit`)
- All create form features plus:
- **Pre-filled Data**: Form loads existing album data
- **Thumbnail Selection**: Visual grid to select album thumbnail from existing photos
- **Update Clients**: Modify assigned client list
- **Permission Check**: Only album owner/admin can access

### 5. Delete Confirmation Modal
- **Warning Message**: Clear explanation of delete consequences
- **Confirmation Dialog**: Prevents accidental deletion
- **Cancel/Confirm Actions**: Safe interaction pattern

## Components Created

### `AlbumForm.svelte`
Reusable form component for both create and edit modes with:
- Form state management
- Dynamic custom fields
- Client multi-select with checkboxes
- Thumbnail selection (edit mode only)
- Loading states
- Error handling

### `ConfirmModal.svelte`
Generic confirmation modal for dangerous actions:
- Customizable title, message, and button text
- Backdrop click to cancel
- Keyboard navigation support

## API Extensions

### Frontend (`frontend/src/lib/api/`)
- **albums.ts**: Added `assignUsers()` and `getUsers()` methods
- **photos.ts**: New module with `listByAlbum()` and `create()` methods

### Backend
- **handlers/album.go**: Added `GetAlbumUsers()` handler
- **main.go**: Added `GET /api/albums/:id/users` route
- **services/album.go**: Updated `AssignUsersToAlbum()` to replace assignments
- **repository/album_user_postgres.go**: Added `DeleteByAlbum()` method
- **repository/interfaces.go**: Updated `AlbumUserRepository` interface

## Type Definitions

### Updated Types (`frontend/src/lib/types/album.ts`)
- Added `customFields` to `TAlbum` type
- Added `customFields` to `TCreateAlbumRequest` and `TUpdateAlbumRequest`
- Added `TPhoto` type for photo data

## Routing Structure

```
/albums                     - Album list (grid/table view)
/albums/new                 - Create new album
/albums/[id]                - Album detail view
/albums/[id]/edit           - Edit album form
```

## Features by User Role

### Admin
- View all albums
- Create new albums
- Edit any album
- Delete any album
- Assign clients to any album
- Upload photos to any album

### Photographer
- View own albums
- View albums where assigned as client
- Create new albums
- Edit own albums
- Delete own albums
- Assign clients to own albums
- Upload photos to own albums

### Client
- View albums where assigned
- View photos in assigned albums
- Cannot create, edit, or delete albums
- Cannot upload photos

## Key Technologies Used

- **SvelteKit**: Frontend framework
- **TanStack Query**: Data fetching and caching
- **DaisyUI**: UI component library
- **Iconify**: Icon library
- **TypeScript**: Type safety
- **Go Fiber**: Backend API framework
- **PostgreSQL**: Database

## Notes

- Album users are replaced (not appended) when updating client assignments
- Thumbnail selection is only available when editing an album with photos
- Custom fields support any key-value pairs
- All forms include loading states and error handling
- API endpoints use proper authentication and authorization
- Breadcrumb navigation helps users understand their location
