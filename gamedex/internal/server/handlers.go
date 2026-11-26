package server

import (
	"net/http"
	"strconv"

	"gamedex/internal/api"
	"gamedex/internal/catalog"
)

// ListPage contient les données de la page de la liste des jeux.
type ListPage struct {
	Games []api.Game
	Query string // le texte de ?q=
}

// JeuPage contient les données de la page d'un jeu.
type JeuPage struct {
	Game        api.Game
	Tournaments []api.Tournament
}

// ErreurPage contient les données de la page d'erreur.
type ErreurPage struct {
	Code    int
	Message string
}

// accueil répond sur "GET /{$}" : la page d'accueil.
func (s *Server) accueil(w http.ResponseWriter, r *http.Request) {
	afficherPage(w, http.StatusOK, "accueil.html", nil)
}

// listeJeux répond sur "GET /jeux" : la liste des jeux, filtrée par ?q=.
func (s *Server) listeJeux(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	page := ListPage{
		Games: catalog.Filter(s.Games, "", query),
		Query: query,
	}
	afficherPage(w, http.StatusOK, "liste.html", page)
}

// detailJeu répond sur "GET /jeux/{id}" : la fiche d'un jeu.
func (s *Server) detailJeu(w http.ResponseWriter, r *http.Request) {
	// L'identifiant doit être un nombre à partir de 1
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		afficherErreur(w, http.StatusBadRequest, "L'identifiant doit être un nombre supérieur ou égal à 1.")
		return
	}

	game, trouve := s.findGame(id)
	if !trouve {
		afficherErreur(w, http.StatusNotFound, "Jeu introuvable.")
		return
	}

	page := JeuPage{
		Game:        game,
		Tournaments: s.tournamentsOf(id),
	}
	afficherPage(w, http.StatusOK, "jeu.html", page)
}

// findGame cherche un jeu par son identifiant. Le booléen dit s'il existe.
func (s *Server) findGame(id int) (api.Game, bool) {
	for _, g := range s.Games {
		if g.ID == id {
			return g, true
		}
	}
	return api.Game{}, false
}

// tournamentsOf renvoie les tournois d'un jeu.
func (s *Server) tournamentsOf(gameID int) []api.Tournament {
	var result []api.Tournament
	for _, t := range s.Tournaments {
		if t.GameID == gameID {
			result = append(result, t)
		}
	}
	return result
}
