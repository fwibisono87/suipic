package repository

import (
	"context"

	"github.com/suipic/backend/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
	FindClientsByUsername(ctx context.Context, username string) ([]*models.User, error)
}

type AlbumRepository interface {
	Create(ctx context.Context, album *models.Album) error
	GetByID(ctx context.Context, id int) (*models.Album, error)
	Update(ctx context.Context, album *models.Album) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*models.Album, error)
	GetByPhotographer(ctx context.Context, photographerID int) ([]*models.Album, error)
}

type PhotoRepository interface {
	Create(ctx context.Context, photo *models.Photo) error
	GetByID(ctx context.Context, id int) (*models.Photo, error)
	Update(ctx context.Context, photo *models.Photo) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*models.Photo, error)
	GetByAlbum(ctx context.Context, albumID int) ([]*models.Photo, error)
}

type AlbumUserRepository interface {
	Create(ctx context.Context, albumUser *models.AlbumUser) error
	GetByID(ctx context.Context, id int) (*models.AlbumUser, error)
	Delete(ctx context.Context, id int) error
	DeleteByAlbumAndUser(ctx context.Context, albumID, userID int) error
	List(ctx context.Context, limit, offset int) ([]*models.AlbumUser, error)
	GetByAlbum(ctx context.Context, albumID int) ([]*models.AlbumUser, error)
	GetByUser(ctx context.Context, userID int) ([]*models.AlbumUser, error)
}

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id int) (*models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*models.Comment, error)
	GetByPhoto(ctx context.Context, photoID int) ([]*models.Comment, error)
	GetThreads(ctx context.Context, photoID int) ([]*models.Comment, error)
	GetReplies(ctx context.Context, parentCommentID int) ([]*models.Comment, error)
}
