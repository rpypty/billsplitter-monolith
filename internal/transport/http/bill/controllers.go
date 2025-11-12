package bill

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	domainbill "billsplitter-monolith/internal/domain/bill"
	vo "billsplitter-monolith/internal/domain/valueobject"
	apperrors "billsplitter-monolith/internal/errors"
	"billsplitter-monolith/internal/transport/http/middleware"
	billuc "billsplitter-monolith/internal/usecase/bill"
	hu "billsplitter-monolith/internal/utils/http"
)

type Controller interface {
	CreateBill(w http.ResponseWriter, r *http.Request)
	FetchEventBills(w http.ResponseWriter, r *http.Request)
}

type controllerImpl struct {
	billUC billuc.UseCase
	logger *slog.Logger
}

func NewController(
	billUC billuc.UseCase,
	logger *slog.Logger,
) Controller {
	return &controllerImpl{
		billUC: billUC,
		logger: logger,
	}
}

// CreateBill godoc
// @Summary      Создает чек внутри ивента
// @Description  Чек создается от имени пользователя из сессии
// @Tags         bill
// @Accept       json
// @Produce      json
// @Param        request  body      CreateBillRq   true  "Данные чека"
// @Success      201      {object}  hu.ResponseID
// @Failure      400      {object}  hu.ErrorResponse
// @Failure      401      {object}  hu.ErrorResponse
// @Failure      404      {object}  hu.ErrorResponse
// @Failure      500      {object}  hu.ErrorResponse
// @Security     SessionAuth
// @Router       /bills [post]
func (c *controllerImpl) CreateBill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "CreateBill")

	sessionInfo, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	rq, err := hu.DecodeReq[CreateBillRq](r)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	participants := make([]domainbill.Participant, 0, len(rq.Participants))

	for _, p := range rq.Participants {
		participants = append(participants, domainbill.Participant{
			MemberID: p.MemberID,
			Amount:   p.Amount,
		})
	}

	createRq := billuc.CreateBillRq{
		EventID:      rq.EventID,
		Name:         rq.Name,
		CreatedBy:    sessionInfo.UserID,
		TotalAmount:  vo.Amount(rq.TotalAmount),
		Currency:     vo.CurrencyCode(rq.Currency),
		SplitType:    vo.SplitType(rq.SplitType),
		Participants: participants,
	}

	billID, err := c.billUC.CreateBill(ctx, createRq)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, apperrors.ErrEventNotFound) {
			status = http.StatusNotFound
		}
		hu.RespondErrWithStatus(w, status, err.Error())
		l.Error(fmt.Sprintf("billController.CreateBill error: %s", err))
		return
	}

	hu.RespondJsonWithStatus(w, http.StatusCreated, &hu.ResponseID{
		ID: billID,
	})
}

// FetchEventBills godoc
// @Summary      Получить список чеков ивента
// @Tags         bill
// @Accept       json
// @Produce      json
// @Param        event_id  query     int     true  "ID ивента"
// @Success      200       {array}   Bill
// @Failure      400       {object}  hu.ErrorResponse
// @Failure      401       {object}  hu.ErrorResponse
// @Failure      404       {object}  hu.ErrorResponse
// @Failure      500       {object}  hu.ErrorResponse
// @Security     SessionAuth
// @Router       /bills [get]
func (c *controllerImpl) FetchEventBills(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := c.l().With("method", "FetchEventBills")

	_, err := middleware.SessionFromContext(ctx)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, err.Error())
		l.Error(err.Error())
		return
	}

	eventIDParam := r.URL.Query().Get("event_id")
	if eventIDParam == "" {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "event_id is required")
		return
	}

	eventID, err := strconv.ParseInt(eventIDParam, 10, 64)
	if err != nil {
		hu.RespondErrWithStatus(w, http.StatusBadRequest, "invalid query param event_id")
		return
	}

	bills, err := c.billUC.FetchEventBills(ctx, eventID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, apperrors.ErrEventNotFound) {
			status = http.StatusNotFound
		}
		hu.RespondErrWithStatus(w, status, err.Error())
		l.Error(fmt.Sprintf("billController.FetchEventBills error: %s", err))
		return
	}

	resp := make([]Bill, 0, len(bills))

	for _, b := range bills {
		resp = append(resp, fromDomainBill(b))
	}

	hu.RespondJson(w, resp)
}

func (c *controllerImpl) l() *slog.Logger {
	return c.logger.WithGroup("bill-controller")
}
