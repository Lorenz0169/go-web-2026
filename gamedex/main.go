// GameDex : le catalogue de jeux et de leurs tournois (mini TP 3 : appel de la mini API).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

// afficherCatalogue écrit dans le terminal chaque jeu avec son nombre de tournois.
func afficherCatalogue() {
	// On récupère les jeux ; en cas d'erreur, le programme s'arrête
	games, err := FetchGames()
	if err != nil {
		log.Fatal("Erreur avec la mini API (est-elle lancée avec go run ./api ?) : ", err)
	}

	// Même chose pour les tournois
	tournaments, err := FetchTournaments()
	if err != nil {
		log.Fatal("Erreur avec la mini API (est-elle lancée avec go run ./api ?) : ", err)
	}

	// On range les tournois par jeu : la clé de la map est l'identifiant du jeu
	byGame := make(map[int][]Tournament)
	for _, t := range tournaments {
		byGame[t.GameID] = append(byGame[t.GameID], t)
	}

	// Pour chaque jeu : son nom, son année et son nombre de tournois
	for _, g := range games {
		fmt.Printf("%s (%d) : %d tournoi(s)\n", g.Name, g.Year, len(byGame[g.ID]))
	}
}

// accueil répond sur "GET /{$}" : la page d'accueil, en HTML.
func accueil(w http.ResponseWriter, r *http.Request) {
	// On prévient le navigateur qu'on envoie du HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// La page utilise la feuille de style servie par /static/
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="utf-8">
    <title>GameDex</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <h1>Bienvenue sur GameDex</h1>
    <p>Le catalogue de jeux et de leurs tournois.</p>
    <p><a href="/jeux">Voir la liste des jeux</a></p>
</body>
</html>`)
}

// listeJeux répond sur "GET /jeux" : la liste des jeux, en texte.
func listeJeux(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Liste des jeux (bientôt le vrai catalogue)")
}

// detailJeu répond sur "GET /jeux/{id}" : la fiche d'un jeu.
func detailJeu(w http.ResponseWriter, r *http.Request) {
	// On lit l'identifiant dans l'URL (c'est du texte)
	texte := r.PathValue("id")

	// On le convertit en nombre ; si ce n'est pas un nombre : erreur 400
	id, err := strconv.Atoi(texte)
	if err != nil {
		http.Error(w, "l'identifiant doit être un nombre", http.StatusBadRequest)
		return
	}

	// Un identifiant plus petit que 1 n'a pas de sens : erreur 400 aussi
	if id < 1 {
		http.Error(w, "l'identifiant doit être supérieur ou égal à 1", http.StatusBadRequest)
		return
	}

	// Les jeux vont de 1 à 12 ; au-delà, le jeu n'existe pas : erreur 404
	if id > 12 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Fiche du jeu numéro %d\n", id)
}

func main() {
	// Option -port : go run . -port 9000 (8080 par défaut)
	port := flag.Int("port", 8080, "port d'écoute du serveur")
	flag.Parse()

	// Avant de lancer le serveur, on affiche le catalogue dans le terminal
	afficherCatalogue()

	// Le routeur : il choisit le handler selon la méthode et l'URL
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /jeux", listeJeux)
	mux.HandleFunc("GET /jeux/{id}", detailJeu)

	// Fichiers statiques : /static/css/style.css -> le fichier static/css/style.css
	// (le chemin "static" est relatif au dossier de lancement)
	fichiers := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fichiers))

	// Adresse d'écoute, par exemple ":8080"
	adresse := ":" + strconv.Itoa(*port)
	fmt.Println("GameDex écoute sur http://localhost" + adresse)

	// ListenAndServe ne rend la main qu'en cas d'erreur
	err := http.ListenAndServe(adresse, mux)
	log.Fatal(err)
}
