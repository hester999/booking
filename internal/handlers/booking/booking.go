package booking

import (
	"booking/apperr"
	command "booking/internal/application/command/booking"
	"booking/internal/handlers/booking/dto"
	"booking/pkg"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type BookingHandler struct {
	booking booking
	log     pkg.Logger
}

func NewBookingHandler(booking booking, log pkg.Logger) *BookingHandler {
	return &BookingHandler{booking: booking, log: log}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	createCommand, err := req.ToCreateCommand()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	res, err := h.booking.CreateBooking(r.Context(), createCommand)
	if err != nil {
		writeBookingError(w, err)
		return
	}

	var response dto.CreateBookingResponse
	response.CreateQueryResultToResponse(res)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *BookingHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	err := h.booking.Confirm(r.Context(), command.ConfirmBookingCommand{ID: id})
	if err != nil {
		writeBookingError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	err := h.booking.Cancel(r.Context(), command.CancelBookingCommand{ID: id})
	if err != nil {
		writeBookingError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BookingHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	res, err := h.booking.GetBooking(r.Context(), id)
	if err != nil {
		writeBookingError(w, err)
		return
	}

	var response dto.GetBookingResponse
	response.CreateQueryResultToResponse(res)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func writeBookingError(w http.ResponseWriter, err error) {
	writeError(w, err, bookingErrorStatus(err))
}

func writeError(w http.ResponseWriter, err error, status int) {
	errResp := apperr.NewAppErr(err, err, status)

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errResp.GetMapError())
}

func bookingErrorStatus(err error) int {
	switch {
	case errors.Is(err, apperr.NotFoundErr):
		return http.StatusNotFound
	case errors.Is(err, apperr.InvalidStatusErr),
		errors.Is(err, apperr.InvalidDateErr),
		errors.Is(err, apperr.InvalidUserIDErr),
		errors.Is(err, apperr.InvalidListingID),
		errors.Is(err, apperr.InvalidPeriodErr),
		errors.Is(err, apperr.IntersectionErr),
		errors.Is(err, apperr.InvalidPriceErr):
		return http.StatusBadRequest
	case errors.Is(err, apperr.BookingAlreadyExistsErr),
		errors.Is(err, apperr.StatusCancel),
		errors.Is(err, apperr.StatusConfirmed),
		errors.Is(err, apperr.StatusExpire),
		errors.Is(err, apperr.UnsupportedStatus),
		errors.Is(err, apperr.NotExpiredErr):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
