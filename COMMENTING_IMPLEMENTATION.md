# Commenting System Implementation

## Overview
Fully implemented commenting system with threaded replies, real-time preview, user attribution badges, timestamps, and auto-refresh functionality.

## Backend Changes

### 1. Updated Comment Service (`backend/services/comment.go`)
- Added `userRepo` dependency to fetch user information
- Created `CommentWithUser` struct that includes user data
- Created `ThreadedComment` struct with nested replies and user information
- Updated `GetThreadedComments` to include user data for each comment
- Added `GetCommentWithUser` method to fetch a single comment with user data
- Modified `NewCommentService` to accept `UserRepository` parameter

### 2. Updated Main (`backend/main.go`)
- Updated `CommentService` initialization to pass both `commentRepo` and `userRepo`

### 3. Updated Photo Handler (`backend/handlers/photo.go`)
- Modified `CreateComment` endpoint to return comment with user data
- `GetComments` endpoint already returns threaded comments with user data

## Frontend Changes

### 1. Types (`frontend/src/lib/types/index.ts`)
- Added `TComment` type with user information and optional replies array
- Includes nested user object with all user fields (id, username, email, friendlyName, role)

### 2. Comments API (`frontend/src/lib/api/comments.ts`)
- Created `commentsApi` with two methods:
  - `getByPhoto(photoId)`: Fetches threaded comments for a photo
  - `create(photoId, text, parentCommentId?)`: Creates a new comment or reply

### 3. CommentSection Component (`frontend/src/lib/components/CommentSection.svelte`)

#### Features:
- **Threaded Display**: Shows comments with nested replies (one level deep)
- **User Attribution**: Displays badges for:
  - Photographer (primary badge)
  - Client (secondary badge)
  - Admin (accent badge)
- **Timestamps**: Uses relative time formatting (e.g., "2 minutes ago", "1 hour ago")
- **Comment Form**:
  - Textarea with character count
  - Real-time preview toggle
  - Reply indicator showing who you're replying to
  - Keyboard shortcuts (Ctrl/Cmd+Enter to submit)
- **Auto-refresh**: Polls for new comments every 10 seconds
  - Shows subtle loading spinner during refresh
  - Silent refresh doesn't interrupt user experience
- **Avatar Placeholders**: Shows first letter of user's name in colored circle
- **Empty States**: Friendly message when no comments exist
- **Error Handling**: Displays error messages with dismissible alerts

#### Props:
- `photo`: TPhoto object
- `photographerId`: Number to determine photographer badge

### 4. Updated Components

#### PhotoInteractionPanel (`frontend/src/lib/components/PhotoInteractionPanel.svelte`)
- Added `photographerId` prop
- Integrated CommentSection below EXIF data
- Added divider for visual separation

#### Lightbox (`frontend/src/lib/components/Lightbox.svelte`)
- Added `photographerId` prop
- Passes photographer ID to PhotoInteractionPanel

#### PhotoGallery (`frontend/src/lib/components/PhotoGallery.svelte`)
- Added `photographerId` prop
- Passes photographer ID to Lightbox

#### Album Page (`frontend/src/routes/albums/[id]/+page.svelte`)
- Passes `photographerId` from album data to PhotoGallery
- Added missing `isUploading` state variable

### 5. Format Utilities (`frontend/src/lib/utils/format.ts`)
- Enhanced `formatRelativeTime` to properly handle singular/plural:
  - "1 minute ago" vs "2 minutes ago"
  - "1 hour ago" vs "2 hours ago"
  - "1 day ago" vs "2 days ago"

## Key Features

### 1. Threaded Comments
- Top-level comments displayed in chronological order
- Replies nested under parent comments with visual indent
- Left border on reply section for clear hierarchy

### 2. Real-time Preview
- Toggle button to show/hide preview
- Preview displays exactly as comment will appear
- Preserves whitespace and line breaks
- Dashed border to distinguish from actual comments

### 3. User Badges
- **Photographer Badge**: Primary styled, shown for album owner
- **Client Badge**: Secondary styled, shown for client users
- **Admin Badge**: Accent styled, shown for admin users
- Automatically determined by comparing user ID with photographer ID and user role

### 4. Reply System
- Reply button on each comment and nested reply
- Visual indicator shows who you're replying to
- Cancel button to clear reply context
- Placeholder text updates based on reply context
- Replies are threaded under the parent comment

### 5. Auto-refresh & Polling
- Polls API every 10 seconds for new comments
- Uses silent refresh to avoid disrupting user
- Shows subtle spinner during refresh
- Automatically updates comment count
- Preserves user's input and UI state during refresh

### 6. Responsive Design
- Mobile-friendly avatar sizes (smaller on replies)
- Responsive badge sizes
- Proper text wrapping and truncation
- Smooth animations for comment appearance

### 7. Accessibility
- Proper semantic HTML structure
- ARIA labels where needed
- Keyboard navigation support
- Focus management on reply action

## API Endpoints Used

- `GET /api/photos/:id/comments` - Fetch threaded comments
- `POST /api/photos/:id/comments` - Create new comment
  - Body: `{ text: string, parentCommentId?: number }`

## Comment Count
- Displays total count including replies
- Updates automatically with polling
- Shows singular/plural correctly ("1 comment" vs "2 comments")

## Animations
- Fade-in animation for new comments
- Slide-down animation for reply threads
- Smooth transitions for all interactions

## Future Enhancements (Not Implemented)
- Edit comment functionality
- Delete comment functionality
- Markdown support in comments
- @mentions and notifications
- Comment reactions/likes
- Sorting options (newest/oldest)
- Load more pagination for many comments
