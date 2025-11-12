package bill

import (
	"context"
	"errors"
	"testing"

	domainbill "billsplitter-monolith/internal/domain/bill"
	domainevent "billsplitter-monolith/internal/domain/event"
	vo "billsplitter-monolith/internal/domain/valueobject"
	apperrors "billsplitter-monolith/internal/errors"
	billmock "billsplitter-monolith/internal/mocks/domain/bill"
	eventmock "billsplitter-monolith/internal/mocks/domain/event"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errBillTest = errors.New("bill-test-error")

func TestUseCase_CreateBill(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)

		ev := &domainevent.Event{
			ID:      1,
			Members: []domainevent.Member{{ID: 10}, {ID: 11}},
		}
		eventSvc.EXPECT().GetByID(mock.Anything, int64(1)).Return(ev, nil)

		rq := CreateBillRq{
			EventID:     1,
			Name:        "Dinner",
			CreatedBy:   vo.UserID(5),
			TotalAmount: vo.Amount(1000),
			Currency:    vo.CurrencyCode("BYN"),
			SplitType:   vo.SplitTypeEven,
			Participants: []domainbill.Participant{
				{MemberID: 10, Amount: 600},
				{MemberID: 11, Amount: 400},
			},
		}

		billSvc.EXPECT().Create(mock.Anything, mock.MatchedBy(func(b domainbill.Bill) bool {
			return b.EventID == rq.EventID && len(b.Participants) == 2 && b.Participants[0].MemberID == 10
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
			Members: []domainevent.Member{{ID: 10}},
		}
		eventSvc.EXPECT().GetByID(mock.Anything, int64(1)).Return(ev, nil)

		_, err := uc.CreateBill(context.Background(), CreateBillRq{
			EventID:     1,
			TotalAmount: vo.Amount(80),
			Currency:    vo.CurrencyCode("BYN"),
			SplitType:   vo.SplitTypeEven,
			Participants: []domainbill.Participant{
				{MemberID: 999, Amount: 80},
			},
		})
		require.Error(t, err)
	})

	t.Run("bill service error", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)

		ev := &domainevent.Event{
			ID:      1,
			Members: []domainevent.Member{{ID: 10}},
		}
		eventSvc.EXPECT().GetByID(mock.Anything, int64(1)).Return(ev, nil)

		billSvc.EXPECT().Create(mock.Anything, mock.Anything).Return(int64(0), errBillTest)

		_, err := uc.CreateBill(context.Background(), CreateBillRq{
			EventID:     1,
			TotalAmount: vo.Amount(50),
			Currency:    vo.CurrencyCode("BYN"),
			SplitType:   vo.SplitTypeEven,
			Participants: []domainbill.Participant{
				{MemberID: 10, Amount: 50},
			},
		})
		require.ErrorIs(t, err, errBillTest)
	})
}

func TestUseCase_FetchEventBills(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)

		eventSvc.EXPECT().GetByID(mock.Anything, int64(1)).Return(&domainevent.Event{ID: 1}, nil)
		expected := []domainbill.Bill{{ID: 1}}
		billSvc.EXPECT().FetchByEventID(mock.Anything, int64(1)).Return(expected, nil)

		bills, err := uc.FetchEventBills(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, expected, bills)
	})

	t.Run("event missing", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)

		eventSvc.EXPECT().GetByID(mock.Anything, int64(2)).Return(nil, nil)

		_, err := uc.FetchEventBills(context.Background(), 2)
		require.ErrorIs(t, err, apperrors.ErrEventNotFound)
	})

	t.Run("bill service error", func(t *testing.T) {
		billSvc := billmock.NewMockService(t)
		eventSvc := eventmock.NewMockService(t)
		uc := New(billSvc, eventSvc)

		eventSvc.EXPECT().GetByID(mock.Anything, int64(3)).Return(&domainevent.Event{ID: 3}, nil)
		billSvc.EXPECT().FetchByEventID(mock.Anything, int64(3)).Return(nil, errBillTest)

		_, err := uc.FetchEventBills(context.Background(), 3)
		require.ErrorIs(t, err, errBillTest)
	})
}
