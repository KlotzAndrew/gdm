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

import "gdms/lib"

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
    mux.HandleFunc(pat.Get("/games"), allGames())
    mux.HandleFunc(pat.Post("/games"), addGame())
    mux.HandleFunc(pat.Get("/games/:game_id"), gameById())
    mux.HandleFunc(pat.Post("/games/:game_id"), deleteGame())
    http.ListenAndServe("localhost:8080", mux)
}

func allGames() func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        games, err := repo.AllGames()
        if err != nil {
            ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all games: ", err)
            return
        }

        respBody, err := json.MarshalIndent(games, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func addGame() func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        var game repo.Game
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&game)
        if err != nil {
            ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        err = repo.AddGame(game)

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

func gameById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		GameId := pat.Param(r, "game_id")

    game, err := repo.GameById(GameId)
    if err != nil {
    	ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
    	log.Println("Failed find game: ", err)
    	return
    }

		if game.GameId == "" {
			ErrorWithJSON(w, "Game not found", http.StatusNotFound)
			return
		}

		respBody, err := json.MarshalIndent(game, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func deleteGame() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		GameId := pat.Param(r, "game_id")

    err := repo.DeleteGame(GameId)

		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed delete game: ", err)
				return
			case mgo.ErrNotFound:
				ErrorWithJSON(w, "Game not found", http.StatusNotFound)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
