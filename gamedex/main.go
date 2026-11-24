// GameDex : le catalogue de jeux et de leurs tournois (mini TP 2 : premier serveur).
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

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
