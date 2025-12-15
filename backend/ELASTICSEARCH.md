# ElasticSearch Integration

## Overview

This application integrates ElasticSearch for advanced photo search capabilities, indexing photo metadata, EXIF data, album information, and comments.

## Features

- **Real-time Indexing**: Photos are automatically indexed when created or updated
- **Bulk Indexing**: Ability to index all photos in an album at once
- **Advanced Search**: Search across multiple fields with filters
- **Comment Indexing**: Photo comments are indexed for full-text search

## Configuration

Set the following environment variables in `.env`:

```bash
ES_ADDRESSES=http://localhost:9200
ES_USERNAME=
ES_PASSWORD=
```

For multiple ElasticSearch nodes, separate addresses with commas:
```bash
ES_ADDRESSES=http://es-node1:9200,http://es-node2:9200
```

## Indexed Fields

Each photo document includes:

- `id`: Photo ID
- `album_id`: Album ID
- `title`: Photo title
- `date_time`: Photo date/time (from EXIF or manual)
- `exif_data`: Complete EXIF metadata as nested object
- `album_title`: Album title
- `album_location`: Album location
- `album_custom_fields`: Album custom fields as nested object
- `comments`: Array of comment texts
- `pick_reject_state`: Photo state (none, pick, reject)
- `stars`: Star rating (0-5)
- `created_at`: Photo creation timestamp
- `updated_at`: Photo last update timestamp

## API Endpoints

### Search Photos

```
GET /api/search
```

**Query Parameters:**
- `q`: Search query (searches title, album title, location, comments, EXIF data)
- `album`: Filter by album ID
- `dateFrom`: Filter by date range start (RFC3339 format)
- `dateTo`: Filter by date range end (RFC3339 format)
- `minStars`: Minimum star rating (0-5)
- `maxStars`: Maximum star rating (0-5)
- `state`: Filter by pick/reject state (none, pick, reject)
- `limit`: Number of results to return (default: 50, max: 1000)
- `offset`: Number of results to skip (default: 0)

**Example:**
```bash
# Search for "wedding" photos with 4+ stars
GET /api/search?q=wedding&minStars=4

# Search photos in album 5 taken in 2024
GET /api/search?album=5&dateFrom=2024-01-01T00:00:00Z&dateTo=2024-12-31T23:59:59Z

# Get all "pick" photos
GET /api/search?state=pick
```

**Response:**
```json
{
  "total": 42,
  "photos": [
    {
      "id": 123,
      "albumId": 5,
      "filename": "abc123.jpg",
      "title": "Beautiful sunset",
      "dateTime": "2024-03-15T18:30:00Z",
      "exifData": {
        "Make": "Canon",
        "Model": "EOS R5",
        "ISO": 100,
        "FocalLength": 24
      },
      "pickRejectState": "pick",
      "stars": 5,
      "createdAt": "2024-03-20T10:00:00Z",
      "updatedAt": "2024-03-20T10:00:00Z"
    }
  ]
}
```

### Bulk Index Album

```
POST /api/albums/:albumId/index
```

Triggers bulk indexing of all photos in an album. Useful for initial indexing or re-indexing after updates.

**Example:**
```bash
POST /api/albums/5/index
```

**Response:**
```json
{
  "message": "album photos indexed successfully"
}
```

## Automatic Indexing

Photos are automatically indexed in the following scenarios:

1. **Photo Creation**: When a new photo is uploaded
2. **Photo Update**: When photo metadata (title, stars, state) is updated
3. **Photo Deletion**: Photo is removed from the index
4. **Comment Creation**: Photo is re-indexed when a comment is added

## Search Behavior

- **Multi-field Search**: The query searches across title, album title, album location, comments, and EXIF data
- **Field Boosting**: Title has 3x weight, album title has 2x weight
- **Date Range**: Uses the `date_time` field from photos
- **Rating Range**: Filters by star rating (0-5)
- **State Filter**: Exact match on pick_reject_state field
- **Default Sorting**: Results are sorted by date_time (descending), then created_at (descending)

## Index Management

### Index Creation

The index is automatically created on service initialization with the appropriate mapping.

### Re-indexing

To re-index all photos in an album:
```bash
POST /api/albums/:albumId/index
```

### Index Name

The default index name is `photos`. This is defined in `services/elasticsearch.go`.

## Error Handling

ElasticSearch failures are logged but do not prevent photo operations from completing. If ElasticSearch is unavailable:

- Photos can still be created, updated, and deleted
- Search endpoint will return errors
- Warnings are logged for indexing failures

## Graceful Degradation

The application can run without ElasticSearch. If connection fails during startup:
- A warning is logged
- The service continues to operate
- Search functionality will be unavailable
- Photo operations work normally
