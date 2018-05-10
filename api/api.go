package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// InitAPI initialiaze the transmitter http api
func InitAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/event", customEventHandler).Methods("POST")
	http.Handle("/api/v1", r)
}

func customEventHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

}
