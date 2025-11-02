package main

//Aquí creo los endpoints
//Crear una función que retorne en JSON los documentos de la colección containers que pertenezcan al usuario autenticado GET

import (
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/websocket"
)

type handler struct {
	store *store
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust this to ensure your origin checks are adequate for your security requirements
	},
}

func NewHandler(store *store) *handler {
	return &handler{store: store}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /alerts", h.HandleGetAlerts)
	mux.HandleFunc("POST /history/alerts", h.HandleGetHistoryAlerts)
	mux.HandleFunc("GET /ws/joinRoom", h.JoinRoom)

}

func (h *handler) HandleGetAlerts(w http.ResponseWriter, r *http.Request) {

	response, err := h.store.GetAlerts()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get alerts: "+err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, response)

}

func (h *handler) HandleGetHistoryAlerts(w http.ResponseWriter, r *http.Request) {

	var payload HistoryPayload

	if err := ParseJSON(r, &payload); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		formattedErrors := FormatValidationErrors(errors)
		WriteError(w, http.StatusBadRequest, "invalid payload: "+formattedErrors)
		return
	}

	response, err := h.store.GetAlertsByRange(payload.Start, payload.End)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get alerts: "+err.Error())
		return
	}

	Stats := h.store.CalcularEstadisticas(response)

	finalResponse := StatsPayload{
		Alerts:     response,
		AlertStats: Stats,
	}

	WriteJSON(w, http.StatusOK, finalResponse)

}

func (h *handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		WriteError(w, http.StatusUnauthorized, "upgrader fail: "+err.Error())
		return
	}
	defer conn.Close()

	log.Printf("A new user has joined")

	done := make(chan struct{})

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("Client disconnected:", err)
				close(done)
				return
			}
		}
	}()

	// 🔹 Obtener el último alerta antes de entrar al bucle
	lastAlert, err := h.store.GetLastAlert()
	if err != nil {
		log.Printf("Initial DB error: %v", err)
		lastAlert = []Alert{}
	}

	// 🔹 Si no hay alerta previa, inicializar vacío
	if len(lastAlert) > 0 {
		lastAlert[0].New = true // primer dato se marca como nuevo
		if err := conn.WriteJSON(lastAlert); err != nil {
			log.Printf("Failed to send initial message: %v", err)
			return
		}
	}

	for {
		select {
		case <-done:
			log.Println("Stopping sender loop")
			return
		default:
			time.Sleep(10 * time.Second)

			alert, err := h.store.GetLastAlert()
			if err != nil {
				log.Printf("DB error: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			if len(alert) == 0 {
				continue
			}

			if len(lastAlert) == 0 || alert[0].Ts != lastAlert[0].Ts {
				alert[0].New = true
				lastAlert = alert
			} else {
				alert[0].New = false
			}

			if err := conn.WriteJSON(alert); err != nil {
				log.Printf("Failed to send message: %v", err)
				return
			}
		}
	}
}
