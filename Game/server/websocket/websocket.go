package ws_handler

import (
	"game_server/internal/game"
    "game_server/server/handlers"
    "game_server/server/middleware"
	"net/http"
	_ "os"
	_ "path/filepath"

	
	"github.com/gorilla/websocket"
)

type WSmanager struct{
	Server *game.Server
	upgrader *websocket.Upgrader
}

func CreateWShandler(s *game.Server) *WSmanager{
    return &WSmanager{
        Server:s,
        upgrader: nil,
    }
}

func (wsh *WSmanager) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
    _, err := wsh.upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    token := r.URL.Query().Get("auth") //jwt token here

    
}