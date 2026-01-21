package event

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	domainbill "billsplitter-monolith/internal/domain/bill"
	domainevent "billsplitter-monolith/internal/domain/event"
	"billsplitter-monolith/internal/domain/session"
	domainuser "billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	appErrors "billsplitter-monolith/internal/errors"
	billmock "billsplitter-monolith/internal/mocks/domain/bill"
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
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
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
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
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
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
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
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
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
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
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
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
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
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
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
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
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
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
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

func TestUseCase_CalculateSummary(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)

		var (
			meetID   int64 = 1
			andreyID int64 = 1
			nastyID  int64 = 2
			artemID  int64 = 3
			danikID  int64 = 4

			ctx = middleware.InjectUserInSession(context.Background(), &session.Session{UserID: vo.UserID(andreyID)})
		)

		expectedEvent := &domainevent.Event{
			ID: meetID,
			Members: []domainevent.Member{
				{
					ID:     andreyID,
					Name:   "Андрей",
					UserID: utils.Ptr[vo.UserID](vo.UserID(andreyID)),
				},
				{
					ID:     nastyID,
					Name:   "Настя",
					UserID: utils.Ptr[vo.UserID](vo.UserID(nastyID)),
				},
				{
					ID:     artemID,
					Name:   "Артем",
					UserID: utils.Ptr[vo.UserID](vo.UserID(artemID)),
				},
				{
					ID:     danikID,
					Name:   "Даник",
					UserID: utils.Ptr[vo.UserID](vo.UserID(danikID)),
				},
			},
		}

		expectedBills := []domainbill.Bill{
			{
				Name:        "Калик",
				EventID:     meetID,
				PaidBy:      andreyID,
				TotalAmount: 60,
				Participants: []domainbill.Participant{
					{
						MemberID: andreyID,
						Amount:   20,
					},
					{
						MemberID: nastyID,
						Amount:   20,
					},
					{
						MemberID: danikID,
						Amount:   20,
					},
				},
			},
			{
				Name:        "Еда",
				EventID:     meetID,
				PaidBy:      danikID,
				TotalAmount: 100,
				Participants: []domainbill.Participant{
					{
						MemberID: andreyID,
						Amount:   25,
					},
					{
						MemberID: nastyID,
						Amount:   25,
					},
					{
						MemberID: artemID,
						Amount:   25,
					},
					{
						MemberID: danikID,
						Amount:   25,
					},
				},
			},
			{
				Name:        "Такси",
				EventID:     meetID,
				PaidBy:      artemID,
				TotalAmount: 20,
				Participants: []domainbill.Participant{
					{
						MemberID: andreyID,
						Amount:   5,
					},
					{
						MemberID: nastyID,
						Amount:   5,
					},
					{
						MemberID: artemID,
						Amount:   5,
					},
					{
						MemberID: danikID,
						Amount:   5,
					},
				},
			},
		}

		eventSvc.EXPECT().
			GetByID(mock.Anything, meetID).
			Return(expectedEvent, nil)

		billSvc.EXPECT().
			FetchByEventID(mock.Anything, meetID).
			Return(expectedBills, nil)

		summary, err := uc.CalculateSummary(ctx, meetID)
		require.NoError(t, err)
		require.NotNil(t, summary)

		expectedBalances := []Balance{
			{
				MemberID:   andreyID,
				UserID:     utils.Ptr[vo.UserID](vo.UserID(andreyID)),
				Name:       "Андрей",
				TotalPaid:  60,
				TotalShare: 50,
				Balance:    10,
			},
			{
				MemberID:   nastyID,
				UserID:     utils.Ptr[vo.UserID](vo.UserID(nastyID)),
				Name:       "Настя",
				TotalPaid:  0,
				TotalShare: 50,
				Balance:    -50,
			},
			{
				MemberID:   artemID,
				UserID:     utils.Ptr[vo.UserID](vo.UserID(artemID)),
				Name:       "Артем",
				TotalPaid:  20,
				TotalShare: 30,
				Balance:    -10,
			},
			{
				MemberID:   danikID,
				UserID:     utils.Ptr[vo.UserID](vo.UserID(danikID)),
				Name:       "Даник",
				TotalPaid:  100,
				TotalShare: 50,
				Balance:    50,
			},
		}

		expectedSettlements := []Settlement{
			{
				FromMemberID: artemID,
				ToMemberID:   andreyID,
				Amount:       10,
			},
			{
				FromMemberID: nastyID,
				ToMemberID:   danikID,
				Amount:       50,
			},
		}

		require.Equal(t, expectedBalances, summary.Balances)
		require.Equal(t, expectedSettlements, summary.Settlements)
	})

	t.Run("get meet error", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
		ctx := context.Background()
		meetID := int64(42)
		svcErr := stderrors.New("get failed")

		eventSvc.EXPECT().
			GetByID(mock.Anything, meetID).
			Return(nil, svcErr)

		summary, err := uc.CalculateSummary(ctx, meetID)

		require.Nil(t, summary)
		require.ErrorIs(t, err, svcErr)
	})

	t.Run("meet not found", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
		ctx := context.Background()
		meetID := int64(100)

		eventSvc.EXPECT().
			GetByID(mock.Anything, meetID).
			Return(nil, nil)

		summary, err := uc.CalculateSummary(ctx, meetID)

		require.Nil(t, summary)
		require.ErrorIs(t, err, appErrors.ErrEventNotFound)
	})

	t.Run("bill fetch error", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
		ctx := middleware.InjectUserInSession(context.Background(), &session.Session{UserID: vo.UserID(5)})
		meetID := int64(7)
		svcErr := stderrors.New("bill fetch failed")

		eventSvc.EXPECT().
			GetByID(mock.Anything, meetID).
			Return(&domainevent.Event{
				ID: meetID,
				Members: []domainevent.Member{
					{
						ID:     5,
						Name:   "A",
						UserID: utils.Ptr[vo.UserID](vo.UserID(5)),
					},
				},
			}, nil)

		billSvc.EXPECT().
			FetchByEventID(mock.Anything, meetID).
			Return(nil, svcErr)

		summary, err := uc.CalculateSummary(ctx, meetID)

		require.Nil(t, summary)
		require.ErrorIs(t, err, svcErr)
	})

	t.Run("no bills", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
		ctx := middleware.InjectUserInSession(context.Background(), &session.Session{UserID: vo.UserID(1)})
		meetID := int64(8)

		me := &domainevent.Event{
			ID: meetID,
			Members: []domainevent.Member{
				{
					ID:     1,
					Name:   "One",
					UserID: utils.Ptr[vo.UserID](vo.UserID(1)),
				},
				{
					ID:     2,
					Name:   "Two",
					UserID: utils.Ptr[vo.UserID](vo.UserID(2)),
				},
			},
		}

		eventSvc.EXPECT().
			GetByID(mock.Anything, meetID).
			Return(me, nil)

		billSvc.EXPECT().
			FetchByEventID(mock.Anything, meetID).
			Return(nil, nil)

		summary, err := uc.CalculateSummary(ctx, meetID)

		require.NoError(t, err)
		require.NotNil(t, summary)
		require.Len(t, summary.Balances, len(me.Members))
		for idx, balance := range summary.Balances {
			require.Zero(t, balance.TotalPaid)
			require.Zero(t, balance.TotalShare)
			require.Zero(t, balance.Balance)
			require.Equal(t, me.Members[idx].UserID, balance.UserID)
		}
		require.Empty(t, summary.Settlements)
	})
}

func TestUseCase_AssignMemberToUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
		ctx := context.Background()

		meetID := int64(10)
		memberID := int64(5)
		userID := vo.UserID(7)
		meet := &domainevent.Event{
			ID: meetID,
			Members: []domainevent.Member{
				{ID: memberID},
			},
		}

		eventSvc.EXPECT().
			GetByID(mock.Anything, meetID).
			Return(meet, nil)
		eventSvc.EXPECT().
			AssignMemberUser(mock.Anything, meetID, memberID, userID).
			Return(nil)

		err := uc.AssignMemberToUser(ctx, meetID, memberID, userID)

		require.NoError(t, err)
	})

	t.Run("meet id required", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)

		err := uc.AssignMemberToUser(context.Background(), 0, 1, vo.UserID(1))

		require.ErrorContains(t, err, "validation error")
	})

	t.Run("member id required", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)

		err := uc.AssignMemberToUser(context.Background(), 1, 0, vo.UserID(1))

		require.ErrorContains(t, err, "validation error")
	})

	t.Run("event not found", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
		ctx := context.Background()

		eventSvc.EXPECT().
			GetByID(mock.Anything, int64(1)).
			Return(nil, nil)

		err := uc.AssignMemberToUser(ctx, 1, 2, vo.UserID(3))

		require.ErrorIs(t, err, appErrors.ErrEventNotFound)
	})

	t.Run("member not in event", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
		ctx := context.Background()

		meet := &domainevent.Event{
			ID: 1,
			Members: []domainevent.Member{
				{ID: 10},
			},
		}

		eventSvc.EXPECT().
			GetByID(mock.Anything, int64(1)).
			Return(meet, nil)

		err := uc.AssignMemberToUser(ctx, 1, 2, vo.UserID(3))

		require.ErrorContains(t, err, "validation error")
	})

	t.Run("member already assigned", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
		ctx := context.Background()

		meet := &domainevent.Event{
			ID: 1,
			Members: []domainevent.Member{
				{ID: 10, UserID: utils.Ptr(vo.UserID(5))},
			},
		}

		eventSvc.EXPECT().
			GetByID(mock.Anything, int64(1)).
			Return(meet, nil)

		err := uc.AssignMemberToUser(ctx, 1, 10, vo.UserID(3))

		require.ErrorContains(t, err, "validation error")
	})

	t.Run("assign error", func(t *testing.T) {
		eventSvc := eventmock.NewMockService(t)
		userSvc := usermock.NewMockService(t)
		billSvc := billmock.NewMockService(t)
		uc := New(eventSvc, userSvc, billSvc)
		ctx := context.Background()

		meet := &domainevent.Event{
			ID: 1,
			Members: []domainevent.Member{
				{ID: 10},
			},
		}
		assignErr := stderrors.New("assign failed")

		eventSvc.EXPECT().
			GetByID(mock.Anything, int64(1)).
			Return(meet, nil)
		eventSvc.EXPECT().
			AssignMemberUser(mock.Anything, int64(1), int64(10), vo.UserID(3)).
			Return(assignErr)

		err := uc.AssignMemberToUser(ctx, 1, 10, vo.UserID(3))

		require.ErrorIs(t, err, assignErr)
	})
}
