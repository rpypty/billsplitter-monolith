package meet

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	apperrors "billsplitter-monolith/internal/errors"
	"billsplitter-monolith/internal/transport/http/middleware"
	"billsplitter-monolith/internal/usecase/event"
	hu "billsplitter-monolith/internal/utils/http"
)

type Controller interface {
	// CreateMeet - создание мита
	CreateMeet(w http.ResponseWriter, r *http.Request)

	// FetchUserMeets - получение персонального списка митов для юзера
	FetchUserMeets(w http.ResponseWriter, r *http.Request)

	// GetMeetDetailsByID - получение инфы мита по айди
	GetMeetDetailsByID(w http.ResponseWriter, r *http.Request)

	// GetSummary - получение сбивки по ивенту
	GetSummary(w http.ResponseWriter, r *http.Request)

	// AssignMemberToUser - связывает участника с текущим пользователем
	AssignMemberToUser(w http.ResponseWriter, r *http.Request)
}

type controllerImpl struct {
	eventUC event.UseCase
	logger  *slog.Logger
}

func NewController(
	eventUC event.UseCase,
	logger *slog.Logger,
) Controller {
	return &controllerImpl{
		eventUC: eventUC,
		logger:  logger,
	}
}

// CreateMeet godoc
// @Summary      Создает новый мит
// @Description  Создает новый мит с участниками
// @Tags         meet
// @Accept       json
// @Produce      json
// @Param        request  body      CreateEventRq  true  "Данные мита"
// @Success      201      {object}  hu.ResponseID
// @Failure      401      {object}  hu.ErrorResponse
// @Failure      500      {object}  hu.ErrorResponse
// @Security     SessionAuth
// @Router       /meets [post]
func (c *controllerImpl) CreateMeet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "CreateMeet")

	sessionInfo, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	rq, err := hu.DecodeReq[CreateEventRq](r)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	createMeetRq := event.CreateMeetRq{
		EventName:       rq.EventName,
		Date:            rq.Date,
		CreatedByUserID: sessionInfo.UserID,
		Members:         rq.Members,
	}

	meetID, err := c.eventUC.CreateMeet(ctx, createMeetRq)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(fmt.Sprintf("meetController.CreateMeetRq error: %s", err))
		return
	}

	hu.RespondJson(w, &hu.ResponseID{
		ID: meetID,
	})
}

// FetchUserMeets godoc
// @Summary      Получение персонального списка ивентов
// @Description  ID юзера берется из сессии
// @Tags         meet
// @Accept       json
// @Produce      json
// @Success      200  {array}   Event
// @Failure      401  {object}  hu.ErrorResponse
// @Failure      500  {object}  hu.ErrorResponse
// @Security     SessionAuth
// @Router       /meets [get]
func (c *controllerImpl) FetchUserMeets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "FetchUserMeets")

	sessionInfo, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	if sessionInfo == nil {
		hu.RespondErrWithStatus(w, http.StatusUnauthorized, "Unauthorized")
		l.Error("User not authorized")
		return
	}

	meets, err := c.eventUC.FetchUserMeets(ctx, sessionInfo.UserID)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusInternalServerError, err.Error())
		l.Error(fmt.Sprintf("meetController.FetchUserMeets error: %s", err))
		return
	}

	out := make([]Event, 0, len(meets))

	for _, meet := range meets {
		out = append(out, fromDomainEvent(meet))
	}

	hu.RespondJson(w, out)
}

// GetMeetDetailsByID godoc
// @Summary      Получение инфы по ивенту
// @Tags         meet
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID мита"
// @Success      200  {object}  EventSummary
// @Failure      401  {object}  hu.ErrorResponse
// @Failure      404  {object}  hu.ErrorResponse
// @Failure      500  {object}  hu.ErrorResponse
// @Security     SessionAuth
// @Router       /meets/{id} [get]
func (c *controllerImpl) GetMeetDetailsByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "GetMeetDetailsByID")

	meetID, err := hu.GetQueryParamInt(r, "id")
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid query param id")
		return
	}

	meet, err := c.eventUC.GetMeetByID(ctx, meetID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, apperrors.ErrForbiden) {
			status = http.StatusForbidden
		}
		hu.RespondErrWithStatus(w, status, err.Error())
		l.Error(fmt.Sprintf("meetController.GetMeetDetailsByID error: %s", err))
		return
	}

	if meet == nil {
		hu.RespondErrWithStatus(w, http.StatusNotFound, "meet not found")
		l.Error(fmt.Sprintf("meetController.GetMeetDetailsByID error: %s", err))
		return
	}

	hu.RespondJson(w, fromDomainEvent(*meet))
}

// GetSummary godoc
// @Summary      Получение сбивки по ивенту
// @Tags         meet
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID мита"
// @Success      200  {object}  EventSummary
// @Failure      401  {object}  hu.ErrorResponse
// @Failure      404  {object}  hu.ErrorResponse
// @Failure      500  {object}  hu.ErrorResponse
// @Security     SessionAuth
// @Router       /meets/{id}/summary [get]
func (c *controllerImpl) GetSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "GetSummary")

	meetID, err := hu.GetQueryParamInt(r, "id")
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid query param id")
		return
	}

	summary, err := c.eventUC.CalculateSummary(ctx, meetID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, apperrors.ErrForbiden) {
			status = http.StatusForbidden
		} else if errors.Is(err, apperrors.ErrEventNotFound) {
			status = http.StatusNotFound
		}
		hu.RespondErrWithStatus(w, status, err.Error())
		l.Error(fmt.Sprintf("meetController.GetSummary error: %s", err))
		return
	}

	hu.RespondJson(w, fromDomainSummary(*summary))
}

// AssignMemberToUser godoc
// @Summary      Связать участника с текущим пользователем
// @Description  user_id берется из сессии
// @Tags         meet
// @Accept       json
// @Produce      json
// @Param        id       path      int             true  "ID мита"
// @Param        request  body      SelectMemberRq  true  "member_id"
// @Success      200      {object}  hu.ResponseOK
// @Failure      400      {object}  hu.ErrorResponse
// @Failure      401      {object}  hu.ErrorResponse
// @Failure      404      {object}  hu.ErrorResponse
// @Failure      500      {object}  hu.ErrorResponse
// @Security     SessionAuth
// @Router       /meets/{id}/assign [put]
func (c *controllerImpl) AssignMemberToUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "AssignMemberToUser")

	sessionInfo, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	if sessionInfo == nil {
		hu.RespondErrWithStatus(w, http.StatusUnauthorized, "Unauthorized")
		l.Error("User not authorized")
		return
	}

	meetID, err := hu.GetQueryParamInt(r, "id")
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid query param id")
		return
	}

	rq, err := hu.DecodeReq[SelectMemberRq](r)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	if rq.MemberID == 0 {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "member_id is required")
		return
	}

	err = c.eventUC.AssignMemberToUser(ctx, meetID, rq.MemberID, sessionInfo.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, apperrors.ErrEventNotFound):
			status = http.StatusNotFound
		case strings.HasPrefix(err.Error(), "validation error"):
			status = http.StatusBadRequest
		case errors.Is(err, apperrors.ErrForbiden):
			status = http.StatusForbidden
		}
		hu.RespondErrWithStatus(w, status, err.Error())
		l.Error(fmt.Sprintf("meetController.AssignMemberToUser error: %s", err))
		return
	}

	hu.RespondOK(w)
}

func (c *controllerImpl) l() *slog.Logger {
	return c.logger.WithGroup("auth-controller")
}
