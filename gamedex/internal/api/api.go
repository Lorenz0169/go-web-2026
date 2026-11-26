// Package api parle à la mini API : les données et les appels HTTP.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Game est un jeu renvoyé par l'API.
// Les tags `json:"..."` donnent le nom de la clé dans le JSON.
type Game struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Year   int      `json:"year"`
	Studio string   `json:"studio"`
	Genres []string `json:"genres"`
}

// Tournament est un tournoi renvoyé par l'API.
// GameID (clé "gameId") indique à quel jeu appartient le tournoi.
type Tournament struct {
	ID     int     `json:"id"`
	GameID int     `json:"gameId"`
	Name   string  `json:"name"`
	City   string  `json:"city"`
	Date   string  `json:"date"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
}

// Les adresses de la mini API
const gamesURL = "http://localhost:8081/api/games"
const tournamentsURL = "http://localhost:8081/api/tournaments"

// Le client HTTP abandonne si l'API met plus de 3 secondes à répondre
var client = &http.Client{Timeout: 3 * time.Second}

// FetchGames demande la liste des jeux à l'API.
func FetchGames() ([]Game, error) {
	// 1. On envoie la requête
	resp, err := client.Get(gamesURL)
	if err != nil {
		return nil, fmt.Errorf("appel des jeux impossible : %w", err)
	}
	// Le corps de la réponse sera fermé à la fin de la fonction
	defer resp.Body.Close()

	// 2. On vérifie le code de réponse
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("appel des jeux : l'API a répondu le code %d", resp.StatusCode)
	}

	// 3. On décode le JSON dans une liste de jeux
	var games []Game
	err = json.NewDecoder(resp.Body).Decode(&games)
	if err != nil {
		return nil, fmt.Errorf("JSON des jeux illisible : %w", err)
	}
	return games, nil
}

// SearchGames demande à l'API les jeux filtrés par nom (q) et par genre.
// Un filtre vide ("") est ignoré par l'API.
func SearchGames(q, genre string) ([]Game, error) {
	// url.Values construit la query string et protège les caractères spéciaux
	params := url.Values{}
	params.Set("q", q)
	params.Set("genre", genre)

	// 1. On envoie la requête
	resp, err := client.Get(gamesURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("recherche des jeux impossible : %w", err)
	}
	defer resp.Body.Close()

	// 2. On vérifie le code de réponse
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("recherche des jeux : l'API a répondu le code %d", resp.StatusCode)
	}

	// 3. On décode le JSON dans une liste de jeux
	var games []Game
	err = json.NewDecoder(resp.Body).Decode(&games)
	if err != nil {
		return nil, fmt.Errorf("JSON de la recherche illisible : %w", err)
	}
	return games, nil
}

// FetchTournaments demande la liste de tous les tournois à l'API.
func FetchTournaments() ([]Tournament, error) {
	// 1. On envoie la requête
	resp, err := client.Get(tournamentsURL)
	if err != nil {
		return nil, fmt.Errorf("appel des tournois impossible : %w", err)
	}
	defer resp.Body.Close()

	// 2. On vérifie le code de réponse
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("appel des tournois : l'API a répondu le code %d", resp.StatusCode)
	}

	// 3. On décode le JSON dans une liste de tournois
	var tournaments []Tournament
	err = json.NewDecoder(resp.Body).Decode(&tournaments)
	if err != nil {
		return nil, fmt.Errorf("JSON des tournois illisible : %w", err)
	}
	return tournaments, nil
}
