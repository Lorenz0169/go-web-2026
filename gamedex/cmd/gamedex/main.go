// GameDex : le catalogue de jeux et de leurs tournois (mini TP 5 : code rangé).
// Lancement, depuis le dossier gamedex/ : go run ./cmd/gamedex
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gamedex/internal/api"
	"gamedex/internal/server"
)

// afficherCatalogue écrit dans le terminal chaque jeu avec son nombre de tournois.
func afficherCatalogue(games []api.Game, tournaments []api.Tournament) {
	// On range les tournois par jeu : la clé de la map est l'identifiant du jeu
	byGame := make(map[int][]api.Tournament)
	for _, t := range tournaments {
		byGame[t.GameID] = append(byGame[t.GameID], t)
	}

	// Pour chaque jeu : son nom, son année et son nombre de tournois
	for _, g := range games {
		fmt.Printf("%s (%d) : %d tournoi(s)\n", g.Name, g.Year, len(byGame[g.ID]))
	}
}

func main() {
	// Option -port : go run ./cmd/gamedex -port 9000 (8080 par défaut)
	port := flag.Int("port", 8080, "port d'écoute du serveur")
	flag.Parse()

	// 1. On charge les données une seule fois, au démarrage
	games, err := api.FetchGames()
	if err != nil {
		log.Fatal("Erreur avec la mini API (est-elle lancée avec go run ./api ?) : ", err)
	}
	tournaments, err := api.FetchTournaments()
	if err != nil {
		log.Fatal("Erreur avec la mini API (est-elle lancée avec go run ./api ?) : ", err)
	}
	afficherCatalogue(games, tournaments)

	// 2. On assemble le serveur avec ces données
	s := &server.Server{
		Games:       games,
		Tournaments: tournaments,
	}

	// 3. On lance le serveur
	adresse := ":" + strconv.Itoa(*port)
	fmt.Println("GameDex écoute sur http://localhost" + adresse)

	// ListenAndServe ne rend la main qu'en cas d'erreur
	err = http.ListenAndServe(adresse, s.Routes())
	log.Fatal(err)
}
