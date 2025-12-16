# Photo Upload - Quick Start Guide

## Installation

Add the exifr dependency:

```bash
pnpm install exifr@^7.1.3
# or
npm install exifr@^7.1.3
```

## Basic Usage

### 1. Using PhotoUploadModal in Your Page

```svelte
<script lang="ts">
  import { PhotoUploadModal } from '$lib/components';
  
  let showUploadModal = false;
  let albumId = 1; // Your album ID
  
  function handleUploadComplete(event) {
    const photos = event.detail;
    console.log('Uploaded photos:', photos);
    // Refresh your photo list or update UI
  }
</script>

<button on:click={() => showUploadModal = true}>
  Upload Photos
</button>

<PhotoUploadModal
  isOpen={showUploadModal}
  {albumId}
  on:close={() => showUploadModal = false}
  on:uploadComplete={handleUploadComplete}
/>
```

### 2. Using PhotoUploadButton (Without Modal)

```svelte
<script lang="ts">
  import { PhotoUploadButton } from '$lib/components';
  
  function handleFiles(event) {
    const files = event.detail;
    // Handle files yourself (e.g., custom upload logic)
    files.forEach(file => {
      console.log(file.name, file.size);
    });
  }
</script>

<PhotoUploadButton
  label="Select Photos"
  on:filesSelected={handleFiles}
/>
```

## Component Props

### PhotoUploadModal

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `isOpen` | boolean | Yes | Controls modal visibility |
| `albumId` | number | Yes | Album ID for uploads |

### PhotoUploadButton

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `accept` | string | image types | Accepted file types |
| `multiple` | boolean | `true` | Allow multiple files |
| `disabled` | boolean | `false` | Disable button |
| `size` | string | `'md'` | Button size (xs/sm/md/lg) |
| `variant` | string | `'primary'` | Button variant |
| `label` | string | 'Upload Photos' | Button text |

## Events

### PhotoUploadModal

- `close`: Fired when modal closes
- `uploadComplete`: Fired on successful upload (payload: `TPhoto[]`)

### PhotoUploadButton

- `filesSelected`: Fired when files selected (payload: `File[]`)

## Utility Functions

### EXIF Extraction

```typescript
import { extractExifData } from '$lib/utils/exif';

const exif = await extractExifData(file);
console.log(exif.make, exif.model); // "Canon", "EOS R5"
```

### Image Preview

```typescript
import { generateImagePreview } from '$lib/utils/image';

const dataUrl = await generateImagePreview(file, 300, 0.8);
// Use dataUrl in <img src={dataUrl} />
```

## Configuration

### Change Concurrent Upload Limit

Edit `src/lib/components/PhotoUploadModal.svelte`:

```typescript
const MAX_CONCURRENT_UPLOADS = 5; // Default is 3
```

### Add Supported File Types

Edit `src/lib/components/PhotoUploadModal.svelte`:

```typescript
const ACCEPTED_TYPES = [
  'image/jpeg',
  'image/jpg',
  'image/png',
  'image/webp',
  'image/gif',
  'image/heic' // Add HEIC
];
```

## Testing

Visit `/upload-demo` in your running application to test the components.

## Troubleshooting

**EXIF data not showing?**
- Some images don't have EXIF data
- Check browser console for errors
- Ensure exifr is installed

**Uploads failing?**
- Check authentication token
- Verify albumId is valid
- Check network tab for API errors
- Ensure backend is running

**Previews not generating?**
- Check browser console for Canvas API errors
- Ensure file is a valid image
- Try with JPEG images first

## Next Steps

1. Integrate into your album pages (see `src/routes/albums/[id]/+page.svelte` for reference)
2. Customize styling to match your theme
3. Add custom upload logic if needed
4. Test with various image formats and sizes

For detailed documentation, see `PHOTO_UPLOAD.md`.
