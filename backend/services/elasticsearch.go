package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/suipic/backend/config"
	"github.com/suipic/backend/models"
)

const photosIndex = "photos"

type ElasticsearchService struct {
	client *elasticsearch.Client
}

type PhotoDocument struct {
	ID                int                    `json:"id"`
	AlbumID           int                    `json:"album_id"`
	Title             string                 `json:"title"`
	DateTime          *time.Time             `json:"date_time,omitempty"`
	ExifData          map[string]interface{} `json:"exif_data,omitempty"`
	AlbumTitle        string                 `json:"album_title"`
	AlbumLocation     string                 `json:"album_location"`
	AlbumCustomFields map[string]interface{} `json:"album_custom_fields,omitempty"`
	Comments          []string               `json:"comments"`
	PickRejectState   string                 `json:"pick_reject_state"`
	Stars             int                    `json:"stars"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

func NewElasticsearchService(cfg *config.ElasticsearchConfig) (*ElasticsearchService, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
	}

	if cfg.Username != "" {
		esCfg.Username = cfg.Username
		esCfg.Password = cfg.Password
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch returned error: %s", res.String())
	}

	service := &ElasticsearchService{
		client: client,
	}

	if err := service.createIndex(); err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	return service, nil
}

func (s *ElasticsearchService) createIndex() error {
	mapping := `{
		"mappings": {
			"properties": {
				"id": { "type": "integer" },
				"album_id": { "type": "integer" },
				"title": { "type": "text" },
				"date_time": { "type": "date" },
				"exif_data": { "type": "object", "enabled": true },
				"album_title": { "type": "text" },
				"album_location": { "type": "text" },
				"album_custom_fields": { "type": "object", "enabled": true },
				"comments": { "type": "text" },
				"pick_reject_state": { "type": "keyword" },
				"stars": { "type": "integer" },
				"created_at": { "type": "date" },
				"updated_at": { "type": "date" }
			}
		}
	}`

	res, err := s.client.Indices.Exists([]string{photosIndex})
	if err != nil {
		return fmt.Errorf("failed to check if index exists: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	res, err = s.client.Indices.Create(
		photosIndex,
		s.client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index: %s", res.String())
	}

	return nil
}

func (s *ElasticsearchService) IndexPhoto(ctx context.Context, photo *models.Photo, album *models.Album, comments []*models.Comment) error {
	doc := s.buildPhotoDocument(photo, album, comments)

	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      photosIndex,
		DocumentID: fmt.Sprintf("%d", photo.ID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to index document: %s", res.String())
	}

	return nil
}

func (s *ElasticsearchService) BulkIndexPhotos(ctx context.Context, photos []*models.Photo, albums map[int]*models.Album, commentsMap map[int][]*models.Comment) error {
	if len(photos) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, photo := range photos {
		album := albums[photo.AlbumID]
		comments := commentsMap[photo.ID]

		doc := s.buildPhotoDocument(photo, album, comments)

		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": photosIndex,
				"_id":    fmt.Sprintf("%d", photo.ID),
			},
		}

		metaData, err := json.Marshal(meta)
		if err != nil {
			return fmt.Errorf("failed to marshal meta: %w", err)
		}

		docData, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal document: %w", err)
		}

		buf.Write(metaData)
		buf.WriteByte('\n')
		buf.Write(docData)
		buf.WriteByte('\n')
	}

	res, err := s.client.Bulk(bytes.NewReader(buf.Bytes()), s.client.Bulk.WithContext(ctx), s.client.Bulk.WithRefresh("true"))
	if err != nil {
		return fmt.Errorf("failed to bulk index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to bulk index: %s", res.String())
	}

	return nil
}

func (s *ElasticsearchService) DeletePhoto(ctx context.Context, photoID int) error {
	req := esapi.DeleteRequest{
		Index:      photosIndex,
		DocumentID: fmt.Sprintf("%d", photoID),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("failed to delete document: %s", res.String())
	}

	return nil
}

func (s *ElasticsearchService) buildPhotoDocument(photo *models.Photo, album *models.Album, comments []*models.Comment) *PhotoDocument {
	doc := &PhotoDocument{
		ID:              photo.ID,
		AlbumID:         photo.AlbumID,
		Title:           "",
		DateTime:        photo.DateTime,
		ExifData:        photo.ExifData,
		PickRejectState: string(photo.PickRejectState),
		Stars:           photo.Stars,
		CreatedAt:       photo.CreatedAt,
		UpdatedAt:       photo.UpdatedAt,
	}

	if photo.Title != nil {
		doc.Title = *photo.Title
	}

	if album != nil {
		doc.AlbumTitle = album.Title
		if album.Location != nil {
			doc.AlbumLocation = *album.Location
		}
		doc.AlbumCustomFields = album.CustomFields
	}

	var commentTexts []string
	for _, comment := range comments {
		commentTexts = append(commentTexts, comment.Text)
	}
	doc.Comments = commentTexts

	return doc
}

type SearchFilter struct {
	Query    string
	AlbumID  *int
	DateFrom *time.Time
	DateTo   *time.Time
	MinStars *int
	MaxStars *int
	State    *string
	Limit    int
	Offset   int
}

type SearchResult struct {
	Total  int             `json:"total"`
	Photos []*models.Photo `json:"photos"`
}

func (s *ElasticsearchService) Search(ctx context.Context, filter *SearchFilter) (*SearchResult, error) {
	query := s.buildSearchQuery(filter)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(photosIndex),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithSize(filter.Limit),
		s.client.Search.WithFrom(filter.Offset),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	hits := result["hits"].(map[string]interface{})
	total := int(hits["total"].(map[string]interface{})["value"].(float64))

	var photos []*models.Photo
	for _, hit := range hits["hits"].([]interface{}) {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})

		photo := &models.Photo{
			ID:              int(source["id"].(float64)),
			AlbumID:         int(source["album_id"].(float64)),
			PickRejectState: models.PickRejectState(source["pick_reject_state"].(string)),
			Stars:           int(source["stars"].(float64)),
		}

		if title, ok := source["title"].(string); ok && title != "" {
			photo.Title = &title
		}

		if dateTime, ok := source["date_time"].(string); ok && dateTime != "" {
			t, _ := time.Parse(time.RFC3339, dateTime)
			photo.DateTime = &t
		}

		if exifData, ok := source["exif_data"].(map[string]interface{}); ok {
			photo.ExifData = exifData
		}

		if createdAt, ok := source["created_at"].(string); ok {
			t, _ := time.Parse(time.RFC3339, createdAt)
			photo.CreatedAt = t
		}

		if updatedAt, ok := source["updated_at"].(string); ok {
			t, _ := time.Parse(time.RFC3339, updatedAt)
			photo.UpdatedAt = t
		}

		photos = append(photos, photo)
	}

	return &SearchResult{
		Total:  total,
		Photos: photos,
	}, nil
}

func (s *ElasticsearchService) buildSearchQuery(filter *SearchFilter) map[string]interface{} {
	must := []interface{}{}

	if filter.Query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": filter.Query,
				"fields": []string{
					"title^3",
					"album_title^2",
					"album_location",
					"comments",
					"exif_data.*",
				},
			},
		})
	}

	if filter.AlbumID != nil {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"album_id": *filter.AlbumID,
			},
		})
	}

	if filter.DateFrom != nil || filter.DateTo != nil {
		rangeQuery := map[string]interface{}{}
		if filter.DateFrom != nil {
			rangeQuery["gte"] = filter.DateFrom.Format(time.RFC3339)
		}
		if filter.DateTo != nil {
			rangeQuery["lte"] = filter.DateTo.Format(time.RFC3339)
		}
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				"date_time": rangeQuery,
			},
		})
	}

	if filter.MinStars != nil || filter.MaxStars != nil {
		rangeQuery := map[string]interface{}{}
		if filter.MinStars != nil {
			rangeQuery["gte"] = *filter.MinStars
		}
		if filter.MaxStars != nil {
			rangeQuery["lte"] = *filter.MaxStars
		}
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				"stars": rangeQuery,
			},
		})
	}

	if filter.State != nil {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"pick_reject_state": *filter.State,
			},
		})
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"sort": []interface{}{
			map[string]interface{}{
				"date_time": map[string]interface{}{
					"order":   "desc",
					"missing": "_last",
				},
			},
			map[string]interface{}{
				"created_at": "desc",
			},
		},
	}

	if len(must) == 0 {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	return query
}
