package meet

import (
	"fmt"
	"log/slog"
	"net/http"

	"billsplitter-monolith/internal/transport/http/middleware"
	"billsplitter-monolith/internal/usecase/event"
	hu "billsplitter-monolith/internal/utils/http"
)

type Controller interface {
	// CreateMeet - создание мита
	CreateMeet(w http.ResponseWriter, r *http.Request)
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
// @Param        request  body      CreateEventRq   true  "Данные мита"
// @Success      201      {object}  hu.ResponseID        "ID созданного мита"
// @Failure      400      {object}  hu.ErrorResponse     "Ошибка валидации"
// @Failure      401      {object}  hu.ErrorResponse     "Неавторизован"
// @Failure      500      {object}  hu.ErrorResponse     "Внутренняя ошибка сервера"
// @Router       /event [post]
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
		Members:         nil,
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

func (c *controllerImpl) l() *slog.Logger {
	return c.logger.WithGroup("auth-controller")
}
