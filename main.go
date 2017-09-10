package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"goji.io"
	"goji.io/pat"
	"gopkg.in/mgo.v2"
)

import "gds/lib"

func errorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func responseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/games"), allGames())
	mux.HandleFunc(pat.Post("/games"), addGame())
	mux.HandleFunc(pat.Get("/games/:game_id"), gameByID())
	mux.HandleFunc(pat.Post("/games/:game_id"), deleteGame())
	http.ListenAndServe("localhost:8080", mux)
}

func allGames() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		games, err := repo.AllGames()
		if err != nil {
			errorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all games: ", err)
			return
		}

		respBody, err := json.MarshalIndent(games, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		responseWithJSON(w, respBody, http.StatusOK)
	}
}

func addGame() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var game repo.Game
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&game)
		if err != nil {
			errorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		err = repo.AddGame(game)

		if err != nil {
			if mgo.IsDup(err) {
				errorWithJSON(w, "Game with this game_id already exists", http.StatusBadRequest)
				return
			}

			errorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert game: ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

func gameByID() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		GameID := pat.Param(r, "game_id")

		game, err := repo.GameByID(GameID)
		if err != nil {
			errorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed find game: ", err)
			return
		}

		if game.GameID == "" {
			errorWithJSON(w, "Game not found", http.StatusNotFound)
			return
		}

		respBody, err := json.MarshalIndent(game, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		responseWithJSON(w, respBody, http.StatusOK)
	}
}

func deleteGame() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		GameID := pat.Param(r, "game_id")

		err := repo.DeleteGame(GameID)

		if err != nil {
			switch err {
			default:
				errorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed delete game: ", err)
				return
			case mgo.ErrNotFound:
				errorWithJSON(w, "Game not found", http.StatusNotFound)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
