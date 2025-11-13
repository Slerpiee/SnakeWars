package ws_handler

import (
	"game_server/internal/game"
	"net/http"
	_ "os"
	_ "path/filepath"

	"github.com/gorilla/websocket"
)

type WShandler struct{
	Server *game.Server
	upgrader *websocket.Upgrader
}

func (wsh *WShandler) websocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := wsh.upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    
    userID := r.URL.Query().Get("auth") //jwt token here
	
    
    room := roomManager.GetOrCreateRoom(roomID)
    user := NewUser(userID, conn)
    
    room.AddUser(user)
}