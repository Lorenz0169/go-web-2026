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

GameDex charge les jeux et les tournois une seule fois, au démarrage : l'API doit être lancée avant lui.

## Routes

| Route | Ce qu'elle fait |
| --- | --- |
| `GET /` | La page d'accueil |
| `GET /jeux` | La liste des jeux ; `?q=` cherche dans le nom |
| `GET /jeux/{id}` | La fiche d'un jeu et ses tournois (400 si l'id est invalide, 404 si le jeu est inconnu) |
| `GET /static/` | Les fichiers statiques (CSS) |

## Organisation

```
cmd/gamedex/       le programme : charge les données et lance le serveur
internal/api/      les structs et les appels à la mini API
internal/catalog/  les règles de filtrage
internal/server/   les routes, les handlers et les templates
web/templates/     les pages HTML
web/static/        le CSS
```
