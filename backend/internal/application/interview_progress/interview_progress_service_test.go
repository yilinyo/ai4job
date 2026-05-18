package interview_progress_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yilin/ai-for-backend/internal/application/activitylog"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	appsvc "github.com/yilin/ai-for-backend/internal/application/interview_progress"
	domainactivity "github.com/yilin/ai-for-backend/internal/domain/activity_log"
	domainprogress "github.com/yilin/ai-for-backend/internal/domain/interview_progress"
	domainapp "github.com/yilin/ai-for-backend/internal/domain/resume_application"
	domainerrors "github.com/yilin/ai-for-backend/pkg/errors"
)

type spyActivityLogger struct {
	Entries []activitylog.Entry
	Err     error
}

func (l *spyActivityLogger) Log(ctx context.Context, entry activitylog.Entry) error {
	l.Entries = append(l.Entries, entry)
	return l.Err
}

type mockInterviewProgressRepo struct {
	mock.Mock
}

func (m *mockInterviewProgressRepo) Create(ctx context.Context, p *domainprogress.InterviewProgress) error {
	return m.Called(ctx, p).Error(0)
}

func (m *mockInterviewProgressRepo) FindByIDAndUserID(ctx context.Context, id, userID string) (*domainprogress.InterviewProgress, error) {
	args := m.Called(ctx, id, userID)
	progress, _ := args.Get(0).(*domainprogress.InterviewProgress)
	return progress, args.Error(1)
}

func (m *mockInterviewProgressRepo) Update(ctx context.Context, p *domainprogress.InterviewProgress) error {
	return m.Called(ctx, p).Error(0)
}

func (m *mockInterviewProgressRepo) Delete(ctx context.Context, id, userID string) error {
	return m.Called(ctx, id, userID).Error(0)
}

func (m *mockInterviewProgressRepo) ListByApplicationAndUser(ctx context.Context, applicationID, userID string) ([]domainprogress.InterviewProgress, error) {
	args := m.Called(ctx, applicationID, userID)
	progresses, _ := args.Get(0).([]domainprogress.InterviewProgress)
	return progresses, args.Error(1)
}

func (m *mockInterviewProgressRepo) DeleteByApplicationID(ctx context.Context, applicationID, userID string) error {
	return m.Called(ctx, applicationID, userID).Error(0)
}

type mockApplicationService struct {
	mock.Mock
}

func (m *mockApplicationService) Get(ctx context.Context, userID, id string) (*domainapp.ResumeApplication, error) {
	args := m.Called(ctx, userID, id)
	application, _ := args.Get(0).(*domainapp.ResumeApplication)
	return application, args.Error(1)
}

func (m *mockApplicationService) UpdateCurrentStatus(ctx context.Context, userID, id string, status domainapp.Status) error {
	return m.Called(ctx, userID, id, status).Error(0)
}

func (m *mockApplicationService) SyncCurrentStatus(ctx context.Context, userID, id string, status domainapp.Status) error {
	return m.Called(ctx, userID, id, status).Error(0)
}

func TestInterviewProgressService_CreateSyncsFailedResultToRejectedAndLogsOnlyProgressActivity(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	logger := &spyActivityLogger{}
	svc := appsvc.NewInterviewProgressService(repo, apps).WithActivityLogger(logger)
	occurredAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	req := dto.CreateInterviewProgressRequest{
		ApplicationID:     "app-1",
		Stage:             domainprogress.StageFirstInterview,
		OccurredAt:        dto.FlexibleTime{Time: occurredAt},
		Result:            domainprogress.ResultFailed,
		InterviewerOrTeam: "Engineering",
		Note:              "did not pass",
	}

	apps.On("Get", mock.Anything, "user-1", "app-1").Return(&domainapp.ResumeApplication{ID: "app-1", UserID: "user-1"}, nil).Once()
	repo.On("Create", mock.Anything, mock.MatchedBy(func(p *domainprogress.InterviewProgress) bool {
		return p.ID != "" && p.UserID == "user-1" && p.ApplicationID == req.ApplicationID && p.Stage == req.Stage && p.OccurredAt.Equal(occurredAt) && p.Result == req.Result && p.InterviewerOrTeam == req.InterviewerOrTeam && p.Note == req.Note
	})).Return(nil).Once()
	repo.On("ListByApplicationAndUser", mock.Anything, "app-1", "user-1").Return([]domainprogress.InterviewProgress{{
		ID:            "progress-1",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         domainprogress.StageFirstInterview,
		OccurredAt:    occurredAt,
		Result:        domainprogress.ResultFailed,
	}}, nil).Once()
	apps.On("SyncCurrentStatus", mock.Anything, "user-1", "app-1", domainapp.StatusRejected).Return(nil).Once()

	got, err := svc.Create(context.Background(), "user-1", req)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "user-1", got.UserID)
	assert.Equal(t, req.ApplicationID, got.ApplicationID)
	require.Len(t, logger.Entries, 1)
	assert.Equal(t, domainactivity.ActionCreated, logger.Entries[0].Action)
	assert.Equal(t, domainactivity.TargetTypeInterview, logger.Entries[0].TargetType)
	repo.AssertExpectations(t)
	apps.AssertExpectations(t)
}

func TestInterviewProgressService_CreateSyncsLatestNonFailedStageToMatchingApplicationStatus(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	svc := appsvc.NewInterviewProgressService(repo, apps)
	older := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)
	latest := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	req := dto.CreateInterviewProgressRequest{
		ApplicationID: "app-1",
		Stage:         domainprogress.StageOffer,
		OccurredAt:    dto.FlexibleTime{Time: latest},
		Result:        domainprogress.ResultPassed,
	}

	apps.On("Get", mock.Anything, "user-1", "app-1").Return(&domainapp.ResumeApplication{ID: "app-1", UserID: "user-1"}, nil).Once()
	repo.On("Create", mock.Anything, mock.AnythingOfType("*interview_progress.InterviewProgress")).Return(nil).Once()
	repo.On("ListByApplicationAndUser", mock.Anything, "app-1", "user-1").Return([]domainprogress.InterviewProgress{
		{ID: "old", UserID: "user-1", ApplicationID: "app-1", Stage: domainprogress.StageViewed, OccurredAt: older, Result: domainprogress.ResultPassed},
		{ID: "new", UserID: "user-1", ApplicationID: "app-1", Stage: domainprogress.StageOffer, OccurredAt: latest, Result: domainprogress.ResultPending},
	}, nil).Once()
	apps.On("SyncCurrentStatus", mock.Anything, "user-1", "app-1", domainapp.StatusOffer).Return(nil).Once()

	got, err := svc.Create(context.Background(), "user-1", req)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "user-1", got.UserID)
	assert.Equal(t, req.ApplicationID, got.ApplicationID)
	repo.AssertExpectations(t)
	apps.AssertExpectations(t)
}

func TestInterviewProgressService_CreateIgnoresActivityLoggingErrorAfterCreateAndSync(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	logger := &spyActivityLogger{Err: errors.New("activity log failed")}
	svc := appsvc.NewInterviewProgressService(repo, apps).WithActivityLogger(logger)
	occurredAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	req := dto.CreateInterviewProgressRequest{
		ApplicationID: "app-1",
		Stage:         domainprogress.StageOffer,
		OccurredAt:    dto.FlexibleTime{Time: occurredAt},
		Result:        domainprogress.ResultPassed,
	}

	apps.On("Get", mock.Anything, "user-1", "app-1").Return(&domainapp.ResumeApplication{ID: "app-1", UserID: "user-1"}, nil).Once()
	repo.On("Create", mock.Anything, mock.AnythingOfType("*interview_progress.InterviewProgress")).Return(nil).Once()
	repo.On("ListByApplicationAndUser", mock.Anything, "app-1", "user-1").Return([]domainprogress.InterviewProgress{{
		ID:            "progress-1",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         domainprogress.StageOffer,
		OccurredAt:    occurredAt,
		Result:        domainprogress.ResultPassed,
	}}, nil).Once()
	apps.On("SyncCurrentStatus", mock.Anything, "user-1", "app-1", domainapp.StatusOffer).Return(nil).Once()

	got, err := svc.Create(context.Background(), "user-1", req)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Len(t, logger.Entries, 1)
	repo.AssertExpectations(t)
	apps.AssertExpectations(t)
}

func TestInterviewProgressService_CreateRejectsUnownedOrMissingApplicationBeforeCreate(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	svc := appsvc.NewInterviewProgressService(repo, apps)
	apps.On("Get", mock.Anything, "user-1", "app-1").Return(nil, domainerrors.ErrNotFound).Once()

	got, err := svc.Create(context.Background(), "user-1", dto.CreateInterviewProgressRequest{
		ApplicationID: "app-1",
		Stage:         domainprogress.StageFirstInterview,
		OccurredAt:    dto.FlexibleTime{Time: time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)},
		Result:        domainprogress.ResultPending,
	})

	assert.Nil(t, got)
	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	apps.AssertExpectations(t)
}

func TestInterviewProgressService_ListRejectsEmptyApplicationID(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	svc := appsvc.NewInterviewProgressService(repo, apps)

	progresses, err := svc.List(context.Background(), "user-1", "")

	assert.Nil(t, progresses)
	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
	repo.AssertNotCalled(t, "ListByApplicationAndUser", mock.Anything, mock.Anything, mock.Anything)
}

func TestInterviewProgressService_UpdateRejectsBlankIDBeforeRepositoryAccess(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	svc := appsvc.NewInterviewProgressService(repo, apps)

	got, err := svc.Update(context.Background(), "user-1", "  \t", dto.UpdateInterviewProgressRequest{})

	assert.Nil(t, got)
	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
	repo.AssertNotCalled(t, "FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything)
}

func TestInterviewProgressService_UpdateReturnsBadParamsWhenStatusUpdaterMissing(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	svc := appsvc.NewInterviewProgressService(repo, nil)
	occurredAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	existing := &domainprogress.InterviewProgress{
		ID:            "progress-1",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         domainprogress.StageViewed,
		OccurredAt:    occurredAt,
		Result:        domainprogress.ResultPassed,
	}

	repo.On("FindByIDAndUserID", mock.Anything, "progress-1", "user-1").Return(existing, nil).Once()
	repo.On("Update", mock.Anything, existing).Return(nil).Once()

	got, err := svc.Update(context.Background(), "user-1", "progress-1", dto.UpdateInterviewProgressRequest{})

	assert.Nil(t, got)
	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
	repo.AssertExpectations(t)
}

func TestInterviewProgressService_UpdatePreservesOmittedFieldsAppliesProvidedFieldsAndSyncsStatus(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	svc := appsvc.NewInterviewProgressService(repo, apps)
	oldOccurredAt := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)
	newOccurredAt := time.Date(2026, 5, 11, 11, 0, 0, 0, time.UTC)
	stage := domainprogress.StageSecondInterview
	note := "updated note"
	existing := &domainprogress.InterviewProgress{
		ID:                "progress-1",
		UserID:            "user-1",
		ApplicationID:     "app-1",
		Stage:             domainprogress.StageFirstInterview,
		OccurredAt:        oldOccurredAt,
		Result:            domainprogress.ResultPending,
		InterviewerOrTeam: "Engineering",
		Note:              "original note",
	}
	req := dto.UpdateInterviewProgressRequest{
		Stage:      &stage,
		OccurredAt: &dto.FlexibleTime{Time: newOccurredAt},
		Note:       &note,
	}

	repo.On("FindByIDAndUserID", mock.Anything, "progress-1", "user-1").Return(existing, nil).Once()
	repo.On("Update", mock.Anything, mock.MatchedBy(func(p *domainprogress.InterviewProgress) bool {
		return p.ID == "progress-1" && p.UserID == "user-1" && p.ApplicationID == "app-1" && p.Stage == stage && p.OccurredAt.Equal(newOccurredAt) && p.Result == domainprogress.ResultPending && p.InterviewerOrTeam == "Engineering" && p.Note == note
	})).Return(nil).Once()
	repo.On("ListByApplicationAndUser", mock.Anything, "app-1", "user-1").Return([]domainprogress.InterviewProgress{{
		ID:            "progress-1",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         stage,
		OccurredAt:    newOccurredAt,
		Result:        domainprogress.ResultPending,
	}}, nil).Once()
	apps.On("SyncCurrentStatus", mock.Anything, "user-1", "app-1", domainapp.StatusSecondInterview).Return(nil).Once()

	got, err := svc.Update(context.Background(), "user-1", "progress-1", req)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "progress-1", got.ID)
	assert.Equal(t, stage, got.Stage)
	repo.AssertExpectations(t)
	apps.AssertExpectations(t)
}

func TestInterviewProgressService_UpdateChangingApplicationIDSyncsBothOldAndNewApplicationIDs(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	svc := appsvc.NewInterviewProgressService(repo, apps)
	occurredAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	newApplicationID := "app-2"
	existing := &domainprogress.InterviewProgress{
		ID:            "progress-1",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         domainprogress.StageViewed,
		OccurredAt:    occurredAt,
		Result:        domainprogress.ResultPassed,
	}
	req := dto.UpdateInterviewProgressRequest{ApplicationID: &newApplicationID}

	repo.On("FindByIDAndUserID", mock.Anything, "progress-1", "user-1").Return(existing, nil).Once()
	apps.On("Get", mock.Anything, "user-1", "app-2").Return(&domainapp.ResumeApplication{ID: "app-2", UserID: "user-1"}, nil).Once()
	repo.On("Update", mock.Anything, mock.MatchedBy(func(p *domainprogress.InterviewProgress) bool {
		return p.ID == "progress-1" && p.ApplicationID == "app-2" && p.Stage == domainprogress.StageViewed && p.Result == domainprogress.ResultPassed
	})).Return(nil).Once()
	repo.On("ListByApplicationAndUser", mock.Anything, "app-1", "user-1").Return([]domainprogress.InterviewProgress{{
		ID:            "old-latest",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         domainprogress.StageApplied,
		OccurredAt:    occurredAt.Add(-time.Hour),
		Result:        domainprogress.ResultPassed,
	}}, nil).Once()
	repo.On("ListByApplicationAndUser", mock.Anything, "app-2", "user-1").Return([]domainprogress.InterviewProgress{{
		ID:            "progress-1",
		UserID:        "user-1",
		ApplicationID: "app-2",
		Stage:         domainprogress.StageViewed,
		OccurredAt:    occurredAt,
		Result:        domainprogress.ResultPassed,
	}}, nil).Once()
	apps.On("SyncCurrentStatus", mock.Anything, "user-1", "app-1", domainapp.StatusApplied).Return(nil).Once()
	apps.On("SyncCurrentStatus", mock.Anything, "user-1", "app-2", domainapp.StatusViewed).Return(nil).Once()

	got, err := svc.Update(context.Background(), "user-1", "progress-1", req)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "app-2", got.ApplicationID)
	repo.AssertExpectations(t)
	apps.AssertExpectations(t)
}

func TestInterviewProgressService_UpdateChangingApplicationIDRejectsUnownedOrMissingNewApplicationBeforeUpdate(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	svc := appsvc.NewInterviewProgressService(repo, apps)
	newApplicationID := "app-2"
	existing := &domainprogress.InterviewProgress{
		ID:            "progress-1",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         domainprogress.StageViewed,
		OccurredAt:    time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC),
		Result:        domainprogress.ResultPassed,
	}
	repo.On("FindByIDAndUserID", mock.Anything, "progress-1", "user-1").Return(existing, nil).Once()
	apps.On("Get", mock.Anything, "user-1", "app-2").Return(nil, domainerrors.ErrNotFound).Once()

	got, err := svc.Update(context.Background(), "user-1", "progress-1", dto.UpdateInterviewProgressRequest{ApplicationID: &newApplicationID})

	assert.Nil(t, got)
	assert.ErrorIs(t, err, domainerrors.ErrNotFound)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	apps.AssertExpectations(t)
}

func TestInterviewProgressService_UpdateMovingOnlyProgressAwaySyncsOldApplicationToApplied(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	svc := appsvc.NewInterviewProgressService(repo, apps)
	occurredAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	newApplicationID := "app-2"
	existing := &domainprogress.InterviewProgress{
		ID:            "progress-1",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         domainprogress.StageFinalInterview,
		OccurredAt:    occurredAt,
		Result:        domainprogress.ResultPassed,
	}
	req := dto.UpdateInterviewProgressRequest{ApplicationID: &newApplicationID}

	repo.On("FindByIDAndUserID", mock.Anything, "progress-1", "user-1").Return(existing, nil).Once()
	apps.On("Get", mock.Anything, "user-1", "app-2").Return(&domainapp.ResumeApplication{ID: "app-2", UserID: "user-1"}, nil).Once()
	repo.On("Update", mock.Anything, mock.MatchedBy(func(p *domainprogress.InterviewProgress) bool {
		return p.ID == "progress-1" && p.ApplicationID == "app-2"
	})).Return(nil).Once()
	repo.On("ListByApplicationAndUser", mock.Anything, "app-1", "user-1").Return([]domainprogress.InterviewProgress{}, nil).Once()
	repo.On("ListByApplicationAndUser", mock.Anything, "app-2", "user-1").Return([]domainprogress.InterviewProgress{*existing}, nil).Once()
	apps.On("SyncCurrentStatus", mock.Anything, "user-1", "app-1", domainapp.StatusApplied).Return(nil).Once()
	apps.On("SyncCurrentStatus", mock.Anything, "user-1", "app-2", domainapp.StatusFinalInterview).Return(nil).Once()

	got, err := svc.Update(context.Background(), "user-1", "progress-1", req)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "app-2", got.ApplicationID)
	repo.AssertExpectations(t)
	apps.AssertExpectations(t)
}

func TestInterviewProgressService_DeleteRejectsBlankIDBeforeRepositoryAccess(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	svc := appsvc.NewInterviewProgressService(repo, apps)

	err := svc.Delete(context.Background(), "user-1", "  \n")

	assert.ErrorIs(t, err, domainerrors.ErrBadParams)
	repo.AssertNotCalled(t, "FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything)
}

func TestInterviewProgressService_DeleteSyncsOldApplicationStatus(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	svc := appsvc.NewInterviewProgressService(repo, apps)
	occurredAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	existing := &domainprogress.InterviewProgress{
		ID:            "progress-1",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         domainprogress.StageFinalInterview,
		OccurredAt:    occurredAt,
		Result:        domainprogress.ResultPassed,
	}

	repo.On("FindByIDAndUserID", mock.Anything, "progress-1", "user-1").Return(existing, nil).Once()
	repo.On("Delete", mock.Anything, "progress-1", "user-1").Return(nil).Once()
	repo.On("ListByApplicationAndUser", mock.Anything, "app-1", "user-1").Return([]domainprogress.InterviewProgress{{
		ID:            "remaining",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         domainprogress.StageSecondInterview,
		OccurredAt:    occurredAt.Add(-time.Hour),
		Result:        domainprogress.ResultPassed,
	}}, nil).Once()
	apps.On("SyncCurrentStatus", mock.Anything, "user-1", "app-1", domainapp.StatusSecondInterview).Return(nil).Once()

	err := svc.Delete(context.Background(), "user-1", "progress-1")

	require.NoError(t, err)
	repo.AssertExpectations(t)
	apps.AssertExpectations(t)
}

func TestInterviewProgressService_DeleteOnlyProgressSyncsApplicationToApplied(t *testing.T) {
	repo := &mockInterviewProgressRepo{}
	apps := &mockApplicationService{}
	svc := appsvc.NewInterviewProgressService(repo, apps)
	occurredAt := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	existing := &domainprogress.InterviewProgress{
		ID:            "progress-1",
		UserID:        "user-1",
		ApplicationID: "app-1",
		Stage:         domainprogress.StageFinalInterview,
		OccurredAt:    occurredAt,
		Result:        domainprogress.ResultPassed,
	}

	repo.On("FindByIDAndUserID", mock.Anything, "progress-1", "user-1").Return(existing, nil).Once()
	repo.On("Delete", mock.Anything, "progress-1", "user-1").Return(nil).Once()
	repo.On("ListByApplicationAndUser", mock.Anything, "app-1", "user-1").Return([]domainprogress.InterviewProgress{}, nil).Once()
	apps.On("SyncCurrentStatus", mock.Anything, "user-1", "app-1", domainapp.StatusApplied).Return(nil).Once()

	err := svc.Delete(context.Background(), "user-1", "progress-1")

	require.NoError(t, err)
	repo.AssertExpectations(t)
	apps.AssertExpectations(t)
}
