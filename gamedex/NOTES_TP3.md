# Notes du mini TP 3

Pour les deux premiers tests, on remplace `/api/games` dans la constante `gamesURL`, puis on remet la bonne URL.

## `/api/broken`

```
Erreur avec la mini API (est-elle lancée avec go run ./api ?) : appel des jeux : l'API a répondu le code 500
```

La requête aboutit (`err` vaut `nil`) : c'est la vérification de `resp.StatusCode != 200` qui détecte le problème.

## `/api/slow`

```
Erreur avec la mini API (est-elle lancée avec go run ./api ?) : appel des jeux impossible : Get "http://localhost:8081/api/slow": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

L'API répond après 5 secondes : le timeout de 3 secondes du client se déclenche avant.

## API coupée

```
Erreur avec la mini API (est-elle lancée avec go run ./api ?) : appel des jeux impossible : Get "http://localhost:8081/api/games": dial tcp [::1]:8081: connectex: Aucune connexion n’a pu être établie car l’ordinateur cible l’a expressément refusée.
```

Dans les trois cas, `log.Fatal` affiche le message et arrête le programme, sans panique.
