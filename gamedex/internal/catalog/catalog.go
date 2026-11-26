// Package catalog contient les règles du catalogue.
package catalog

import (
	"slices"
	"strings"

	"gamedex/internal/api"
)

// Genres renvoie tous les genres des jeux, sans doublon, par ordre alphabétique.
func Genres(games []api.Game) []string {
	var result []string
	for _, g := range games {
		for _, genre := range g.Genres {
			// On ajoute le genre seulement s'il n'est pas déjà dans la liste
			if !slices.Contains(result, genre) {
				result = append(result, genre)
			}
		}
	}
	slices.Sort(result)
	return result
}

// Filter garde les jeux qui ont le bon genre et dont le nom contient query.
// Un filtre vide ("") ne retire aucun jeu.
func Filter(games []api.Game, genre, query string) []api.Game {
	// On met le texte cherché en minuscules, une seule fois
	query = strings.ToLower(query)

	result := []api.Game{}
	for _, g := range games {
		// Le jeu n'a pas le genre demandé : on passe au suivant
		if genre != "" && !slices.Contains(g.Genres, genre) {
			continue
		}
		// Le nom (en minuscules) ne contient pas le texte cherché : on passe
		if query != "" && !strings.Contains(strings.ToLower(g.Name), query) {
			continue
		}
		result = append(result, g)
	}
	return result
}
