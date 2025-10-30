// handlers/api.go
package handlers

import (
    "encoding/json"
    "game_server/internal/game"
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
    var request struct {
        Name string `json:"name"`
    }

	
    
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
	if len(request.Name) <= 0 || len(request.Name) > 20{
		http.Error(w, "Incorrect room name", http.StatusBadRequest)
	}
    
	newRoom := game.CreateRoom(request.Name)

	h.Server.AddRoom(newRoom)

	json.NewEncoder(w)


    //	json.NewEncoder(w).Encode(room)
}



func (h *APIHandlers) GetRoom(w http.ResponseWriter, r *http.Request) {
	
}