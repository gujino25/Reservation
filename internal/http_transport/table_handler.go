package httptransport

import (
	"context"
	"net/http"
	"strconv"

	"reservation/internal/domain"
)

type tableRepo interface {
	Create(ctx context.Context, t domain.Table) (domain.Table, error)
	ListByRestaurant(ctx context.Context, restaurantID int64) ([]domain.Table, error)
}

type TableHandler struct {
	repo tableRepo
}

func NewTableHandler(repo tableRepo) *TableHandler {
	return &TableHandler{repo: repo}
}

func (h *TableHandler) Create(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.ParseInt(pathParam(r, "id"), 10, 64)
	if err != nil || restaurantID <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	var req createTableRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, errInvalidBody)
		return
	}

	tbl, err := domain.NewTable(restaurantID, req.Number, req.Capacity)
	if err != nil {
		writeError(w, err)
		return
	}

	created, err := h.repo.Create(r.Context(), tbl)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toTableResponse(created))
}

func (h *TableHandler) ListByRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.ParseInt(pathParam(r, "id"), 10, 64)
	if err != nil || restaurantID <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	tables, err := h.repo.ListByRestaurant(r.Context(), restaurantID)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]tableResponse, 0, len(tables))
	for _, t := range tables {
		resp = append(resp, toTableResponse(t))
	}
	writeJSON(w, http.StatusOK, resp)
}
