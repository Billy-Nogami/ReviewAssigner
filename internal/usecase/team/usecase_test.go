package team

import (
	"testing"
	"ReviewAssigner/internal/domain/schemas"
	pkgerrors "ReviewAssigner/internal/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock для TeamRepository
type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) Create(team *schemas.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

func (m *MockTeamRepository) GetByName(name string) (*schemas.Team, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.Team), args.Error(1)
}

func (m *MockTeamRepository) Exists(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}

func TestUsecase_CreateTeam_Success(t *testing.T) {
	mockRepo := &MockTeamRepository{}
	usecase := NewUsecase(mockRepo)

	team := &schemas.Team{Name: "test", Members: []schemas.User{{ID: "u1"}}}
	mockRepo.On("Exists", "test").Return(false, nil)
	mockRepo.On("Create", team).Return(nil)
	mockRepo.On("GetByName", "test").Return(team, nil)

	result, err := usecase.CreateTeam(team)
	assert.NoError(t, err)
	assert.Equal(t, team, result)
	mockRepo.AssertExpectations(t)
}

func TestUsecase_CreateTeam_Exists(t *testing.T) {
	mockRepo := &MockTeamRepository{}
	usecase := NewUsecase(mockRepo)

	mockRepo.On("Exists", "test").Return(true, nil)

	_, err := usecase.CreateTeam(&schemas.Team{Name: "test"})
	assert.Equal(t, pkgerrors.ErrTeamExists, err)
}

func TestUsecase_GetTeam_Success(t *testing.T) {
	mockRepo := &MockTeamRepository{}
	usecase := NewUsecase(mockRepo)

	team := &schemas.Team{Name: "test"}
	mockRepo.On("GetByName", "test").Return(team, nil)

	result, err := usecase.GetTeam("test")
	assert.NoError(t, err)
	assert.Equal(t, team, result)
}

func TestUsecase_GetTeam_NotFound(t *testing.T) {
	mockRepo := &MockTeamRepository{}
	usecase := NewUsecase(mockRepo)

	mockRepo.On("GetByName", "test").Return(nil, nil)

	_, err := usecase.GetTeam("test")
	assert.Equal(t, pkgerrors.ErrNotFound, err)
}
