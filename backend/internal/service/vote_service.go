package service

import (
	"fmt"

	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
)

type VoteService struct {
	Repo *data.Repository
}

func NewVoteService(repo *data.Repository) *VoteService {
	return &VoteService{
		Repo: repo,
	}
}

// VoteOnPost allows a user to vote on a post
func (voteService *VoteService) VoteOnPost(userID, postID, voteType int) error {
	// Validate voteType
	if voteType != 1 && voteType != -1 {
		return fmt.Errorf("invalid vote type: %d", voteType)
	}

	// Validate IDs
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	if postID <= 0 {
		return fmt.Errorf("invalid post ID: %d", postID)
	}

	// Check currennt vote status
	currentVote, err := voteService.Repo.GetUserVoteOnPost(userID, postID)
	if err != nil {
		return fmt.Errorf("failed to get current vote status: %w", err)
	}

	// If user clicks same vote, remove the vote
	if currentVote != nil && *currentVote == voteType {
		err = voteService.Repo.RemovePostVote(userID, postID)
		if err != nil {
			return fmt.Errorf("failed to remove vote: %w", err)
		}

		return nil
	}

	// Else add/update the vote
	err = voteService.Repo.VotePost(userID, postID, voteType)
	if err != nil {
		return fmt.Errorf("failed to cast vote: %w", err)
	}

	return nil
}

// RemoveVoteFromPost removes a user's vote from a post
func (voteService *VoteService) RemoveVoteFromPost(userID, postID int) error {
	// Validate IDs
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	if postID <= 0 {
		return fmt.Errorf("invalid post ID: %d", postID)
	}

	// Delegate call to repository layer
	err := voteService.Repo.RemovePostVote(userID, postID)
	if err != nil {
		return fmt.Errorf("failed to remove vote: %w", err)
	}

	return nil
}

// VoteOnComment allows a user to vote on a comment
func (voteService *VoteService) VoteOnComment(userID, commentID, voteType int) error {
	if voteType != 1 && voteType != -1 {
		return fmt.Errorf("invalid vote type: %d", voteType)
	}

	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	if commentID <= 0 {
		return fmt.Errorf("invalid comment ID: %d", commentID)
	}

	currentVote, err := voteService.Repo.GetUserVoteOnComment(userID, commentID)
	if err != nil {
		return fmt.Errorf("failed to get current vote status: %w", err)
	}

	if currentVote != nil && *currentVote == voteType {
		err = voteService.Repo.RemoveCommentVote(userID, commentID)
		if err != nil {
			return fmt.Errorf("failed to remove vote: %w", err)
		}

		return nil
	}

	err = voteService.Repo.VoteComment(userID, commentID, voteType)
	if err != nil {
		return fmt.Errorf("failed to cast vote: %w", err)
	}

	return nil
}

// RemoveVoteFromComment removes a user's vote from a comment
func (voteService *VoteService) RemoveVoteFromComment(userID, commentID int) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	if commentID <= 0 {
		return fmt.Errorf("invalid comment ID: %d", commentID)
	}

	err := voteService.Repo.RemoveCommentVote(userID, commentID)
	if err != nil {
		return fmt.Errorf("failed to remove vote: %w", err)
	}

	return nil
}
