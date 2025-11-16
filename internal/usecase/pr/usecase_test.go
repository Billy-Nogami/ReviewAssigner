package pr

import (
	"testing"
	"time"
	"ReviewAssigner/internal/domain/interfaces"
	"ReviewAssigner/internal/domain/schemas"
	pkgerrors "ReviewAssigner/internal/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock для UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(userID string) (*schemas.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.User), args.Error(1)
}

func (m *MockUserRepository) UpdateIsActive(userID string, isActive bool) (*schemas.User, error) {
	args := m.Called(userID, isActive)
	return args.Get(0).(*schemas.User), args.Error(1)
}

func (m *MockUserRepository) GetActiveByTeam(teamName string, excludeUserID string) ([]schemas.User, error) {
	args := m.Called(teamName, excludeUserID)
	return args.Get(0).([]schemas.User), args.Error(1)
}

// Mock для PullRequestRepository
type MockPullRequestRepository struct {
	mock.Mock
}

func (m *MockPullRequestRepository) Create(pr *schemas.PullRequest) error {
	args := m.Called(pr)
	return args.Error(0)
}

func (m *MockPullRequestRepository) GetByID(id string) (*schemas.PullRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.PullRequest), args.Error(1)
}

func (m *MockPullRequestRepository) UpdateStatus(id string, status string, mergedAt *time.Time) (*schemas.PullRequest, error) {
	args := m.Called(id, status, mergedAt)
	return args.Get(0).(*schemas.PullRequest), args.Error(1)
}

func (m *MockPullRequestRepository) UpdateReviewers(id string, reviewers []string) error {
	args := m.Called(id, reviewers)
	return args.Error(0)
}

func (m *MockPullRequestRepository) GetByReviewerID(userID string) ([]schemas.PullRequestShort, error) {
	args := m.Called(userID)
	return args.Get(0).([]schemas.PullRequestShort), args.Error(1)
}

func (m *MockPullRequestRepository) Exists(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func TestUsecase_CreatePR_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockPRRepo := &MockPullRequestRepository{}
	usecase := NewUsecase(mockUserRepo, mockPRRepo)

	author := &schemas.User{ID: "u1", TeamName: "backend"}
	candidates := []schemas.User{{ID: "u2"}}
	pr := &schemas.PullRequest{ID: "pr1", AuthorID: "u1", Status: "OPEN", AssignedReviewers: []string{"u2"}}

	mockPRRepo.On("Exists", "pr1").Return(false, nil)
	mockUserRepo.On("GetByID", "u1").Return(author, nil)
	mockUserRepo.On("GetActiveByTeam", "backend", "u1").Return(candidates, nil)
	mockPRRepo.On("Create", pr).Return(nil)

	result, err := usecase.CreatePR("pr1", "Test", "u1")
	assert.NoError(t, err)
	assert.Equal(t, pr.ID, result.ID)
}

func TestUsecase_MergePR_Idempotent(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockPRRepo := &MockPullRequestRepository{}
	usecase := NewUsecase(mockUserRepo, mockPRRepo)

	pr := &schemas.PullRequest{ID: "pr1", Status: "MERGED"}
	mockPRRepo.On("GetByID", "pr1").Return(pr, nil)

	result, err := usecase.MergePR("pr1")
	assert.NoError(t, err)
	assert.Equal(t, "MERGED", result.Status)
}

func TestUsecase_ReassignPR_NoCandidate(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockPRRepo := &MockPullRequestRepository{}
	usecase := NewUsecase(mockUserRepo, mockPRRepo)

	pr := &schemas.PullRequest{ID: "pr1", Status: "OPEN", AssignedReviewers: []string{"u2"}}
	oldUser := &schemas.User{ID: "u2", TeamName: "backend"}

	mockPRRepo.On("GetByID", "pr1").Return(pr, nil)
	mockUserRepo.On("GetByID", "u2").Return(oldUser, nil)
	mockUserRepo.On("GetActiveByTeam", "backend", "u1").Return([]schemas.User{}, nil)

	_, _, err := usecase.ReassignPR("pr1", "u2")
	assert.Equal(t, pkgerrors.ErrNoCandidate, err)
}
