package user

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Register(r *mux.Router, userUC UserUseCase) {
	h := NewHandler(userUC)

	r.HandleFunc("/users/setIsActive", h.SetIsActive).Methods(http.MethodPost)
	r.HandleFunc("/users/getReview/{id}", h.GetReview).Methods(http.MethodGet)
}
