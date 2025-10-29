package handlers

import (
	"net/http"
	_ "path/filepath"

	_ "github.com/gorilla/mux"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
    filePath := "./templates/main.html"
    http.ServeFile(w, r, filePath)
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request){
	filepath := "./templates/room.html"
	http.ServeFile(w, r, filepath)
}