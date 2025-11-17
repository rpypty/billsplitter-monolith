package e2e

import (
	"billsplitter-monolith/internal/e2e/sdk"
	"billsplitter-monolith/internal/transport/http/auth"
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
}
