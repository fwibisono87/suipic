package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/suipic/backend/models"
	"github.com/suipic/backend/repository"
)

type AlbumService struct {
	albumRepo     repository.AlbumRepository
	albumUserRepo repository.AlbumUserRepository
}

func NewAlbumService(db *sql.DB) *AlbumService {
	return &AlbumService{
		albumRepo:     repository.NewPostgresAlbumRepository(db),
		albumUserRepo: repository.NewPostgresAlbumUserRepository(db),
	}
}

func (s *AlbumService) CreateAlbum(ctx context.Context, album *models.Album) error {
	return s.albumRepo.Create(ctx, album)
}

func (s *AlbumService) GetAlbumByID(ctx context.Context, id int) (*models.Album, error) {
	return s.albumRepo.GetByID(ctx, id)
}

func (s *AlbumService) UpdateAlbum(ctx context.Context, album *models.Album) error {
	return s.albumRepo.Update(ctx, album)
}

func (s *AlbumService) DeleteAlbum(ctx context.Context, id int) error {
	return s.albumRepo.Delete(ctx, id)
}

func (s *AlbumService) ListAlbums(ctx context.Context, photographerID *int, userID *int) ([]*models.Album, error) {
	if photographerID != nil {
		return s.albumRepo.GetByPhotographer(ctx, *photographerID)
	}

	if userID != nil {
		albumUsers, err := s.albumUserRepo.GetByUser(ctx, *userID)
		if err != nil {
			return nil, err
		}

		var albums []*models.Album
		for _, au := range albumUsers {
			album, err := s.albumRepo.GetByID(ctx, au.AlbumID)
			if err != nil {
				return nil, err
			}
			if album != nil {
				albums = append(albums, album)
			}
		}
		return albums, nil
	}

	return s.albumRepo.List(ctx, 1000, 0)
}

func (s *AlbumService) AssignUsersToAlbum(ctx context.Context, albumID int, userIDs []int) error {
	if err := s.albumUserRepo.DeleteByAlbum(ctx, albumID); err != nil {
		return fmt.Errorf("failed to clear existing users: %w", err)
	}

	for _, userID := range userIDs {
		albumUser := &models.AlbumUser{
			AlbumID: albumID,
			UserID:  userID,
		}
		if err := s.albumUserRepo.Create(ctx, albumUser); err != nil {
			return fmt.Errorf("failed to assign user %d to album: %w", userID, err)
		}
	}
	return nil
}

func (s *AlbumService) GetAlbumUsers(ctx context.Context, albumID int) ([]*models.AlbumUser, error) {
	return s.albumUserRepo.GetByAlbum(ctx, albumID)
}

func (s *AlbumService) CanUserAccessAlbum(ctx context.Context, userID int, albumID int) (bool, error) {
	album, err := s.albumRepo.GetByID(ctx, albumID)
	if err != nil {
		return false, err
	}
	if album == nil {
		return false, fmt.Errorf("album not found")
	}

	if album.PhotographerID == userID {
		return true, nil
	}

	albumUsers, err := s.albumUserRepo.GetByAlbum(ctx, albumID)
	if err != nil {
		return false, err
	}

	for _, au := range albumUsers {
		if au.UserID == userID {
			return true, nil
		}
	}

	return false, nil
}
