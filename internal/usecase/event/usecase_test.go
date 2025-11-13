package event

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	domainevent "billsplitter-monolith/internal/domain/event"
	"billsplitter-monolith/internal/domain/session"
	domainuser "billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	appErrors "billsplitter-monolith/internal/errors"
	eventmock "billsplitter-monolith/internal/mocks/domain/event"
	usermock "billsplitter-monolith/internal/mocks/domain/user"
	"billsplitter-monolith/internal/transport/http/middleware"
	"billsplitter-monolith/internal/utils"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var testUserID = utils.Ptr[vo.UserID](1)

func TestUseCase_CreateMeet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		uc := New(eventSvc, userSvc)
		ctx := context.Background()

		date := time.Now().Add(24 * time.Hour)
		creatorID := vo.UserID(10)
		creator := &domainuser.User{ID: creatorID, Username: "Alice"}
		req := CreateMeetRq{
			EventName:       "Birthday",
			Date:            &date,
			CreatedByUserID: creatorID,
			Members:         []string{"Bob", "Charlie"},
		}

		userSvc.EXPECT().
			GetByID(mock.Anything, creatorID).
			Return(creator, nil)

		expectedCreate := domainevent.CreateEventRq{
			Name:            req.EventName,
			EventDate:       req.Date,
			CreatorUserID:   creator.ID,
			CreatorUsername: creator.Username,
			Members:         req.Members,
		}
		createdMeetID := int64(77)

		eventSvc.EXPECT().
			Create(mock.Anything, expectedCreate).
			Return(createdMeetID, nil)

		id, err := uc.CreateMeet(ctx, req)

		require.NoError(t, err)
		require.Equal(t, createdMeetID, id)
	})

	t.Run("user service error", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		uc := New(eventSvc, userSvc)
		ctx := context.Background()
		req := CreateMeetRq{CreatedByUserID: vo.UserID(1)}
		svcErr := stderrors.New("user svc failure")

		userSvc.EXPECT().
			GetByID(mock.Anything, req.CreatedByUserID).
			Return(nil, svcErr)

		id, err := uc.CreateMeet(ctx, req)

		require.Zero(t, id)
		require.ErrorIs(t, err, svcErr)
	})

	t.Run("creator not found", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		uc := New(eventSvc, userSvc)
		ctx := context.Background()
		req := CreateMeetRq{CreatedByUserID: vo.UserID(1)}

		userSvc.EXPECT().
			GetByID(mock.Anything, req.CreatedByUserID).
			Return(nil, nil)

		id, err := uc.CreateMeet(ctx, req)

		require.Zero(t, id)
		require.ErrorContains(t, err, "get event creator error")
		require.ErrorIs(t, err, appErrors.ErrUserNotFound)
	})

	t.Run("event service error", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		uc := New(eventSvc, userSvc)
		ctx := context.Background()
		creatorID := vo.UserID(2)
		creator := &domainuser.User{ID: creatorID, Username: "Alice"}
		req := CreateMeetRq{EventName: "Meet", CreatedByUserID: creatorID}
		svcErr := stderrors.New("create failed")

		userSvc.EXPECT().
			GetByID(mock.Anything, creatorID).
			Return(creator, nil)

		eventSvc.EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(int64(0), svcErr)

		id, err := uc.CreateMeet(ctx, req)

		require.Zero(t, id)
		require.ErrorIs(t, err, svcErr)
	})
}

func TestUseCase_FetchUserMeets(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		uc := New(eventSvc, userSvc)
		ctx := context.Background()
		userID := vo.UserID(5)
		expected := []domainevent.Event{{ID: 1}}

		eventSvc.EXPECT().
			FetchByUserID(mock.Anything, userID).
			Return(expected, nil)

		events, err := uc.FetchUserMeets(ctx, userID)

		require.NoError(t, err)
		require.Equal(t, expected, events)
	})

	t.Run("event service error", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		uc := New(eventSvc, userSvc)
		ctx := context.Background()
		userID := vo.UserID(5)
		svcErr := stderrors.New("fetch failed")

		eventSvc.EXPECT().
			FetchByUserID(mock.Anything, userID).
			Return(nil, svcErr)

		events, err := uc.FetchUserMeets(ctx, userID)

		require.Nil(t, events)
		require.ErrorIs(t, err, svcErr)
	})
}

func TestUseCase_GetMeetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		uc := New(eventSvc, userSvc)
		ctx := middleware.InjectUserInSession(context.Background(), &session.Session{UserID: *testUserID})

		meetID := int64(42)
		expected := &domainevent.Event{
			ID: meetID,
			Members: []domainevent.Member{
				{
					UserID: testUserID,
				},
			},
		}

		eventSvc.EXPECT().
			GetByID(mock.Anything, meetID).
			Return(expected, nil)

		meet, err := uc.GetMeetByID(ctx, meetID)

		require.NoError(t, err)
		require.Equal(t, expected, meet)
	})

	t.Run("forbidden event", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		uc := New(eventSvc, userSvc)
		ctx := middleware.InjectUserInSession(context.Background(), &session.Session{UserID: *testUserID})
		meetID := int64(42)
		expected := &domainevent.Event{
			ID: meetID,
			Members: []domainevent.Member{
				{
					UserID: testUserID,
				},
			},
		}

		eventSvc.EXPECT().
			GetByID(mock.Anything, meetID).
			Return(expected, nil)

		meet, err := uc.GetMeetByID(ctx, meetID)

		require.NoError(t, err)
		require.Equal(t, expected, meet)
	})

	t.Run("event service error", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		uc := New(eventSvc, userSvc)
		ctx := context.Background()
		meetID := int64(42)
		svcErr := stderrors.New("get failed")

		eventSvc.EXPECT().
			GetByID(mock.Anything, meetID).
			Return(nil, svcErr)

		meet, err := uc.GetMeetByID(ctx, meetID)

		require.Nil(t, meet)
		require.ErrorIs(t, err, svcErr)
	})
}
