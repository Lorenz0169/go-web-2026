// Package server affiche les pages de GameDex.
package server

import (
	"html/template"
	"log"
	"net/http"

	"gamedex/internal/api"
)

// Server garde en mémoire les données chargées au démarrage.
type Server struct {
	Games       []api.Game
	Tournaments []api.Tournament
}

// On charge tous les templates une seule fois.
// Le chemin est relatif au dossier gamedex/, d'où on lance le programme.
// template.Must arrête le programme si un template contient une erreur.
var templates = template.Must(template.ParseGlob("web/templates/*.html"))

// Routes crée le routeur : il choisit le handler selon la méthode et l'URL.
func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.accueil)
	mux.HandleFunc("GET /jeux", s.listeJeux)
	mux.HandleFunc("GET /jeux/{id}", s.detailJeu)
	mux.HandleFunc("POST /jeux/{id}/like", s.likeJeu)
	mux.HandleFunc("POST /tournois/{id}/inscription", s.inscrireTournoi)

	// Fichiers statiques : /static/css/style.css -> web/static/css/style.css
	fichiers := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fichiers))

	return mux
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
