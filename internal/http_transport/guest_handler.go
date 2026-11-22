package httptransport

import (
	"context"
	"net/http"
	"strconv"

	"reservation/internal/domain"
)

type guestRepo interface {
	Create(ctx context.Context, g domain.Guest) (domain.Guest, error)
	GetByID(ctx context.Context, id int64) (domain.Guest, error)
	GetByPhone(ctx context.Context, phone string) (domain.Guest, error)
}

type GuestHandler struct {
	repo guestRepo
}

func NewGuestHandler(repo guestRepo) *GuestHandler {
	return &GuestHandler{repo: repo}
}

func (h *GuestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createGuestRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, errInvalidBody)
		return
	}
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "name is required"})
		return
	}

	guest, err := h.repo.Create(r.Context(), domain.NewGuest(req.Name, req.Phone))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toGuestResponse(guest))
}

func (h *GuestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(pathParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	guest, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toGuestResponse(guest))
}

func (h *GuestHandler) GetByPhone(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "phone query parameter is required"})
		return
	}

	guest, err := h.repo.GetByPhone(r.Context(), phone)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toGuestResponse(guest))
}
