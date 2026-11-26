// GameDex : le catalogue de jeux et de leurs tournois (mini TP 4 : les vraies pages).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
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

// JeuPage contient les données de la page d'un jeu.
type JeuPage struct {
	Game        Game
	Tournaments []Tournament
}

// ListPage contient les données de la page de la liste des jeux.
type ListPage struct {
	Games []Game
	Query string // le texte de ?q=
}

// ErreurPage contient les données de la page d'erreur.
type ErreurPage struct {
	Code    int
	Message string
}

// On charge tous les templates une seule fois.
// template.Must arrête le programme si un template contient une erreur.
var templates = template.Must(template.ParseGlob("templates/*.html"))

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

// afficherPage remplit un template avec des données et l'envoie au navigateur.
func afficherPage(w http.ResponseWriter, code int, nom string, donnees any) {
	// Les en-têtes d'abord, puis le code, et enfin la page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)

	err := templates.ExecuteTemplate(w, nom, donnees)
	if err != nil {
		// Le code est déjà envoyé : on note seulement l'erreur dans le terminal
		log.Println("erreur de template :", err)
	}
}

// afficherErreur envoie la page d'erreur avec le bon code HTTP.
func afficherErreur(w http.ResponseWriter, code int, message string) {
	page := ErreurPage{Code: code, Message: message}
	afficherPage(w, code, "erreur.html", page)
}

// accueil répond sur "GET /{$}" : la page d'accueil.
func accueil(w http.ResponseWriter, r *http.Request) {
	afficherPage(w, http.StatusOK, "accueil.html", nil)
}

// listeJeux répond sur "GET /jeux" : la liste des jeux.
func listeJeux(w http.ResponseWriter, r *http.Request) {
	// On demande les jeux à l'API
	games, err := FetchGames()
	if err != nil {
		log.Println(err)
		afficherErreur(w, http.StatusBadGateway, "La liste des jeux est indisponible : l'API ne répond pas correctement.")
		return
	}

	page := ListPage{
		Games: games,
		Query: r.URL.Query().Get("q"),
	}
	afficherPage(w, http.StatusOK, "liste.html", page)
}

// detailJeu répond sur "GET /jeux/{id}" : la fiche d'un jeu.
func detailJeu(w http.ResponseWriter, r *http.Request) {
	// 1. On lit l'identifiant dans l'URL : il doit être un nombre à partir de 1
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		afficherErreur(w, http.StatusBadRequest, "L'identifiant doit être un nombre supérieur ou égal à 1.")
		return
	}

	// 2. On demande les jeux à l'API
	games, err := FetchGames()
	if err != nil {
		log.Println(err)
		afficherErreur(w, http.StatusBadGateway, "La fiche est indisponible : l'API ne répond pas correctement.")
		return
	}

	// 3. On cherche le jeu qui a cet identifiant
	var page JeuPage
	trouve := false
	for _, g := range games {
		if g.ID == id {
			page.Game = g
			trouve = true
		}
	}
	if !trouve {
		afficherErreur(w, http.StatusNotFound, "Jeu introuvable.")
		return
	}

	// 4. On garde seulement les tournois de ce jeu
	tournaments, err := FetchTournaments()
	if err != nil {
		log.Println(err)
		afficherErreur(w, http.StatusBadGateway, "Les tournois sont indisponibles : l'API ne répond pas correctement.")
		return
	}
	for _, t := range tournaments {
		if t.GameID == id {
			page.Tournaments = append(page.Tournaments, t)
		}
	}

	afficherPage(w, http.StatusOK, "jeu.html", page)
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
