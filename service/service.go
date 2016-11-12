package main

import (
    "log"
    "net/http"
    "encoding/json"
    "time"
)

type Response struct {
    Message string `json:"message"`
}

type Toggle struct {
    errorMode bool
}

func main() {
    var h = Toggle{false}

    log.Println("Starting server...")
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            if h.errorMode  == true {
                time.Sleep(3000 * time.Millisecond)
                http.Error(w, "Fatal Error", 500)
                return
            } else {
                time.Sleep(1000 * time.Millisecond)
                w.Header().Set("Content-Type", "application/json")
                json.NewEncoder(w).Encode(Response{Message: "Success!"})
                return
            }
        } else {
            h.errorMode = !h.errorMode
            log.Println("Toggled error mode")
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(Response{Message: "Toggled error mode"})
        }
    })
    log.Fatal(http.ListenAndServe(":8000", mux))
    log.Printf("HTTP service listening on 8000")
}