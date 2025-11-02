package handlers

import (
	"net/http"
	"os"
	_ "path/filepath"

	_ "github.com/gorilla/mux"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "./templates/main.html"
	http.ServeFile(w, r, filePath)
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	filepath := "./templates/room.html"
	http.ServeFile(w, r, filepath)
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
	filepath := "./templates/game.html"
	http.ServeFile(w, r, filepath)
}

func GetJwt() string {
	if secret := os.Getenv("JWT_KEY"); secret != "" {
		return secret
	}
	return "Shmul_pedik"
}
