package main

//Aquí creo los endpoints
//Crear una función que retorne en JSON los documentos de la colección containers que pertenezcan al usuario autenticado GET

import (
	"net/http"
)

type handler struct {
	store *store
}

func NewHandler(store *store) *handler {
	return &handler{store: store}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /alerts", h.HandleUserRegister)


}

func (h *handler) HandleUserRegister(w http.ResponseWriter, r *http.Request) {

	response, err := h.store.GetAlerts()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get alerts: "+err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, response)

}
