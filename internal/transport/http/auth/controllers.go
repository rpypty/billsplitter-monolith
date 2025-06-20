package auth

import (
	"log/slog"
	"net/http"

	"billsplitter-monolith/internal/domain/auth"
	"billsplitter-monolith/internal/transport/http/middleware"
	hu "billsplitter-monolith/internal/utils/http"
)

type Controller interface {
	// LoginTelegram - создает сессию пользователя использую Telegram.initData
	LoginTelegram(w http.ResponseWriter, r *http.Request)

	// Me - возвращает данные о пользователе по сессии
	Me(w http.ResponseWriter, r *http.Request)
}

type controllerImpl struct {
	svc    auth.Service
	logger *slog.Logger
}

func NewController(svc auth.Service, logger *slog.Logger) Controller {
	return &controllerImpl{
		svc:    svc,
		logger: logger,
	}
}

// LoginTelegram godoc
// @Summary      Авторизация через Telegram
// @Description  Создаёт или получает пользователя по Telegram ID и возвращает sessionID
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginTelegramReq  true  "Данные пользователя из Telegram"
// @Success      200      {object}  LoginTelegramRes
// @Failure      400      {object}  hu.ErrorResponse  "Некорректный запрос"
// @Failure      500      {object}  hu.ErrorResponse  "Internal Server Error, но в debug моде возвращает детали ошибки"
// @Router       /auth/login/telegram [post]
func (c *controllerImpl) LoginTelegram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "LoginTelegram")

	rq, err := hu.DecodeReq[LoginTelegramReq](r)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.svc.CreateOrGetUserByTgID(ctx, rq.TelegramID, &auth.User{
		Username:  rq.Username,
		FirstName: rq.FirstName,
		LastName:  rq.LastName,
		Extra: auth.UserExtra{
			TelegramID: rq.TelegramID,
		},
	})
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(err.Error())
		return
	}

	sessionID, err := c.svc.CreateSession(ctx, &auth.Session{
		UserID: user.ID,
	})
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(err.Error())
		return
	}

	hu.RespondJson(w, &LoginTelegramRes{
		SessionID: sessionID,
	})
}

// Me godoc
// @Summary      Получить данные текущего пользователя
// @Description  Возвращает данные пользователя, извлечённые по sessionID из контекста
// @Tags         auth
// @Produce      json
// @Success      200  {object}  MeRes
// @Failure      400  {object}  hu.ErrorResponse  "Пользователь не найден или сессия невалидна"
// @Router       /auth/me [get]
func (c *controllerImpl) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "Me")

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	hu.RespondJson(w, &MeRes{
		User: user,
	})
}

func (c *controllerImpl) l() *slog.Logger {
	return c.logger.WithGroup("auth-controller")
}
