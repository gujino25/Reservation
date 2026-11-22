package httptransport

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"reservation/internal/domain"
)

type restaurantRepo interface {
	Create(ctx context.Context, r domain.Restaurant) (domain.Restaurant, error)
	GetByID(ctx context.Context, id int64) (domain.Restaurant, error)
}

type RestaurantHandler struct {
	repo restaurantRepo
}

func NewRestaurantHandler(repo restaurantRepo) *RestaurantHandler {
	return &RestaurantHandler{repo: repo}
}

func (h *RestaurantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRestaurantRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, errInvalidBody)
		return
	}
	if req.Name == "" || req.Address == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "name and address are required"})
		return
	}

	rest, err := h.repo.Create(r.Context(), domain.Restaurant{
		Name:      req.Name,
		Address:   req.Address,
		CreatedAt: time.Now(),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toRestaurantResponse(rest))
}

func (h *RestaurantHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(pathParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	rest, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toRestaurantResponse(rest))
}
