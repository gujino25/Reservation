package httptransport

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"reservation/internal/domain"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("unexpected data after JSON body")
		}
		return err
	}
	return nil
}

var errInvalidBody = errors.New("invalid request body")

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, domain.ErrAlreadyExists),
		errors.Is(err, domain.ErrTableBusy),
		errors.Is(err, domain.ErrCapacityExceeded),
		errors.Is(err, domain.ErrReservationAlreadyCancelled),
		errors.Is(err, domain.ErrReservationAlreadyStarted):
		status = http.StatusConflict
	case errors.Is(err, domain.ErrReservationTime),
		errors.Is(err, domain.ErrGuestAmountInvalid),
		errors.Is(err, domain.ErrCapacityInvalid),
		errors.Is(err, domain.ErrNumberOfTableInvalid),
		errors.Is(err, errInvalidBody):
		status = http.StatusBadRequest
	}
	writeJSON(w, status, errorResponse{Error: err.Error()})
}
