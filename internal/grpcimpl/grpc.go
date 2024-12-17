package grpc

import (
	"context"
	"fmt"

	db "example.com/mod/internal/db"
	genproto "example.com/mod/internal/genproto"
	"google.golang.org/grpc"
)

type ExploreServer struct {
	genproto.UnimplementedExploreServiceServer
	DB db.DatabaseInterface
}

// ListLikedYou lists all users who liked the recipient
func (s *ExploreServer) ListLikedYou(ctx context.Context, req *genproto.ListLikedYouRequest) (*genproto.ListLikedYouResponse, error) {
	likers, err := s.DB.GetReceivedLikes(req.RecipientUserId)
	if err != nil {
		return nil, err
	}

	var likerProtos []*genproto.ListLikedYouResponse_Liker
	for _, liker := range likers {
		likerProtos = append(likerProtos, &genproto.ListLikedYouResponse_Liker{
			ActorId:       liker,
			UnixTimestamp: 1617180000,
		})
	}

	return &genproto.ListLikedYouResponse{Likers: likerProtos}, nil
}

// ListNewLikedYou lists users who liked the recipient excluding those who have been liked in return
func (s *ExploreServer) ListNewLikedYou(ctx context.Context, req *genproto.ListLikedYouRequest) (*genproto.ListLikedYouResponse, error) {
	// Get all users who liked the recipient (target_id = recipient_user_id)
	likers, err := s.DB.GetReceivedLikes(req.RecipientUserId)
	if err != nil {
		return nil, fmt.Errorf("error fetching received likes: %v", err)
	}

	// Get all users who the recipient has liked (user_id = recipient_user_id and decision = 'LIKE')
	likedBackUsers, err := s.DB.GetGivenDecisions(req.RecipientUserId)
	if err != nil {
		return nil, fmt.Errorf("error fetching given decisions: %v", err)
	}

	// Create a set of users who the recipient has liked back
	likedBackSet := make(map[string]struct{})
	for _, decision := range likedBackUsers {
		likedBackSet[decision.TargetID] = struct{}{}
	}

	// Filter out the mutual likers from the "likers" list
	var newLikers []*genproto.ListLikedYouResponse_Liker
	for _, liker := range likers {
		if _, found := likedBackSet[liker]; !found {
			newLikers = append(newLikers, &genproto.ListLikedYouResponse_Liker{
				ActorId:       liker,
				UnixTimestamp: 1617180000, // Should use actual timestamp
			})
		}
	}

	//Return the response with the new likers
	return &genproto.ListLikedYouResponse{Likers: newLikers}, nil
}

// CountLikedYou returns the count of users who liked the recipient
func (s *ExploreServer) CountLikedYou(ctx context.Context, req *genproto.CountLikedYouRequest) (*genproto.CountLikedYouResponse, error) {
	// Get all users who liked the recipient (target_id = recipient_user_id)
	likers, err := s.DB.GetReceivedLikes(req.RecipientUserId)
	if err != nil {
		return nil, fmt.Errorf("error fetching received likes: %v", err)
	}

	// Get all users who the recipient has liked (user_id = recipient_user_id and decision = 'LIKE')
	likedBackUsers, err := s.DB.GetGivenDecisions(req.RecipientUserId)
	if err != nil {
		return nil, fmt.Errorf("error fetching given decisions: %v", err)
	}

	// Create a set of users who the recipient has liked back
	likedBackSet := make(map[string]struct{})
	for _, decision := range likedBackUsers {
		likedBackSet[decision.TargetID] = struct{}{}
	}

	// Filter out the mutual likers from the "likers" list
	var newLikers []string
	for _, liker := range likers {
		if _, found := likedBackSet[liker]; !found {
			newLikers = append(newLikers, liker)
		}
	}

	// Return the count of new likers
	return &genproto.CountLikedYouResponse{
		Count: uint64(len(newLikers)),
	}, nil
}

// PutDecision records the decision of the actor to like or pass the recipient
func (s *ExploreServer) PutDecision(ctx context.Context, req *genproto.PutDecisionRequest) (*genproto.PutDecisionResponse, error) {
	// Upsert the user's decision (INSERT or UPDATE)
	err := s.DB.UpsertDecision(req.ActorUserId, req.RecipientUserId, req.LikedRecipient)
	if err != nil {
		return nil, fmt.Errorf("error inserting/updating decision: %v", err)
	}

	// Check if the target has already liked back the user (mutual like)
	// Get the decision made by the target towards the user
	targetDecision, err := s.DB.GetGivenDecisions(req.RecipientUserId)
	if err != nil {
		return nil, fmt.Errorf("error fetching target's decisions: %v", err)
	}

	// Determine if mutual like exists
	mutualLike := false
	for _, decision := range targetDecision {
		// If the target has liked the user back, it's a mutual like
		if decision.TargetID == req.RecipientUserId && decision.Decision == "LIKE" {
			mutualLike = true
			break
		}
	}

	// Return the response with the mutual like status
	return &genproto.PutDecisionResponse{
		MutualLikes: mutualLike, // Indicating if mutual like exists
	}, nil
}

// Register the server
func RegisterServer(s *grpc.Server) {
	server := &ExploreServer{}
	genproto.RegisterExploreServiceServer(s, server)
}
