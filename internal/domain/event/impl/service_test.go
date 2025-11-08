package impl

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"billsplitter-monolith/internal/domain/event"
	vo "billsplitter-monolith/internal/domain/valueobject"
	eventmock "billsplitter-monolith/internal/mocks/domain/event"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := eventmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()

		eventDate := time.Now().Add(24 * time.Hour)
		creatorID := vo.UserID(100)
		req := event.CreateEventRq{
			Name:            "Birthday",
			EventDate:       &eventDate,
			CreatorUserID:   creatorID,
			CreatorUsername: "Alice",
			Members:         []string{"Bob", "Charlie"},
		}

		expectedEvent := event.Event{
			Name:            req.Name,
			CreatedByUserID: req.CreatorUserID,
			Status:          event.StatusDraft,
			Type:            event.TypeMeet,
			EventDate:       req.EventDate,
		}
		createdID := int64(77)

		repo.EXPECT().
			Create(mock.Anything, mock.MatchedBy(func(ev event.Event) bool {
				require.Equal(t, expectedEvent, ev)
				return true
			})).
			Return(createdID, nil)

		expectedMembers := []event.Member{
			{
				UserID: &req.CreatorUserID,
				Name:   req.CreatorUsername,
			},
			{
				UserID: nil,
				Name:   "Bob",
			},
			{
				UserID: nil,
				Name:   "Charlie",
			},
		}

		repo.EXPECT().
			AddMembers(mock.Anything, createdID, mock.MatchedBy(func(members []event.Member) bool {
				require.Equal(t, expectedMembers, members)
				return true
			})).
			Return(nil)

		id, err := svc.Create(ctx, req)

		require.NoError(t, err)
		require.Equal(t, createdID, id)
	})

	t.Run("create error", func(t *testing.T) {
		repo := eventmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		req := event.CreateEventRq{
			Name:            "Test",
			CreatorUserID:   vo.UserID(1),
			CreatorUsername: "Tester",
		}
		repoErr := stderrors.New("create failed")

		repo.EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(int64(0), repoErr)

		id, err := svc.Create(ctx, req)

		require.Zero(t, id)
		require.ErrorIs(t, err, repoErr)
	})

	t.Run("add members error", func(t *testing.T) {
		repo := eventmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		req := event.CreateEventRq{
			Name:            "Test",
			CreatorUserID:   vo.UserID(1),
			CreatorUsername: "Tester",
			Members:         []string{"One"},
		}
		repoErr := stderrors.New("add failed")
		createdID := int64(10)

		repo.EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(createdID, nil)

		repo.EXPECT().
			AddMembers(mock.Anything, createdID, mock.Anything).
			Return(repoErr)

		id, err := svc.Create(ctx, req)

		require.Zero(t, id)
		require.ErrorIs(t, err, repoErr)
	})
}

func TestService_FetchByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := eventmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userID := vo.UserID(55)
		expected := []event.Event{{ID: 1}, {ID: 2}}

		repo.EXPECT().
			FetchByUserID(mock.Anything, userID).
			Return(expected, nil)

		events, err := svc.FetchByUserID(ctx, userID)

		require.NoError(t, err)
		require.Equal(t, expected, events)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := eventmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userID := vo.UserID(55)
		repoErr := stderrors.New("fetch failed")

		repo.EXPECT().
			FetchByUserID(mock.Anything, userID).
			Return(nil, repoErr)

		events, err := svc.FetchByUserID(ctx, userID)

		require.Nil(t, events)
		require.ErrorIs(t, err, repoErr)
	})
}

func TestService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := eventmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		eventID := int64(77)
		expected := &event.Event{ID: eventID}

		repo.EXPECT().
			GetByID(mock.Anything, eventID).
			Return(expected, nil)

		ev, err := svc.GetByID(ctx, eventID)

		require.NoError(t, err)
		require.Equal(t, expected, ev)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := eventmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		eventID := int64(77)
		repoErr := stderrors.New("get failed")

		repo.EXPECT().
			GetByID(mock.Anything, eventID).
			Return(nil, repoErr)

		ev, err := svc.GetByID(ctx, eventID)

		require.Nil(t, ev)
		require.ErrorIs(t, err, repoErr)
	})
}
