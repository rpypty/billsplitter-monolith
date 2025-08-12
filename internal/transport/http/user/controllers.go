package user

import (
	"fmt"
	"log/slog"
	"net/http"

	"billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/transport/http/middleware"
	hu "billsplitter-monolith/internal/utils/http"
)

type Controller interface {
	// UpdateUserProfile - обновляет данные профиля пользователя
	UpdateUserProfile(w http.ResponseWriter, r *http.Request)

	// TODO: payment methods...
}

type controllerImpl struct {
	userSvc user.Service
	logger  *slog.Logger
}

func NewController(
	userSvc user.Service,
	logger *slog.Logger,
) Controller {
	return &controllerImpl{
		userSvc: userSvc,
		logger:  logger,
	}
}

// UpdateUserProfile godoc
// @Summary      Обновление профиля пользователя
// @Description  Обновляет профиль пользователя по переданному ID, если пользователь аутентифицирован и имеет доступ
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id       path      int                     true  "ID пользователя"
// @Param        request  body      UpdateUserProfileReq    true  "Данные для обновления профиля"
// @Success      200      {object}  hu.ResponseOK           "Профиль успешно обновлён"
// @Failure      400      {object}  hu.ErrorResponse        "Некорректный запрос или ошибка сессии"
// @Failure      403      {object}  hu.ErrorResponse        "Доступ запрещён (несовпадение ID пользователя и сессии)"
// @Failure      500      {object}  hu.ErrorResponse        "Внутренняя ошибка сервера (в debug режиме — детали)"
// @Router       /user/{id}/profile [put]
func (c *controllerImpl) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "LoginTelegram")

	sessionInfo, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	userId, err := hu.GetQueryParamInt(r, "id")
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid query param id")
		return
	}

	if sessionInfo.UserID != vo.UserID(userId) {
		hu.RespondErrWithStatus(w, http.StatusForbidden, "Forbidden")
		return
	}

	rq, err := hu.DecodeReq[UpdateUserProfileReq](r)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	userReq := user.User{
		Username:  rq.Username,
		FirstName: rq.FirstName,
		LastName:  rq.LastName,
	}

	err = c.userSvc.Update(ctx, vo.UserID(userId), userReq)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(fmt.Sprintf("userController.UpdateUserProfile error: %s", err))
		return
	}

	hu.RespondOK(w)
}

func (c *controllerImpl) l() *slog.Logger {
	return c.logger.WithGroup("auth-controller")
}
