# Search API Examples

This document provides examples of how to use the photo search API.

## Basic Search

### Search by text query
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?q=sunset"
```

Searches across:
- Photo titles
- Album titles
- Album locations
- Photo comments
- EXIF data fields

## Filtering

### Filter by Album
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?album=5"
```

### Filter by Date Range
```bash
# Photos taken in 2024
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?dateFrom=2024-01-01T00:00:00Z&dateTo=2024-12-31T23:59:59Z"
```

### Filter by Star Rating
```bash
# Photos with 4 or 5 stars
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?minStars=4"

# Photos with exactly 3 stars
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?minStars=3&maxStars=3"
```

### Filter by Pick/Reject State
```bash
# Get all "pick" photos
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?state=pick"

# Get all "reject" photos
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?state=reject"

# Get photos with no state set
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?state=none"
```

## Combined Filters

### Wedding photos with 5 stars
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?q=wedding&minStars=5"
```

### Picked photos in specific album
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?album=5&state=pick"
```

### Photos by camera model in date range
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?q=Canon%20EOS%20R5&dateFrom=2024-01-01T00:00:00Z&dateTo=2024-12-31T23:59:59Z"
```

## Pagination

### Limit results
```bash
# Get first 20 results
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?q=sunset&limit=20"
```

### Offset for pagination
```bash
# Get next page (results 20-40)
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?q=sunset&limit=20&offset=20"
```

## Bulk Indexing

### Index all photos in an album
```bash
curl -X POST -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/albums/5/index"
```

This is useful for:
- Initial indexing when setting up ElasticSearch
- Re-indexing after bulk updates
- Fixing any indexing issues

## Response Format

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
        "FocalLength": 24,
        "FNumber": 2.8,
        "ExposureTime": "1/1000"
      },
      "pickRejectState": "pick",
      "stars": 5,
      "createdAt": "2024-03-20T10:00:00Z",
      "updatedAt": "2024-03-20T10:00:00Z"
    }
  ]
}
```

## Advanced Query Examples

### Search by location
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?q=Paris"
```

### Search by lens model
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?q=24-70mm"
```

### Search by ISO setting
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?q=ISO%203200"
```

### Search by artist/photographer in EXIF
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?q=John%20Doe"
```

### Find all photos with comments containing specific text
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:3000/api/search?q=needs%20retouching"
```

## Tips

1. **URL Encoding**: Remember to URL encode query parameters (spaces become `%20`, etc.)
2. **Date Format**: Use RFC3339 format for dates: `YYYY-MM-DDTHH:MM:SSZ`
3. **Star Range**: Valid values are 0-5 for both minStars and maxStars
4. **State Values**: Must be exactly `none`, `pick`, or `reject`
5. **Limit**: Maximum limit is 1000 results per request
6. **Field Boosting**: Title matches are weighted higher than other fields in search results
