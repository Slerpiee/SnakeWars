// handlers/api.go
package handlers

import (
	"encoding/json"
	"game_server/internal/game"
	"game_server/server/middleware"
	"net/http"

	_ "github.com/gorilla/mux"
)

type APIHandlers struct {
	Server *game.Server
}

func NewAPIHandlers(server *game.Server) *APIHandlers {
	return &APIHandlers{Server: server}
}

func (h *APIHandlers) CreateRoom(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	var request struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid Body format", http.StatusBadRequest)
		return
	}

	if userID == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if len(request.Name) <= 0 || len(request.Name) > 20 {
		http.Error(w, "Incorrect room name", http.StatusBadRequest)
		return
	}

	_, exists := h.Server.GetRoomIdByName(request.Name)
	if exists {
		http.Error(w, "Room already exists", http.StatusBadRequest)
		return
	}

	newRoom := game.CreateRoom(request.Name)

	h.Server.AddRoom(newRoom)

	json.NewEncoder(w).Encode(newRoom)

}

func (h *APIHandlers) GetRoomByName(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	var resp struct {
		Id    string `json:"id"`
		Exist bool   `json:"exist"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Body format", http.StatusBadRequest)
		return
	}
	userId := middleware.GetUserIDFromContext(r.Context())
	if userId == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	Id, exist := h.Server.GetRoomIdByName(req.Name)
	resp.Exist = exist
	resp.Id = Id
	json.NewEncoder(w).Encode(resp)
}

func (h *APIHandlers) GetRooms(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	rooms := h.Server.GetRooms()
	json.NewEncoder(w).Encode(rooms)
}
