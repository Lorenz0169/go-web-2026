# GameDex

Un catalogue de jeux vidéo et de leurs tournois, écrit en Go.

## Lancer

Il faut deux terminaux.

1. La mini API (port 8081), depuis la racine du dépôt :

```
go run ./api
```

2. GameDex (port 8080), depuis le dossier `gamedex/` :

```
go run ./cmd/gamedex
```

Pour changer de port : `go run ./cmd/gamedex -port 9000`.

GameDex charge les jeux et les tournois une seule fois, au démarrage : l'API doit être lancée avant lui. La recherche de `/jeux` rappelle l'API à chaque fois. Les likes et les inscriptions sont gardés en mémoire : ils disparaissent à l'arrêt de GameDex.

## Routes

| Route | Ce qu'elle fait |
| --- | --- |
| `GET /` | La page d'accueil |
| `GET /jeux` | La liste des jeux, les plus likés en premier ; `?q=` et `?genre=` sont transmis à la mini API (502 si elle ne répond pas) |
| `GET /jeux/{id}` | La fiche d'un jeu et ses tournois, avec une carte Leaflet s'il y en a (400 si l'id est invalide, 404 si le jeu est inconnu) |
| `POST /jeux/{id}/like` | Ajoute un like puis redirige vers la fiche (303) |
| `POST /tournois/{id}/inscription` | Inscrit le pseudo envoyé par le formulaire puis redirige (303) |
| `GET /static/` | Les fichiers statiques (CSS et JavaScript) |

## Organisation

```
cmd/gamedex/       le programme : charge les données et lance le serveur
internal/api/      les structs et les appels à la mini API
internal/catalog/  les règles de filtrage
internal/store/    les likes et les inscriptions, en mémoire
internal/server/   les routes, les handlers et les templates
web/templates/     les pages HTML
web/static/        le CSS et le JavaScript de la carte
```
