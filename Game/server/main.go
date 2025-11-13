package main

import (
	"context"
	"flag"
	"game_server/internal/game"
	"game_server/server/handlers"
	"game_server/server/middleware"
	ws_handler "game_server/server/websocket"

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

	port := flag.Int("port", 8070, "Порт для запуска сервера")
	flag.Parse()



	GameServer := game.CreateServer()
	WScontroller := ws_handler.CreateWShandler(GameServer)

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	r.HandleFunc("/auth", handlers.AuthHandler).Methods("GET")
	r.HandleFunc("/", handlers.MainHandler).Methods("GET")
	r.HandleFunc("/createRoom", handlers.CreateRoomHandler).Methods("GET")
	r.HandleFunc("/game", handlers.GameHandler).Methods("GET")

	r.HandleFunc("/ws", WScontroller.WebsocketHandler)

	apiHandlers := handlers.NewAPIHandlers(GameServer)
	secretKey := handlers.GetJwt()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/createRoom", middleware.AuthMiddleware(secretKey, apiHandlers.CreateRoom)).Methods("POST")
	api.HandleFunc("/getRooms", middleware.AuthMiddleware(secretKey, apiHandlers.GetRooms)).Methods("GET")
	api.HandleFunc("/getRoom", middleware.AuthMiddleware(secretKey, apiHandlers.GetRoomByName)).Methods("GET")

	addr := "0.0.0.0:" + strconv.Itoa(*port)

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("GAME SERVER RUNNING ON: %s", srv.Addr)

		if _, err := os.Stat("./static"); err != nil {
			log.Printf("WARNING: Static directory not found: %v", err)
		} else {
			log.Printf("Static directory found")
		}

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
