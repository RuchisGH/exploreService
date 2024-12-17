package grpc

import (
	"context"
	"testing"
	"time"

	genproto "example.com/mod/internal/genproto"
	"example.com/mod/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock db
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetReceivedLikes(userID string) ([]string, error) {
	args := m.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockDatabase) GetGivenDecisions(userID string) ([]models.Decision, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Decision), args.Error(1)
}

func (m *MockDatabase) UpsertDecision(userID, targetID string, decision bool) error {
	args := m.Called(userID, targetID, decision)
	return args.Error(0)
}

func TestListLikedYou(t *testing.T) {
	mockDB := new(MockDatabase)
	server := ExploreServer{DB: mockDB} // Inject the mock database

	mockDB.On("GetReceivedLikes", "user2").Return([]string{"user1", "user3"}, nil)

	req := &genproto.ListLikedYouRequest{
		RecipientUserId: "user2",
	}
	resp, err := server.ListLikedYou(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Likers, 2)
	assert.Equal(t, "user1", resp.Likers[0].ActorId)
	assert.Equal(t, "user3", resp.Likers[1].ActorId)
}

func TestListNewLikedYou(t *testing.T) {
	mockDB := new(MockDatabase)
	server := ExploreServer{DB: mockDB} // Inject the mock database

	// GetReceivedLikes: "user2" has been liked by "user1", "user3"
	mockDB.On("GetReceivedLikes", "user2").Return([]string{"user1", "user3"}, nil)

	// GetGivenDecisions: "user2" has liked back "user3"
	mockDB.On("GetGivenDecisions", "user2").Return([]models.Decision{
		{TargetID: "user3", Decision: "LIKE", Timestamp: time.Now()},
	}, nil)

	req := &genproto.ListLikedYouRequest{
		RecipientUserId: "user2",
	}

	// Call the ListNewLikedYou method
	resp, err := server.ListNewLikedYou(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Likers, 1)
	assert.Equal(t, "user1", resp.Likers[0].ActorId)
	assert.Equal(t, uint64(1617180000), resp.Likers[0].UnixTimestamp)

	mockDB.AssertExpectations(t)
}

func TestCountLikedYou(t *testing.T) {
	mockDB := new(MockDatabase)
	server := ExploreServer{DB: mockDB} // Inject the mock database

	// GetReceivedLikes: "user2" has been liked by "user1", "user3"
	mockDB.On("GetReceivedLikes", "user2").Return([]string{"user1", "user3"}, nil)

	// GetGivenDecisions: "user2" has liked back "user3"
	mockDB.On("GetGivenDecisions", "user2").Return([]models.Decision{
		{TargetID: "user3", Decision: "LIKE", Timestamp: time.Now()},
	}, nil)

	req := &genproto.CountLikedYouRequest{
		RecipientUserId: "user2",
	}

	resp, err := server.CountLikedYou(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// Assert that the count of new likers is 1 (since "user3" liked back)
	assert.Equal(t, uint64(1), resp.Count)

	mockDB.AssertExpectations(t)
}

func TestPutDecision(t *testing.T) {
	mockDB := new(MockDatabase)
	server := ExploreServer{DB: mockDB} // Inject the mock database

	mockDB.On("UpsertDecision", "user1", "user2", true).Return(nil)

	mockDB.On("GetGivenDecisions", "user2").Return([]models.Decision{
		{TargetID: "user1", Decision: "LIKE", Timestamp: time.Now()},
	}, nil)

	req := &genproto.PutDecisionRequest{
		ActorUserId:     "user1",
		RecipientUserId: "user2",
		LikedRecipient:  true, // "user1" liked "user2"
	}

	// Call the PutDecision method
	resp, err := server.PutDecision(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Simulate the case where there's no mutual like
	mockDB.On("UpsertDecision", "user1", "user2", true).Return(nil)
	mockDB.On("GetGivenDecisions", "user2").Return([]models.Decision{}, nil)

	// Call the PutDecision method again for no mutual like case
	resp, err = server.PutDecision(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// Assert that mutual like is false because "user2" hasn't liked "user1" back
	assert.False(t, resp.MutualLikes)

	mockDB.AssertExpectations(t)
}
