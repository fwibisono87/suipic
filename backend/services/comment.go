package services

import (
	"context"
	"fmt"

	"github.com/suipic/backend/models"
	"github.com/suipic/backend/repository"
)

type CommentService struct {
	commentRepo repository.CommentRepository
}

func NewCommentService(commentRepo repository.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
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

func (s *CommentService) GetCommentsByPhoto(ctx context.Context, photoID int) ([]*models.Comment, error) {
	return s.commentRepo.GetByPhoto(ctx, photoID)
}

type ThreadedComment struct {
	*models.Comment
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
		threadedComment := &ThreadedComment{
			Comment: comment,
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
