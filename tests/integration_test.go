package tests

import (
	"testing"
	"ReviewAssigner/internal/repository/inmemory"
	"ReviewAssigner/internal/usecase/pr"
	"ReviewAssigner/internal/usecase/team"
	"ReviewAssigner/internal/usecase/user"
	"ReviewAssigner/internal/domain/schemas"

	"github.com/stretchr/testify/assert"
)

func TestIntegration_FullFlow(t *testing.T) {
	userRepo := inmemory.NewUserRepository()
	teamRepo := inmemory.NewTeamRepository()
	prRepo := inmemory.NewPullRequestRepository()

	// Инициализация данных
	userRepo.(*inmemory.UserRepository).AddUser(&schemas.User{ID: "u1", Username: "Alice", TeamName: "backend", IsActive: true})
	userRepo.(*inmemory.UserRepository).AddUser(&schemas.User{ID: "u2", Username: "Bob", TeamName: "backend", IsActive: true})
	userRepo.(*inmemory.UserRepository).AddUser(&schemas.User{ID: "u3", Username: "Charlie", TeamName: "backend", IsActive: true})
	teamRepo.(*inmemory.TeamRepository).AddTeam(&schemas.Team{Name: "backend", Members: []schemas.User{
		{ID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
		{ID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
		{ID: "u3", Username: "Charlie", TeamName: "backend", IsActive: true},
	}})

	teamUsecase := team.NewUsecase(teamRepo)
	userUsecase := user.NewUsecase(userRepo, prRepo)
	prUsecase := pr.NewUsecase(userRepo, prRepo)

	// Создать команду
	team, err := teamUsecase.CreateTeam(&schemas.Team{Name: "backend", Members: []schemas.User{
		{ID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
	}})
	assert.NoError(t, err)

	// Создать PR
	pr, err := prUsecase.CreatePR("pr1", "Test PR", "u1")
	assert.NoError(t, err)
	assert.Len(t, pr.AssignedReviewers, 2) // u2, u3

	// Merge
	pr, err = prUsecase.MergePR("pr1")
	assert.NoError(t, err)
	assert.Equal(t, "MERGED", pr.Status)

	// Попытка reassign после merge
	_, _, err = prUsecase.ReassignPR("pr1", "u2")
	assert.Equal(t, pkgerrors.ErrPRMerged, err)
}
