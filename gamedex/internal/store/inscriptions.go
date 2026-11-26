package store

import (
	"slices"
	"sync"
)

// Les inscriptions suivent le même modèle que les likes : une map et un verrou
var inscriptionsMutex sync.Mutex

// La clé est l'identifiant du tournoi, la valeur la liste des pseudos inscrits
var inscriptions = make(map[int][]string)

// AddInscription inscrit un pseudo à un tournoi.
// Un pseudo déjà inscrit n'est pas ajouté une deuxième fois.
func AddInscription(tournamentID int, pseudo string) {
	inscriptionsMutex.Lock()
	defer inscriptionsMutex.Unlock()

	if slices.Contains(inscriptions[tournamentID], pseudo) {
		return
	}
	inscriptions[tournamentID] = append(inscriptions[tournamentID], pseudo)
}

// Inscriptions renvoie les pseudos inscrits à un tournoi.
func Inscriptions(tournamentID int) []string {
	inscriptionsMutex.Lock()
	defer inscriptionsMutex.Unlock()

	// On renvoie une copie : la liste d'origine reste protégée par le verrou
	return slices.Clone(inscriptions[tournamentID])
}
