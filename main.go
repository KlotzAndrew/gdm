package main

import (
    "fmt"
    "net/http"

    "goji.io"
    "goji.io/pat"
)

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    w.Write(json)
}

func main() {
    mux := goji.NewMux()
    mux.HandleFunc(pat.Post("/games"), addGame())
    http.ListenAndServe("localhost:8080", mux)
}

func addGame() func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
      ResponseWithJSON(w, []byte(`{msg: "coming soon!"}`), http.StatusBadRequest)
    }
}
