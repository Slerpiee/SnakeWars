package main

import(
	"net/http"
	"github.com/gorilla/mux"
	"time"
	""
)
func main() {


	
    r := mux.NewRouter()


    srv := &http.Server{
        Handler:      r,
        Addr:         "127.0.0.1:8000",
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    go func(){
		go srv.ListenAndServe()
	}()

}
