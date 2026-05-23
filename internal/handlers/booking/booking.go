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
		errorResp := apperr.NewAppErr(err, err, http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResp.GetMapError())
		return
	}

	createCommand, err := req.ToCreateCommand()
	if err != nil {
		resp := apperr.NewAppErr(err, err, http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp.GetMapError())
		return
	}

	res, err := h.booking.CreateBooking(r.Context(), createCommand)

	if err != nil {
		switch {
		default:
			errResp := apperr.NewAppErr(err, err, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errResp.GetMapError())
			return
		}

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

		switch err {

		default:
			errResp := apperr.NewAppErr(err, err, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errResp.GetMapError())
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	err := h.booking.Cancel(r.Context(), command.CancelBookingCommand{ID: id})
	if err != nil {

		switch err {

		default:
			errResp := apperr.NewAppErr(err, err, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errResp.GetMapError())
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *BookingHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	res, err := h.booking.GetBooking(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, apperr.NotFoundErr):
			w.WriteHeader(http.StatusNotFound)
			errorResp := apperr.NewAppErr(errors.New("бронь не найдена"), err, http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(errorResp.GetMapError())
			return
		default:
			errResp := apperr.NewAppErr(err, err, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errResp.GetMapError())
			return
		}

	}
	var response dto.GetBookingResponse

	response.CreateQueryResultToResponse(res)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
