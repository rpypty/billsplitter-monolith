package user

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"billsplitter-monolith/internal/domain/payment_method"
	"billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/transport/http/middleware"
	hu "billsplitter-monolith/internal/utils/http"
	"github.com/go-chi/chi/v5"
)

type Controller interface {
	// UpdateUserProfile - обновляет данные профиля пользователя
	UpdateUserProfile(w http.ResponseWriter, r *http.Request)

	// GetPaymentMethods - получает список платежных методов пользователя
	GetPaymentMethods(w http.ResponseWriter, r *http.Request)

	// CreatePaymentMethod - создает новый платежный метод
	CreatePaymentMethod(w http.ResponseWriter, r *http.Request)

	// UpdatePaymentMethod - обновляет платежный метод
	UpdatePaymentMethod(w http.ResponseWriter, r *http.Request)

	// DeletePaymentMethod - удаляет платежный метод
	DeletePaymentMethod(w http.ResponseWriter, r *http.Request)
}

type controllerImpl struct {
	userSvc          user.Service
	paymentMethodSvc payment_method.Service
	logger           *slog.Logger
}

func NewController(
	userSvc user.Service,
	paymentMethodSvc payment_method.Service,
	logger *slog.Logger,
) Controller {
	return &controllerImpl{
		userSvc:          userSvc,
		paymentMethodSvc: paymentMethodSvc,
		logger:           logger,
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
	l := c.l().With("method", "UpdateUserProfile")

	sessionInfo, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	userID, err := hu.GetQueryParamInt(r, "id")
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid query param id")
		return
	}

	if sessionInfo.UserID != vo.UserID(userID) {
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

	err = c.userSvc.Update(ctx, vo.UserID(userID), userReq)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(fmt.Sprintf("userController.UpdateUserProfile error: %s", err))
		return
	}

	hu.RespondOK(w)
}

// GetPaymentMethods godoc
// @Summary      Получение списка платежных методов
// @Description  Возвращает список всех платежных методов пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id  path      int  false  "ID пользователя (опционально, если не указан - берется из сессии)"
// @Success      200     {object}  PaymentMethodsResponse  "Список платежных методов"
// @Failure      400     {object}  hu.ErrorResponse        "Некорректный запрос или ошибка сессии"
// @Failure      403     {object}  hu.ErrorResponse        "Доступ запрещён (несовпадение ID пользователя и сессии)"
// @Failure      500     {object}  hu.ErrorResponse        "Внутренняя ошибка сервера"
// @Router       /user/{id}/payment_methods [get]
// @Router       /payment_methods [get]
func (c *controllerImpl) GetPaymentMethods(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "GetPaymentMethods")

	sessionInfo, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	// Try to get user ID from URL path first
	userID := sessionInfo.UserID
	if urlUserID := chi.URLParam(r, "id"); urlUserID != "" {
		parsedUserID, err := strconv.ParseInt(urlUserID, 10, 64)
		if err != nil {
			hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid user id")
			return
		}
		userID = vo.UserID(parsedUserID)
	}

	// Verify user has access to this user's data
	if sessionInfo.UserID != userID {
		hu.RespondErrWithStatus(w, http.StatusForbidden, "Forbidden")
		return
	}

	methods, err := c.paymentMethodSvc.FetchByUserID(ctx, userID)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(fmt.Sprintf("userController.GetPaymentMethods error: %s", err))
		return
	}

	response := PaymentMethodsResponse{
		PaymentMethods: make([]PaymentMethodResponse, 0, len(methods)),
	}

	for _, method := range methods {
		response.PaymentMethods = append(response.PaymentMethods, PaymentMethodResponse{
			ID:          method.ID,
			UserID:      method.UserID,
			Name:        method.Name,
			Description: method.Description,
			Recipient:   method.Recipient,
		})
	}

	hu.RespondJson(w, response)
}

// CreatePaymentMethod godoc
// @Summary      Создание платежного метода
// @Description  Создает новый платежный метод для пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id  path      int                    false  "ID пользователя (опционально, если не указан - берется из сессии)"
// @Param        request body      CreatePaymentMethodReq true  "Данные для создания платежного метода"
// @Success      201     {object}  PaymentMethodResponse  "Платежный метод успешно создан"
// @Failure      400     {object}  hu.ErrorResponse       "Некорректный запрос или ошибка сессии"
// @Failure      403     {object}  hu.ErrorResponse       "Доступ запрещён (несовпадение ID пользователя и сессии)"
// @Failure      500     {object}  hu.ErrorResponse       "Внутренняя ошибка сервера"
// @Router       /user/{id}/payment_methods [post]
// @Router       /payment_methods [post]
func (c *controllerImpl) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "CreatePaymentMethod")

	// Debug logging for request details
	l.Info("CreatePaymentMethod request details",
		"contentType", r.Header.Get("Content-Type"),
		"method", r.Method,
		"url", r.URL.String())

	sessionInfo, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	// Try to get user ID from URL path first
	userID := sessionInfo.UserID
	if urlUserID := chi.URLParam(r, "id"); urlUserID != "" {
		parsedUserID, err := strconv.ParseInt(urlUserID, 10, 64)
		if err != nil {
			hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid user id")
			return
		}
		userID = vo.UserID(parsedUserID)
	}

	// Verify user has access to this user's data
	if sessionInfo.UserID != userID {
		hu.RespondErrWithStatus(w, http.StatusForbidden, "Forbidden")
		return
	}

	req, err := hu.DecodeReq[CreatePaymentMethodReq](r)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	// Debug logging to see what's received
	l.Info("CreatePaymentMethod request received",
		"name", req.Name,
		"description", req.Description,
		"recipient", req.Recipient)

	// Manual validation for required fields
	if req.Name == "" {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "name field is required")
		return
	}
	if req.Recipient == "" {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "recipient field is required")
		return
	}

	paymentMethod := payment_method.PaymentMethod{
		UserID:      int64(userID),
		Name:        req.Name,
		Description: req.Description,
		Recipient:   req.Recipient,
	}

	createdMethod, err := c.paymentMethodSvc.Create(ctx, paymentMethod)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(fmt.Sprintf("userController.CreatePaymentMethod error: %s", err))
		return
	}

	response := PaymentMethodResponse{
		ID:          createdMethod.ID,
		UserID:      createdMethod.UserID,
		Name:        createdMethod.Name,
		Description: createdMethod.Description,
		Recipient:   createdMethod.Recipient,
	}

	hu.RespondJsonWithStatus(w, http.StatusCreated, response)
}

// UpdatePaymentMethod godoc
// @Summary      Обновление платежного метода
// @Description  Обновляет существующий платежный метод пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id  path      int                    false  "ID пользователя (опционально, если не указан - берется из сессии)"
// @Param        methodId      path      int                    true  "ID платежного метода"
// @Param        request body      UpdatePaymentMethodReq true  "Данные для обновления платежного метода"
// @Success      200     {object}  hu.ResponseOK          "Платежный метод успешно обновлён"
// @Failure      400     {object}  hu.ErrorResponse       "Некорректный запрос или ошибка сессии"
// @Failure      403     {object}  hu.ErrorResponse       "Доступ запрещён (несовпадение ID пользователя и сессии)"
// @Failure      500     {object}  hu.ErrorResponse       "Внутренняя ошибка сервера"
// @Router       /user/{id}/payment_methods/{methodId} [put]
// @Router       /payment_methods/{methodId} [put]
func (c *controllerImpl) UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "UpdatePaymentMethod")

	sessionInfo, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	// Try to get user ID from URL path first
	userID := sessionInfo.UserID
	if urlUserID := chi.URLParam(r, "id"); urlUserID != "" {
		parsedUserID, err := strconv.ParseInt(urlUserID, 10, 64)
		if err != nil {
			hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid user id")
			return
		}
		userID = vo.UserID(parsedUserID)
	}

	methodID, err := strconv.ParseInt(chi.URLParam(r, "methodId"), 10, 64)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid method id")
		return
	}

	// Verify user has access to this user's data
	if sessionInfo.UserID != userID {
		hu.RespondErrWithStatus(w, http.StatusForbidden, "Forbidden")
		return
	}

	req, err := hu.DecodeReq[UpdatePaymentMethodReq](r)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	paymentMethod := payment_method.PaymentMethod{
		UserID:      int64(userID),
		Name:        req.Name,
		Description: req.Description,
		Recipient:   req.Recipient,
	}

	err = c.paymentMethodSvc.Update(ctx, methodID, paymentMethod)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(fmt.Sprintf("userController.UpdatePaymentMethod error: %s", err))
		return
	}

	hu.RespondOK(w)
}

// DeletePaymentMethod godoc
// @Summary      Удаление платежного метода
// @Description  Удаляет платежный метод пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id  path      int  false  "ID пользователя (опционально, если не указан - берется из сессии)"
// @Param        methodId      path      int  true  "ID платежного метода"
// @Success      200     {object}  hu.ResponseOK    "Платежный метод успешно удалён"
// @Failure      400     {object}  hu.ErrorResponse "Некорректный запрос или ошибка сессии"
// @Failure      403     {object}  hu.ErrorResponse "Доступ запрещён (несовпадение ID пользователя и сессии)"
// @Failure      500     {object}  hu.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /user/{id}/payment_methods/{methodId} [delete]
// @Router       /payment_methods/{methodId} [delete]
func (c *controllerImpl) DeletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "DeletePaymentMethod")

	sessionInfo, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	// Try to get user ID from URL path first
	userID := sessionInfo.UserID
	if urlUserID := chi.URLParam(r, "id"); urlUserID != "" {
		parsedUserID, err := strconv.ParseInt(urlUserID, 10, 64)
		if err != nil {
			hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid user id")
			return
		}
		userID = vo.UserID(parsedUserID)
	}

	methodID, err := strconv.ParseInt(chi.URLParam(r, "methodId"), 10, 64)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid method id")
		return
	}

	// Verify user has access to this user's data
	if sessionInfo.UserID != userID {
		hu.RespondErrWithStatus(w, http.StatusForbidden, "Forbidden")
		return
	}

	err = c.paymentMethodSvc.Delete(ctx, methodID)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(fmt.Sprintf("userController.DeletePaymentMethod error: %s", err))
		return
	}

	hu.RespondOK(w)
}

func (c *controllerImpl) l() *slog.Logger {
	return c.logger.WithGroup("user-controller")
}
