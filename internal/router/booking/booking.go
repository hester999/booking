package booking

import (
	"booking/internal/handlers/booking"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PATCH  = "PATCH"
	PUT    = "PUT"
)

func NewRouter(bookingHandler *booking.BookingHandler) *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	}).Methods(GET)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/booking", bookingHandler.Create).Methods(POST)
	api.HandleFunc("/booking/{id}/confirm", bookingHandler.Confirm).Methods(POST) //todo мокнуть сервис оплаты
	api.HandleFunc("/booking/{id}/cancel", bookingHandler.Cancel).Methods(POST)
	api.HandleFunc("/booking/{id}", bookingHandler.Get).Methods(GET)
	return router
}
