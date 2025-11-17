//go:build e2e
// +build e2e

package e2e

import (
	"billsplitter-monolith/internal/e2e/sdk"
	"billsplitter-monolith/internal/transport/http/auth"
	"billsplitter-monolith/internal/transport/http/bill"
	"billsplitter-monolith/internal/transport/http/meet"
	"billsplitter-monolith/internal/transport/http/user"
	"context"
	"fmt"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL = "http://localhost:5001"
)

func TestE2E_HappyPath(t *testing.T) {
	t.Run("login and me - success", func(t *testing.T) {
		sdkClient := sdk.New(baseURL)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		unix := time.Now().Unix()

		loginReq := auth.LoginTelegramReq{
			Username:   fmt.Sprintf("test_user_%d", unix),
			FirstName:  "Test",
			LastName:   fmt.Sprintf("User_%d", unix),
			TelegramID: unix,
		}

		loginResp, err := sdkClient.Auth(ctx, loginReq)
		require.NoError(t, err)
		assert.Equal(t, loginReq.Username, loginResp.UserInfo.Username)
		assert.Equal(t, loginReq.FirstName, loginResp.UserInfo.FirstName)
		assert.Equal(t, loginReq.LastName, loginResp.UserInfo.LastName)
		assert.Equal(t, loginReq.TelegramID, loginResp.UserInfo.Extra.TelegramID)

		meResp, err := sdkClient.Me(ctx)
		require.NoError(t, err)
		assert.Equal(t, loginReq.Username, meResp.User.Username)
		assert.Equal(t, loginReq.FirstName, meResp.User.FirstName)
		assert.Equal(t, loginReq.LastName, meResp.User.LastName)
		assert.Equal(t, loginReq.TelegramID, meResp.User.Extra.TelegramID)
	})

	t.Run("meet and bill flow - success", func(t *testing.T) {
		sdkClient := sdk.New(baseURL)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		unix := time.Now().Unix()
		loginReq := auth.LoginTelegramReq{
			Username:   fmt.Sprintf("meet_user_%d", unix),
			FirstName:  "Meet",
			LastName:   fmt.Sprintf("User_%d", unix),
			TelegramID: unix,
		}

		loginResp, err := sdkClient.Auth(ctx, loginReq)
		require.NoError(t, err)

		meetReq := meet.CreateEventRq{
			EventName: fmt.Sprintf("Test meet %d", unix),
			Members:   []string{"first_member", "second_member"},
		}

		createdMeet, err := sdkClient.CreateMeet(ctx, meetReq)
		require.NoError(t, err)
		require.NotZero(t, createdMeet.ID)

		userMeets, err := sdkClient.FetchUserMeets(ctx)
		require.NoError(t, err)
		require.NotNil(t, userMeets)

		var fetchedMeet *meet.Event
		for i := range *userMeets {
			if (*userMeets)[i].ID == createdMeet.ID {
				fetchedMeet = &(*userMeets)[i]
				break
			}
		}

		require.NotNil(t, fetchedMeet)
		assert.Equal(t, meetReq.EventName, fetchedMeet.Name)
		assert.Len(t, fetchedMeet.Members, len(meetReq.Members)+1)

		meetDetails, err := sdkClient.GetMeetDetailsByID(ctx, createdMeet.ID)
		require.NoError(t, err)
		require.NotNil(t, meetDetails)

		var creatorMemberID int64
		var secondaryMemberID int64
		for _, member := range meetDetails.Members {
			if member.UserID != nil && *member.UserID == int64(loginResp.UserInfo.ID) {
				creatorMemberID = member.ID
				continue
			}

			if secondaryMemberID == 0 {
				secondaryMemberID = member.ID
			}
		}

		require.NotZero(t, creatorMemberID)
		require.NotZero(t, secondaryMemberID)

		billReq := bill.CreateBillRq{
			EventID:     createdMeet.ID,
			Name:        fmt.Sprintf("Test bill %d", unix),
			TotalAmount: 2000,
			Currency:    "RUB",
			SplitType:   "custom",
			PaidBy:      creatorMemberID,
			Participants: []bill.ParticipantRq{
				{
					MemberID: creatorMemberID,
					Amount:   1200,
				},
				{
					MemberID: secondaryMemberID,
					Amount:   800,
				},
			},
		}

		billResp, err := sdkClient.CreateBill(ctx, billReq)
		require.NoError(t, err)
		require.NotZero(t, billResp.ID)

		bills, err := sdkClient.FetchEventBills(ctx, createdMeet.ID)
		require.NoError(t, err)
		require.NotNil(t, bills)

		var createdBill *bill.Bill
		for i := range *bills {
			if (*bills)[i].ID == billResp.ID {
				createdBill = &(*bills)[i]
				break
			}
		}

		require.NotNil(t, createdBill)
		assert.Equal(t, billReq.Name, createdBill.Name)
		assert.Equal(t, billReq.EventID, createdBill.EventID)
		assert.Equal(t, billReq.PaidBy, createdBill.PaidBy)
		assert.Len(t, createdBill.Participants, len(billReq.Participants))

		summary, err := sdkClient.GetMeetSummary(ctx, createdMeet.ID)
		require.NoError(t, err)
		require.NotNil(t, summary)

		balances := make(map[int64]meet.Balance, len(summary.Balances))
		for _, balance := range summary.Balances {
			balances[balance.MemberID] = balance
		}

		require.Contains(t, balances, creatorMemberID)
		require.Contains(t, balances, secondaryMemberID)

		assert.Equal(t, billReq.TotalAmount, balances[creatorMemberID].TotalPaid)
		assert.Equal(t, billReq.Participants[0].Amount, balances[creatorMemberID].TotalShare)
		assert.Equal(t, billReq.Participants[1].Amount, balances[secondaryMemberID].TotalShare)

		var settlementFound bool
		for _, settlement := range summary.Settlements {
			if settlement.FromMemberID == secondaryMemberID &&
				settlement.ToMemberID == creatorMemberID &&
				settlement.Amount == billReq.Participants[1].Amount {
				settlementFound = true
				break
			}
		}

		assert.True(t, settlementFound)
	})

	t.Run("user profile and payment methods - success", func(t *testing.T) {
		sdkClient := sdk.New(baseURL)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		unix := time.Now().Unix()
		loginReq := auth.LoginTelegramReq{
			Username:   fmt.Sprintf("pm_user_%d", unix),
			FirstName:  "Payment",
			LastName:   fmt.Sprintf("User_%d", unix),
			TelegramID: unix,
		}

		loginResp, err := sdkClient.Auth(ctx, loginReq)
		require.NoError(t, err)
		require.NotNil(t, loginResp)

		profileReq := user.UpdateUserProfileReq{
			Username:  fmt.Sprintf("pm_user_updated_%d", unix),
			FirstName: "PaymentUpdated",
			LastName:  "UserUpdated",
		}

		updateProfileResp, err := sdkClient.UpdateUserProfile(ctx, profileReq)
		require.NoError(t, err)
		assert.Equal(t, "OK", updateProfileResp.Message)

		meResp, err := sdkClient.Me(ctx)
		require.NoError(t, err)
		require.NotNil(t, meResp)
		assert.Equal(t, profileReq.Username, meResp.User.Username)
		assert.Equal(t, profileReq.FirstName, meResp.User.FirstName)
		assert.Equal(t, profileReq.LastName, meResp.User.LastName)

		createPMReq := user.CreatePaymentMethodReq{
			Name:        "SBP",
			Description: "e2e test payment method",
			Recipient:   fmt.Sprintf("+7000%d", unix%1000000),
		}

		createdMethod, err := sdkClient.CreatePaymentMethod(ctx, createPMReq)
		require.NoError(t, err)
		require.NotNil(t, createdMethod)
		require.NotZero(t, createdMethod.ID)
		assert.Equal(t, createPMReq.Name, createdMethod.Name)
		assert.Equal(t, int64(loginResp.UserInfo.ID), createdMethod.UserID)

		methodsResp, err := sdkClient.GetPaymentMethods(ctx)
		require.NoError(t, err)
		require.NotNil(t, methodsResp)

		var storedMethod *user.PaymentMethodResponse
		for i := range methodsResp.PaymentMethods {
			if methodsResp.PaymentMethods[i].ID == createdMethod.ID {
				storedMethod = &methodsResp.PaymentMethods[i]
				break
			}
		}
		require.NotNil(t, storedMethod)

		updatePMReq := user.UpdatePaymentMethodReq{
			Name:        "Card",
			Description: "updated description",
			Recipient:   fmt.Sprintf("+7999%d", unix%1000000),
		}

		updatePMResp, err := sdkClient.UpdatePaymentMethod(ctx, createdMethod.ID, updatePMReq)
		require.NoError(t, err)
		assert.Equal(t, "OK", updatePMResp.Message)

		methodsAfterUpdate, err := sdkClient.GetPaymentMethods(ctx)
		require.NoError(t, err)

		var updatedMethod *user.PaymentMethodResponse
		for i := range methodsAfterUpdate.PaymentMethods {
			if methodsAfterUpdate.PaymentMethods[i].ID == createdMethod.ID {
				updatedMethod = &methodsAfterUpdate.PaymentMethods[i]
				break
			}
		}
		require.NotNil(t, updatedMethod)
		assert.Equal(t, updatePMReq.Name, updatedMethod.Name)
		assert.Equal(t, updatePMReq.Recipient, updatedMethod.Recipient)

		deleteResp, err := sdkClient.DeletePaymentMethod(ctx, createdMethod.ID)
		require.NoError(t, err)
		assert.Equal(t, "OK", deleteResp.Message)

		methodsAfterDelete, err := sdkClient.GetPaymentMethods(ctx)
		require.NoError(t, err)

		for _, method := range methodsAfterDelete.PaymentMethods {
			assert.NotEqual(t, createdMethod.ID, method.ID)
		}
	})
}
