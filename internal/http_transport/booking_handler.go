package httptransport

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"reservation/internal/domain"
)

type bookingService interface {
	Create(ctx context.Context, tableID, guestID int64, startsAt time.Time, guestsCount int) (domain.Reservation, error)
	GetByID(ctx context.Context, id int64) (domain.Reservation, error)
	Cancel(ctx context.Context, id int64) (domain.Reservation, error)
}

type BookingHandler struct {
	svc bookingService
}

func NewBookingHandler(svc bookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createBookingRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, errInvalidBody)
		return
	}
	if req.TableID <= 0 || req.GuestID <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "table_id and guest_id are required"})
		return
	}

	booking, err := h.svc.Create(r.Context(), req.TableID, req.GuestID, req.StartsAt, req.GuestsCount)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toBookingResponse(booking))
}

func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(pathParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	booking, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toBookingResponse(booking))
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(pathParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	booking, err := h.svc.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toBookingResponse(booking))
}
