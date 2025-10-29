package main

import(
	"net/http"
    "game_server/server/handlers"
    "os"
    "os/signal"
    "syscall"
    "time"
    "log"
    "context"

    "github.com/gorilla/mux"
)
func main() {

    r := mux.NewRouter()

    r.HandleFunc("/", handlers.MainHandler).Methods("GET")
    r.HandleFunc("/createRoom", handlers.CreateRoomHandler).Methods("GET")
    r.HandleFunc("/game", handlers.GameHandler).Methods("GET")

    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", 
		http.FileServer(http.Dir("./static"))))


    srv := &http.Server{
        Handler:      r,
        Addr:         "127.0.0.1:8000",
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
