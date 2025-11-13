package bill

import (
	"context"
	"errors"
	"testing"

	domainbill "billsplitter-monolith/internal/domain/bill"
	domainevent "billsplitter-monolith/internal/domain/event"
	"billsplitter-monolith/internal/domain/session"
	vo "billsplitter-monolith/internal/domain/valueobject"
	apperrors "billsplitter-monolith/internal/errors"
	billmock "billsplitter-monolith/internal/mocks/domain/bill"
	eventmock "billsplitter-monolith/internal/mocks/domain/event"
	"billsplitter-monolith/internal/transport/http/middleware"
	"billsplitter-monolith/internal/utils"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errBillTest = errors.New("bill-test-error")
var testCreatedByID = utils.Ptr[vo.UserID](1)

func TestUseCase_CreateBill(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)

		ev := &domainevent.Event{
			ID: 1,
			Members: []domainevent.Member{
				{ID: 10, UserID: testCreatedByID}, {ID: 11}},
		}
		eventSvc.EXPECT().GetByID(mock.Anything, int64(1)).Return(ev, nil)

		rq := CreateBillRq{
			EventID:     1,
			Name:        "Dinner",
			CreatedBy:   *testCreatedByID,
			TotalAmount: vo.Amount(1000),
			Currency:    vo.CurrencyCode("BYN"),
			SplitType:   vo.SplitTypeEven,
			PaidBy:      10,
			Participants: []domainbill.Participant{
				{MemberID: 10, Amount: 600},
				{MemberID: 11, Amount: 400},
			},
		}

		billSvc.EXPECT().Create(mock.Anything, mock.MatchedBy(func(b domainbill.Bill) bool {
			return b.EventID == rq.EventID && b.PaidBy == rq.PaidBy && len(b.Participants) == 2 && b.Participants[0].MemberID == 10
		})).Return(int64(99), nil)

		id, err := uc.CreateBill(context.Background(), rq)
		require.NoError(t, err)
		require.Equal(t, int64(99), id)
	})

	t.Run("event not found", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)

		eventSvc.EXPECT().GetByID(mock.Anything, int64(1)).Return(nil, nil)

		_, err := uc.CreateBill(context.Background(), CreateBillRq{
			EventID:     1,
			TotalAmount: vo.Amount(100),
			Currency:    vo.CurrencyCode("BYN"),
			SplitType:   vo.SplitTypeEven,
			PaidBy:      10,
			Participants: []domainbill.Participant{
				{MemberID: 10, Amount: 100},
			},
		})
		require.ErrorIs(t, err, apperrors.ErrEventNotFound)
	})

	t.Run("invalid participant", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)

		ev := &domainevent.Event{
			ID:      1,
			Members: []domainevent.Member{{ID: 10, UserID: testCreatedByID}},
		}
		eventSvc.EXPECT().GetByID(mock.Anything, int64(1)).Return(ev, nil)

		_, err := uc.CreateBill(context.Background(), CreateBillRq{
			EventID:     1,
			TotalAmount: vo.Amount(80),
			Currency:    vo.CurrencyCode("BYN"),
			SplitType:   vo.SplitTypeEven,
			PaidBy:      999,
			Participants: []domainbill.Participant{
				{MemberID: 999, Amount: 80},
			},
			CreatedBy: *testCreatedByID,
		})
		require.Error(t, err)
	})

	t.Run("invalid paid_by", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)

		ev := &domainevent.Event{
			ID:      1,
			Members: []domainevent.Member{{ID: 10, UserID: testCreatedByID}},
		}
		eventSvc.EXPECT().GetByID(mock.Anything, int64(1)).Return(ev, nil)

		_, err := uc.CreateBill(context.Background(), CreateBillRq{
			EventID:     1,
			TotalAmount: vo.Amount(80),
			Currency:    vo.CurrencyCode("BYN"),
			SplitType:   vo.SplitTypeEven,
			PaidBy:      11,
			Participants: []domainbill.Participant{
				{MemberID: 10, Amount: 80},
			},
			CreatedBy: *testCreatedByID,
		})
		require.Error(t, err)
	})

	t.Run("invalid created_by", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)

		ev := &domainevent.Event{
			ID: 1,
			Members: []domainevent.Member{
				{
					ID: 10,
				},
				{
					ID:     11,
					UserID: utils.Ptr[vo.UserID](33),
				},
			},
		}
		eventSvc.EXPECT().GetByID(mock.Anything, int64(1)).Return(ev, nil)

		_, err := uc.CreateBill(context.Background(), CreateBillRq{
			EventID:     1,
			TotalAmount: vo.Amount(80),
			Currency:    vo.CurrencyCode("BYN"),
			SplitType:   vo.SplitTypeEven,
			PaidBy:      11,
			Participants: []domainbill.Participant{
				{MemberID: 10, Amount: 80},
			},
			CreatedBy: 34, // не принадлежит этому ивенту - ошибка валидации
		})
		require.Error(t, err)
	})

	t.Run("bill service error", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)

		ev := &domainevent.Event{
			ID:      1,
			Members: []domainevent.Member{{ID: 10, UserID: testCreatedByID}},
		}
		eventSvc.EXPECT().GetByID(mock.Anything, int64(1)).Return(ev, nil)

		billSvc.EXPECT().Create(mock.Anything, mock.Anything).Return(int64(0), errBillTest)

		_, err := uc.CreateBill(context.Background(), CreateBillRq{
			EventID:     1,
			TotalAmount: vo.Amount(50),
			Currency:    vo.CurrencyCode("BYN"),
			SplitType:   vo.SplitTypeEven,
			PaidBy:      10,
			Participants: []domainbill.Participant{
				{MemberID: 10, Amount: 50},
			},
			CreatedBy: *testCreatedByID,
		})
		require.ErrorIs(t, err, errBillTest)
	})
}

func TestUseCase_FetchEventBills(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)
		ctx := middleware.InjectUserInSession(context.Background(), &session.Session{UserID: *testCreatedByID})

		expectedEv := &domainevent.Event{
			ID: 1,
			Members: []domainevent.Member{
				{
					UserID: testCreatedByID,
				},
			},
		}

		eventSvc.EXPECT().GetByID(mock.Anything, int64(1)).Return(expectedEv, nil)
		expected := []domainbill.Bill{{ID: 1}}
		billSvc.EXPECT().FetchByEventID(mock.Anything, int64(1)).Return(expected, nil)

		bills, err := uc.FetchEventBills(ctx, 1)
		require.NoError(t, err)
		require.Equal(t, expected, bills)
	})

	t.Run("event missing", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)
		ctx := middleware.InjectUserInSession(context.Background(), &session.Session{UserID: *testCreatedByID})

		eventSvc.EXPECT().GetByID(mock.Anything, int64(2)).Return(nil, nil)

		_, err := uc.FetchEventBills(ctx, 2)
		require.ErrorIs(t, err, apperrors.ErrEventNotFound)
	})

	t.Run("bill service error", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)
		ctx := middleware.InjectUserInSession(context.Background(), &session.Session{UserID: *testCreatedByID})

		expectedEv := &domainevent.Event{
			ID: 1,
			Members: []domainevent.Member{
				{
					UserID: testCreatedByID,
				},
			},
		}

		eventSvc.EXPECT().GetByID(mock.Anything, int64(3)).Return(expectedEv, nil)
		billSvc.EXPECT().FetchByEventID(mock.Anything, int64(3)).Return(nil, errBillTest)

		_, err := uc.FetchEventBills(ctx, 3)
		require.ErrorIs(t, err, errBillTest)
	})
}
