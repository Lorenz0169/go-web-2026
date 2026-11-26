// Package store garde en mémoire ce que font les visiteurs.
// Tout est perdu quand on arrête GameDex.
//
// Attention : ce n'est PAS le fichier snippets/likes.go du cours, qui était
// introuvable. C'est une version reconstruite à partir de l'énoncé du mini TP 6,
// à remplacer par le fichier officiel s'il est fourni.
package store

import "sync"

// Le serveur peut traiter plusieurs requêtes en même temps :
// le verrou empêche deux requêtes de toucher à la map au même moment.
var likesMutex sync.Mutex

// La clé est l'identifiant du jeu, la valeur son nombre de likes
var likes = make(map[int]int)

// AddLike ajoute un like à un jeu.
func AddLike(id int) {
	likesMutex.Lock()
	defer likesMutex.Unlock()
	likes[id]++
}

// Likes renvoie le nombre de likes d'un jeu (0 s'il n'en a pas).
func Likes(id int) int {
	likesMutex.Lock()
	defer likesMutex.Unlock()
	return likes[id]
}
