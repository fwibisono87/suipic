# Photo Upload Feature

## Overview

The photo upload feature provides a comprehensive UI for uploading multiple photos to albums with drag-and-drop support, concurrent uploads, preview generation, and EXIF metadata display.

## Components

### PhotoUploadModal

The main upload component that provides:

- **Drag-and-drop zone**: Drop files anywhere in the modal
- **Multiple file selection**: Browse and select multiple files at once
- **Client-side preview generation**: Thumbnails are generated locally before upload
- **EXIF extraction and display**: Camera settings, lens info, dimensions, and capture date
- **Concurrent uploads**: Configurable concurrent upload limit (default: 3-5)
- **Progress indicators**: Per-file upload status and progress
- **Error handling**: Retry failed uploads individually
- **File validation**: Only accepts supported image formats (JPEG, PNG, WebP, GIF)

#### Usage

```svelte
<script>
  import { PhotoUploadModal } from '$lib/components';

  let showUploadModal = false;
  let albumId = 1;

  function handleUploadComplete(event) {
    const photos = event.detail;
    console.log('Uploaded:', photos);
  }
</script>

<PhotoUploadModal
  isOpen={showUploadModal}
  {albumId}
  on:close={() => (showUploadModal = false)}
  on:uploadComplete={handleUploadComplete}
/>
```

#### Props

- `isOpen`: boolean - Controls modal visibility
- `albumId`: number - The album ID to upload photos to

#### Events

- `close`: Fired when the modal is closed
- `uploadComplete`: Fired when all uploads complete successfully, passes array of uploaded photos

### PhotoUploadButton

A simplified button component for triggering file selection:

```svelte
<script>
  import { PhotoUploadButton } from '$lib/components';

  function handleFilesSelected(event) {
    const files = event.detail;
    // Handle files
  }
</script>

<PhotoUploadButton
  label="Upload Photos"
  multiple={true}
  on:filesSelected={handleFilesSelected}
/>
```

#### Props

- `accept`: string - File types to accept (default: image types)
- `multiple`: boolean - Allow multiple file selection (default: true)
- `disabled`: boolean - Disable the button (default: false)
- `size`: 'xs' | 'sm' | 'md' | 'lg' - Button size (default: 'md')
- `variant`: 'primary' | 'secondary' | 'ghost' | 'outline' - Button style (default: 'primary')
- `label`: string - Button text (default: 'Upload Photos')

## Utilities

### EXIF Extraction (`lib/utils/exif.ts`)

Extracts and formats EXIF metadata from image files:

```typescript
import { extractExifData } from '$lib/utils/exif';

const exif = await extractExifData(file);
// Returns: { make, model, lens, focalLength, aperture, shutterSpeed, iso, dateTime, dimensions, latitude, longitude }
```

### Image Preview Generation (`lib/utils/image.ts`)

Generates optimized preview thumbnails:

```typescript
import { generateImagePreview } from '$lib/utils/image';

const preview = await generateImagePreview(file, 300, 0.8);
// Returns: Data URL of resized image
```

Other utilities:
- `isValidImageType(file, acceptedTypes)`: Validate file type
- `formatFileSize(bytes)`: Format bytes to human-readable size
- `getImageDimensions(file)`: Get original image dimensions

## API Integration

### Photo API (`lib/api/photos.ts`)

```typescript
import { photosApi } from '$lib/api';

// Upload single photo
const photo = await photosApi.create(albumId, file);

// Upload multiple photos (sequential with error handling)
const photos = await photosApi.createBatch(albumId, files);
```

## Features

### 1. Drag-and-Drop Support

- Drop zone highlights when dragging files over it
- Supports dragging multiple files at once
- Visual feedback during drag operations

### 2. Multiple File Selection

- Standard file browser allows selecting multiple files
- Add more files after initial selection
- Remove individual files before upload

### 3. Client-Side Preview Generation

- Thumbnails generated locally (300x300 max)
- Optimized JPEG encoding (80% quality)
- Fallback to blob URLs if canvas fails

### 4. EXIF Metadata Display

Displays before upload:
- Camera make and model
- Lens information
- Focal length
- Aperture (f-number)
- Shutter speed
- ISO
- Capture date/time
- Image dimensions

### 5. Concurrent Upload Management

- Configurable concurrent upload limit (default: 3)
- Queue system for managing uploads
- Upload status tracking per file:
  - `pending`: Waiting to upload
  - `uploading`: Currently uploading
  - `complete`: Successfully uploaded
  - `error`: Upload failed

### 6. Progress Indicators

- Overall statistics (total, pending, uploading, complete, failed)
- Per-file status badges
- Progress bars for active uploads

### 7. Error Handling

- Per-file error messages
- Retry failed uploads individually
- Cancel confirmation if uploads in progress

## Configuration

### Concurrent Upload Limit

Modify `MAX_CONCURRENT_UPLOADS` in `PhotoUploadModal.svelte`:

```typescript
const MAX_CONCURRENT_UPLOADS = 3; // Change to 5 for more concurrent uploads
```

### Accepted File Types

Modify `ACCEPTED_TYPES` in `PhotoUploadModal.svelte`:

```typescript
const ACCEPTED_TYPES = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif'];
```

### Preview Size and Quality

When calling `generateImagePreview`:

```typescript
const preview = await generateImagePreview(
  file,
  300,  // Max dimension in pixels
  0.8   // JPEG quality (0-1)
);
```

## Dependencies

- **exifr**: EXIF metadata extraction library
  - Lightweight and fast
  - Supports all common EXIF fields
  - GPS coordinate extraction

## Browser Support

- Modern browsers with:
  - File API support
  - Canvas API support
  - FileReader API support
  - Drag and Drop API support

## Performance Considerations

1. **Preview Generation**: Async with throttling to avoid blocking UI
2. **EXIF Extraction**: Async and optional (continues if fails)
3. **Concurrent Uploads**: Limited to prevent overwhelming the server
4. **Memory Management**: Blob URLs are cleaned up after use

## Security

- File type validation on client and server
- File size limits enforced by backend
- Authorization checks on upload endpoint
- EXIF stripping option (handled by backend)

## Future Enhancements

- Bulk EXIF editing before upload
- Image rotation/cropping before upload
- Upload pause/resume functionality
- Folder/batch organization
- Duplicate detection
- Automatic album date detection from EXIF
