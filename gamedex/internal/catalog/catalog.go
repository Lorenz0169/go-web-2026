// Package catalog contient les règles du catalogue.
package catalog

import (
	"slices"
	"strings"

	"gamedex/internal/api"
)

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
