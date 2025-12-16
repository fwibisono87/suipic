# Photo Interaction Components

This document describes the new photo interaction UI components added to the Suipic frontend.

## Components Overview

### 1. StarRating Component (`StarRating.svelte`)

Interactive 5-star rating component with hover effects.

**Props:**
- `rating: number` (default: 0) - Current rating (0-5)
- `maxStars: number` (default: 5) - Maximum number of stars
- `size: 'sm' | 'md' | 'lg'` (default: 'md') - Component size
- `readonly: boolean` (default: false) - Whether the rating is read-only
- `showClear: boolean` (default: false) - Show a clear button when rating > 0

**Events:**
- `change: number` - Fired when rating changes

**Usage:**
```svelte
<StarRating 
  rating={3} 
  on:change={(e) => console.log('New rating:', e.detail)} 
  showClear={true}
/>
```

### 2. PickRejectButtons Component (`PickRejectButtons.svelte`)

Color-coded toggle buttons for photo selection workflow (Pick/Unflagged/Reject).

**Props:**
- `state: 'none' | 'pick' | 'reject'` (default: 'none') - Current state
- `disabled: boolean` (default: false) - Whether buttons are disabled
- `size: 'sm' | 'md'` (default: 'md') - Button size

**Events:**
- `change: 'none' | 'pick' | 'reject'` - Fired when state changes

**Usage:**
```svelte
<PickRejectButtons 
  state="pick" 
  on:change={(e) => console.log('New state:', e.detail)}
/>
```

**Color Scheme:**
- Pick: Green (success)
- Unflagged: Gray (ghost)
- Reject: Red (error)

### 3. PhotoMetadata Component (`PhotoMetadata.svelte`)

Display photo information including filename, title, dates, and EXIF data with inline title editing.

**Props:**
- `photo: TPhoto` - Photo object to display
- `editable: boolean` (default: false) - Whether title can be edited inline

**Events:**
- `updateTitle: string | null` - Fired when title is edited

**Usage:**
```svelte
<PhotoMetadata 
  photo={photoObject} 
  editable={true}
  on:updateTitle={(e) => saveTitle(e.detail)}
/>
```

**Features:**
- Click-to-edit title functionality
- EXIF data parsing and display (camera, lens, aperture, shutter speed, ISO, dimensions)
- Formatted date/time display
- Responsive design

### 4. PhotoInteractionPanel Component (`PhotoInteractionPanel.svelte`)

All-in-one panel combining pick/reject buttons, star rating, and metadata display with API integration.

**Props:**
- `photo: TPhoto` - Photo object to interact with

**Events:**
- `update: TPhoto` - Fired when photo is updated via API

**Usage:**
```svelte
<PhotoInteractionPanel 
  photo={photoObject} 
  on:update={(e) => handlePhotoUpdate(e.detail)}
/>
```

**Features:**
- Integrated pick/reject/unflagged toggle
- Star rating (0-5 stars)
- Inline title editing with save/cancel
- Metadata display with EXIF parsing
- Automatic API updates using `photosApi.update()`
- Loading states and error handling

## Lightbox Integration

The `Lightbox` component has been enhanced to include the photo interaction panel:

**New Features:**
- Press `I` key to toggle the info panel
- Info panel slides in from the right on desktop
- Full-screen panel on mobile
- Real-time photo updates propagated to gallery

**Keyboard Shortcuts:**
- `I` - Toggle info panel
- `←` - Previous photo
- `→` - Next photo
- `ESC` - Close panel or lightbox

## API Integration

### photos.ts API Module

New `update` method added:

```typescript
photosApi.update(
  photoId: number, 
  updates: Partial<{ 
    title: string | null; 
    pickRejectState: string | null; 
    stars: number 
  }>
): Promise<TPhoto>
```

**Backend Endpoint:** `PUT /photos/:id`

**Request Body:**
```json
{
  "title": "Updated Title",
  "pickRejectState": "pick",
  "stars": 4
}
```

## Usage Examples

### Standalone Components

```svelte
<script>
  import { StarRating, PickRejectButtons, PhotoMetadata } from '$lib/components';
  
  let rating = 3;
  let state = 'none';
</script>

<StarRating bind:rating on:change={(e) => updateRating(e.detail)} />
<PickRejectButtons bind:state on:change={(e) => updateState(e.detail)} />
<PhotoMetadata photo={photo} editable={true} on:updateTitle={saveTitle} />
```

### Complete Interaction Panel

```svelte
<script>
  import { PhotoInteractionPanel } from '$lib/components';
  
  function handleUpdate(event) {
    const updatedPhoto = event.detail;
    // Photo is already updated via API
    // Optionally refresh gallery or update local state
  }
</script>

<PhotoInteractionPanel photo={currentPhoto} on:update={handleUpdate} />
```

### In Lightbox

```svelte
<script>
  import { Lightbox } from '$lib/components';
  
  function handlePhotoUpdate(photo) {
    // Update photo in photos array
    photos[currentIndex] = photo;
  }
</script>

<Lightbox 
  photos={photos} 
  bind:currentIndex 
  bind:isOpen 
  onClose={closeLightbox}
  onPhotoUpdate={handlePhotoUpdate}
/>
```

## Demo Page

Visit `/components-demo` to see all components in action with interactive examples.

## Styling

Components use DaisyUI classes for consistent theming:
- `btn-success` for Pick state
- `btn-error` for Reject state
- `btn-ghost` for Unflagged state
- `text-warning` for star ratings
- Responsive breakpoints for mobile/desktop views

## EXIF Data Parsing

The components automatically parse and display EXIF data from photos:

- Camera make and model
- Lens model
- Focal length
- Aperture (f-number)
- Shutter speed (exposure time)
- ISO sensitivity
- Image dimensions

EXIF values are formatted for readability (e.g., "f/2.8", "1/250s", "ISO 100").

## Accessibility

- All interactive elements use semantic HTML
- Keyboard navigation support
- ARIA labels for screen readers
- Focus management for inline editing
- Disabled states properly indicated
