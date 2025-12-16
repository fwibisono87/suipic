package services

import (
	"context"
	"fmt"

	"github.com/suipic/backend/models"
	"github.com/suipic/backend/repository"
)

type CommentService struct {
	commentRepo repository.CommentRepository
	userRepo    repository.UserRepository
}

func NewCommentService(commentRepo repository.CommentRepository, userRepo repository.UserRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
	}
}

func (s *CommentService) CreateComment(ctx context.Context, comment *models.Comment) error {
	if comment.ParentCommentID != nil {
		parent, err := s.commentRepo.GetByID(ctx, *comment.ParentCommentID)
		if err != nil {
			return fmt.Errorf("failed to get parent comment: %w", err)
		}
		if parent == nil {
			return fmt.Errorf("parent comment not found")
		}
		if parent.PhotoID != comment.PhotoID {
			return fmt.Errorf("parent comment does not belong to the same photo")
		}
	}

	return s.commentRepo.Create(ctx, comment)
}

func (s *CommentService) GetCommentWithUser(ctx context.Context, commentID int) (*CommentWithUser, error) {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, fmt.Errorf("comment not found")
	}

	user, err := s.userRepo.GetByID(ctx, comment.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user for comment: %w", err)
	}

	return &CommentWithUser{
		Comment: comment,
		User:    user,
	}, nil
}

func (s *CommentService) GetCommentsByPhoto(ctx context.Context, photoID int) ([]*models.Comment, error) {
	return s.commentRepo.GetByPhoto(ctx, photoID)
}

type CommentWithUser struct {
	*models.Comment
	User *models.User `json:"user"`
}

type ThreadedComment struct {
	*models.Comment
	User    *models.User       `json:"user"`
	Replies []*ThreadedComment `json:"replies,omitempty"`
}

func (s *CommentService) GetThreadedComments(ctx context.Context, photoID int) ([]*ThreadedComment, error) {
	allComments, err := s.commentRepo.GetByPhoto(ctx, photoID)
	if err != nil {
		return nil, err
	}

	commentMap := make(map[int]*ThreadedComment)
	var threads []*ThreadedComment

	for _, comment := range allComments {
		user, err := s.userRepo.GetByID(ctx, comment.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user for comment: %w", err)
		}

		threadedComment := &ThreadedComment{
			Comment: comment,
			User:    user,
			Replies: []*ThreadedComment{},
		}
		commentMap[comment.ID] = threadedComment
	}

	for _, comment := range allComments {
		threadedComment := commentMap[comment.ID]
		if comment.ParentCommentID == nil {
			threads = append(threads, threadedComment)
		} else {
			if parent, ok := commentMap[*comment.ParentCommentID]; ok {
				parent.Replies = append(parent.Replies, threadedComment)
			}
		}
	}

	return threads, nil
}
