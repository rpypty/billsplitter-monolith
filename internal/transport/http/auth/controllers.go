package auth

import (
	"log/slog"
	"net/http"

	"billsplitter-monolith/internal/domain/session"
	"billsplitter-monolith/internal/domain/user"
	"billsplitter-monolith/internal/transport/http/middleware"
	useruc "billsplitter-monolith/internal/usecase/user"
	hu "billsplitter-monolith/internal/utils/http"
)

type Controller interface {
	// LoginTelegram - создает сессию пользователя использую Telegram.initData
	LoginTelegram(w http.ResponseWriter, r *http.Request)

	// Me - возвращает данные о пользователе по сессии
	Me(w http.ResponseWriter, r *http.Request)
}

type controllerImpl struct {
	userUC     useruc.UseCase
	userSvc    user.Service
	sessionSvc session.Service
	logger     *slog.Logger
}

func NewController(
	userUC useruc.UseCase,
	userSvc user.Service,
	sessionSvc session.Service,
	logger *slog.Logger,
) Controller {
	return &controllerImpl{
		userUC:     userUC,
		userSvc:    userSvc,
		sessionSvc: sessionSvc,
		logger:     logger,
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

	userReq := user.User{
		Username:  rq.Username,
		FirstName: rq.FirstName,
		LastName:  rq.LastName,
		Extra: user.ExtraInfo{
			TelegramID: rq.TelegramID,
		},
	}

	userInfo, err := c.userUC.GetByTgIDOrCreate(ctx, userReq)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(err.Error())
		return
	}

	sessionID, err := c.sessionSvc.Create(ctx, userInfo.ID)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(err.Error())
		return
	}

	hu.RespondJson(w, &LoginTelegramRes{
		SessionID: sessionID,
		UserInfo:  userInfo,
	})
}

// Me godoc
// @Summary      Получить данные текущего пользователя
// @Description  Возвращает данные пользователя, извлечённые по sessionID из контекста
// @Tags         auth
// @Produce      json
// @Security     SessionAuth
// @Success      200  {object}  MeRes
// @Failure      400  {object}  hu.ErrorResponse  "Пользователь не найден или сессия невалидна"
// @Router       /auth/me [get]
func (c *controllerImpl) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "Me")

	sessionInfo, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	userInfo, err := c.userSvc.GetByID(ctx, sessionInfo.UserID)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	hu.RespondJson(w, &MeRes{
		User: userInfo,
	})
}

func (c *controllerImpl) l() *slog.Logger {
	return c.logger.WithGroup("auth-controller")
}
