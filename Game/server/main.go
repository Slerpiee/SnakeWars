package main

import (
	"context"
	"flag"
	"game_server/internal/game"
	"game_server/server/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	port := flag.Int("port", 8000, "Порт для запуска сервера")
	flag.Parse()

	GameServer := game.CreateServer()

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.MainHandler).Methods("GET")
	r.HandleFunc("/createRoom", handlers.CreateRoomHandler).Methods("GET")
	r.HandleFunc("/game", handlers.GameHandler).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))

	apiHandlers := handlers.NewAPIHandlers(GameServer)
	secretKey := handlers.GetJwt()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/createRoom", handlers.AuthMiddleware(secretKey, apiHandlers.CreateRoom)).Methods("POST")
	api.HandleFunc("/getRooms", handlers.AuthMiddleware(secretKey, apiHandlers.GetRooms)).Methods("GET")
	api.HandleFunc("/getRoom", handlers.AuthMiddleware(secretKey, apiHandlers.GetRoomByName)).Methods("GET")

	addr := "127.0.0.1:" + strconv.Itoa(*port)

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("Сервер запускается на %s", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("Получен сигнал: %v. Начинаем graceful shutdown...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	srv.Shutdown(ctx)
	defer cancel()

}
