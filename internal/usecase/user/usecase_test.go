package user

import (
	"testing"
	"ReviewAssigner/internal/domain/schema"
	pkgerrors "ReviewAssigner/internal/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock для UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(userID string) (*schema.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.User), args.Error(1)
}

func (m *MockUserRepository) UpdateIsActive(userID string, isActive bool) (*schema.User, error) {
	args := m.Called(userID, isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.User), args.Error(1)
}

func (m *MockUserRepository) GetActiveByTeam(teamName string, excludeUserID string) ([]schema.User, error) {
	args := m.Called(teamName, excludeUserID)
	return args.Get(0).([]schema.User), args.Error(1)
}

// Mock для PullRequestRepository (добавлены все методы интерфейса)
type MockPullRequestRepository struct {
	mock.Mock
}

func (m *MockPullRequestRepository) Create(pr *schema.PullRequest) error {
	args := m.Called(pr)
	return args.Error(0)
}

func (m *MockPullRequestRepository) GetByID(id string) (*schema.PullRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.PullRequest), args.Error(1)
}

func (m *MockPullRequestRepository) UpdateStatus(id string, status string, mergedAt *time.Time) (*schema.PullRequest, error) {
	args := m.Called(id, status, mergedAt)
	return args.Get(0).(*schema.PullRequest), args.Error(1)
}

func (m *MockPullRequestRepository) UpdateReviewers(id string, reviewers []string) error {
	args := m.Called(id, reviewers)
	return args.Error(0)
}

func (m *MockPullRequestRepository) GetByReviewerID(userID string) ([]schema.PullRequestShort, error) {
	args := m.Called(userID)
	return args.Get(0).([]schema.PullRequestShort), args.Error(1)
}

func (m *MockPullRequestRepository) Exists(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func TestUsecase_SetIsActive_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockPRRepo := &MockPullRequestRepository{}
	usecase := NewUsecase(mockUserRepo, mockPRRepo)

	user := &schema.User{ID: "u1", IsActive: true}
	mockUserRepo.On("UpdateIsActive", "u1", false).Return(user, nil)

	result, err := usecase.SetIsActive("u1", false)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestUsecase_SetIsActive_NotFound(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockPRRepo := &MockPullRequestRepository{}
	usecase := NewUsecase(mockUserRepo, mockPRRepo)

	mockUserRepo.On("UpdateIsActive", "u1", false).Return(nil, nil)

	_, err := usecase.SetIsActive("u1", false)
	assert.Equal(t, pkgerrors.ErrNotFound, err)
}

func TestUsecase_GetUserReviews_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockPRRepo := &MockPullRequestRepository{}
	usecase := NewUsecase(mockUserRepo, mockPRRepo)

	user := &schema.User{ID: "u1"}
	prs := []schema.PullRequestShort{{ID: "pr1"}}
	mockUserRepo.On("GetByID", "u1").Return(user, nil)
	mockPRRepo.On("GetByReviewerID", "u1").Return(prs, nil)

	resultUser, resultPRs, err := usecase.GetUserReviews("u1")
	assert.NoError(t, err)
	assert.Equal(t, user, resultUser)
	assert.Equal(t, prs, resultPRs)
}
