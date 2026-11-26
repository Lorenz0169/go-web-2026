package server

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"gamedex/internal/api"
	"gamedex/internal/catalog"
	"gamedex/internal/store"
)

// ListPage contient les données de la page de la liste des jeux.
type ListPage struct {
	Games  []api.Game
	Query  string   // le texte de ?q=
	Genre  string   // le genre choisi (?genre=)
	Genres []string // tous les genres, pour la liste déroulante
}

// JeuPage contient les données de la page d'un jeu.
type JeuPage struct {
	Game        api.Game
	Tournaments []api.Tournament
	Likes       int
	Inscrits    map[int][]string // identifiant du tournoi -> pseudos inscrits
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

// listeJeux répond sur "GET /jeux" : la liste des jeux, filtrée par ?q= et ?genre=.
func (s *Server) listeJeux(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	genre := r.URL.Query().Get("genre")

	// C'est la mini API qui fait la recherche
	games, err := api.SearchGames(query, genre)
	if err != nil {
		log.Println(err)
		afficherErreur(w, http.StatusBadGateway, "La liste des jeux est indisponible : l'API ne répond pas correctement.")
		return
	}

	// Les jeux les plus likés en premier ; à égalité, l'ordre ne change pas
	slices.SortStableFunc(games, func(a, b api.Game) int {
		return store.Likes(b.ID) - store.Likes(a.ID)
	})

	page := ListPage{
		Games:  games,
		Query:  query,
		Genre:  genre,
		Genres: catalog.Genres(s.Games),
	}
	afficherPage(w, http.StatusOK, "liste.html", page)
}

// detailJeu répond sur "GET /jeux/{id}" : la fiche d'un jeu.
func (s *Server) detailJeu(w http.ResponseWriter, r *http.Request) {
	game, ok := s.lireJeu(w, r)
	if !ok {
		return
	}

	page := JeuPage{
		Game:        game,
		Tournaments: s.tournamentsOf(game.ID),
		Likes:       store.Likes(game.ID),
		Inscrits:    make(map[int][]string),
	}
	for _, t := range page.Tournaments {
		page.Inscrits[t.ID] = store.Inscriptions(t.ID)
	}
	afficherPage(w, http.StatusOK, "jeu.html", page)
}

// likeJeu répond sur "POST /jeux/{id}/like" : ajoute un like au jeu.
func (s *Server) likeJeu(w http.ResponseWriter, r *http.Request) {
	game, ok := s.lireJeu(w, r)
	if !ok {
		return
	}

	store.AddLike(game.ID)

	// 303 : le navigateur va ensuite chercher la fiche avec un GET.
	// Un F5 recharge ce GET, il ne renvoie pas le formulaire.
	http.Redirect(w, r, "/jeux/"+strconv.Itoa(game.ID), http.StatusSeeOther)
}

// inscrireTournoi répond sur "POST /tournois/{id}/inscription".
func (s *Server) inscrireTournoi(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		afficherErreur(w, http.StatusBadRequest, "L'identifiant doit être un nombre supérieur ou égal à 1.")
		return
	}

	tournament, trouve := s.findTournament(id)
	if !trouve {
		afficherErreur(w, http.StatusNotFound, "Tournoi introuvable.")
		return
	}

	// On lit le champ "pseudo" du formulaire, sans les espaces au début et à la fin
	pseudo := strings.TrimSpace(r.FormValue("pseudo"))
	if pseudo == "" || len(pseudo) > 30 {
		afficherErreur(w, http.StatusBadRequest, "Le pseudo est obligatoire (30 caractères au maximum).")
		return
	}

	store.AddInscription(tournament.ID, pseudo)

	// Comme pour les likes : redirection 303 vers la fiche du jeu
	http.Redirect(w, r, "/jeux/"+strconv.Itoa(tournament.GameID), http.StatusSeeOther)
}

// lireJeu lit l'identifiant dans l'URL et cherche le jeu.
// En cas de problème, elle affiche l'erreur et renvoie false.
func (s *Server) lireJeu(w http.ResponseWriter, r *http.Request) (api.Game, bool) {
	// L'identifiant doit être un nombre à partir de 1
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		afficherErreur(w, http.StatusBadRequest, "L'identifiant doit être un nombre supérieur ou égal à 1.")
		return api.Game{}, false
	}

	game, trouve := s.findGame(id)
	if !trouve {
		afficherErreur(w, http.StatusNotFound, "Jeu introuvable.")
		return api.Game{}, false
	}
	return game, true
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

// findTournament cherche un tournoi par son identifiant.
func (s *Server) findTournament(id int) (api.Tournament, bool) {
	for _, t := range s.Tournaments {
		if t.ID == id {
			return t, true
		}
	}
	return api.Tournament{}, false
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
