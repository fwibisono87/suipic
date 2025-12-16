# Photo Gallery UI Implementation

## Overview
This document describes the implementation of the photo gallery UI for the album detail page with grid/masonry layout, virtual scrolling, lightbox viewer, and responsive design.

## Components

### 1. PhotoGallery.svelte
Main gallery component that displays photos in either grid or masonry layout.

**Features:**
- Virtual scrolling for performance with large photo collections
- Grid layout with responsive columns (2-5 columns based on screen size)
- Masonry layout with CSS columns for a Pinterest-like appearance
- Automatic viewport detection and scroll handling
- Optimized rendering with visible range calculation

**Props:**
- `photos: TPhoto[]` - Array of photos to display
- `layout: 'grid' | 'masonry'` - Layout mode (default: 'grid')

### 2. PhotoThumbnail.svelte
Individual photo thumbnail component with lazy loading.

**Features:**
- Intersection Observer-based lazy loading
- Progressive loading with skeleton placeholder
- Error state handling with broken image indicator
- Hover/active states with photo metadata overlay
- Touch-friendly with active state on mobile
- Displays photo title, filename, stars, and pick/reject state

**Props:**
- `photo: TPhoto` - Photo object to display
- `onClick: () => void` - Click handler
- `lazyLoad: boolean` - Enable lazy loading (default: true)

### 3. Lightbox.svelte
Full-screen photo viewer with navigation.

**Features:**
- Full-screen modal overlay
- Keyboard navigation (Arrow keys, ESC)
- Touch gesture navigation (swipe left/right)
- Image preloading for adjacent photos
- Photo metadata display (title, description, location, date, stars, pick/reject)
- Loading states with spinner
- Error handling with fallback UI
- Mobile-optimized controls and layout
- Photo counter (current/total)

**Props:**
- `photos: TPhoto[]` - Array of all photos
- `currentIndex: number` - Currently displayed photo index
- `isOpen: boolean` - Modal open state
- `onClose: () => void` - Close handler

## Album Detail Page Updates

### Layout Toggle
Added toggle buttons to switch between grid and masonry layouts:
- Grid layout: Uniform square thumbnails in responsive grid
- Masonry layout: Variable height with CSS columns

### Responsive Design
- Mobile: 2 columns (grid), 2 columns (masonry)
- Tablet (sm): 3 columns
- Desktop (md): 3 columns
- Large (lg): 4 columns
- XL: 5 columns

### Photo Count Indicator
Added photo count display in album header showing total number of photos.

## Performance Optimizations

### Virtual Scrolling
- Only renders photos in viewport + 2 viewport heights buffer
- Dramatically reduces DOM nodes for large albums
- Recalculates visible range on scroll and resize
- Only active in grid layout (masonry uses all photos due to CSS columns)

### Lazy Loading
- Images load only when entering viewport
- 100px root margin for smooth user experience
- Fallback for browsers without Intersection Observer

### Image Preloading
- Lightbox preloads adjacent images for faster navigation
- Maintains preload cache to avoid redundant requests
- Clears cache on component destroy

## Keyboard Shortcuts

When lightbox is open:
- `Arrow Left` / `Arrow Right` - Navigate between photos
- `Escape` - Close lightbox

## Touch Gestures

When lightbox is open:
- Swipe left - Next photo
- Swipe right - Previous photo
- Tap outside image - Close lightbox

## API Integration

### Photo Endpoints
- `GET /api/photos/{id}` - Full-size photo
- `GET /api/thumbnails/{id}` - Thumbnail (300x300)

### Configuration
Photos are served through the configured API URL in `frontend/src/lib/config.ts`.

## Styling

### Tailwind Classes
Uses DaisyUI components and custom Tailwind utilities:
- Responsive grid/columns
- Custom shadows and hover effects
- Gradient overlays
- Mobile-first responsive design

### Custom CSS
- Line clamping utilities for text truncation
- Lightbox animations (fade in)
- Button hover effects
- Shadow utilities for depth

## Browser Support
- Modern browsers with ES6 support
- Intersection Observer (with fallback)
- Touch events for mobile
- CSS Grid and CSS Columns

## Future Enhancements
- Infinite scroll for very large albums
- Photo selection/bulk actions
- Drag and drop reordering
- Zoom controls in lightbox
- Download button in lightbox
- Share functionality
