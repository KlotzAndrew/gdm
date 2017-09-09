package main

import (
    "fmt"
    "net/http"
    "log"
    "encoding/json"

    "goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
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
  session, err := mgo.Dial("localhost")
    if err != nil { panic(err) }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    ensureIndex(session)

    mux := goji.NewMux()
    mux.HandleFunc(pat.Post("/games"), addGame(session))
    http.ListenAndServe("localhost:8080", mux)
}

func ensureIndex(s *mgo.Session) {
    session := s.Copy()
    defer session.Close()

    c := session.DB("store").C("games")

    index := mgo.Index{
        Key:        []string{"game_id"},
        Unique:     true,
        DropDups:   true,
        Background: true,
        Sparse:     true,
    }
    err := c.EnsureIndex(index)
    if err != nil { panic(err) }
}

type Game struct {
    GameId   string `json:"game_id"`
    PlayerId string `json:"player_id"`
    Meta      string `json:"meta"`
}

func addGame(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var game Game
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&game)
        if err != nil {
            ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("store").C("games")

        err = c.Insert(game)
        if err != nil {
            if mgo.IsDup(err) {
                ErrorWithJSON(w, "Game with this game_id already exists", http.StatusBadRequest)
                return
            }

            ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed insert game: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
    }
}
