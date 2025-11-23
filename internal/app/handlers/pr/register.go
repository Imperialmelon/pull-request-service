package pr

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Register(r *mux.Router, prUC PRUseCase) {
	h := NewHandler(prUC)

	r.HandleFunc("/pullRequest/merge", h.Merge).Methods(http.MethodPost)
	r.HandleFunc("/pullRequest/create", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/pullRequest/reassign", h.Reassign).Methods(http.MethodPost)
}
