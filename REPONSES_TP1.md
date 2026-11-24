# Mini TP 1 : exploration de la mini API

API lancée avec `go run ./api` (écoute sur `localhost:8081`). Les commandes ont été faites avec `curl.exe`. Les valeurs ci-dessous sont celles réellement obtenues.

## 1. Codes de `/api/games/999` et `/api/games/abc`

| Commande | Code HTTP | Réponse |
| --- | --- | --- |
| `curl.exe -i http://localhost:8081/api/games/999` | **404 Not Found** | `{"error": "jeu introuvable"}` |
| `curl.exe -i http://localhost:8081/api/games/abc` | **400 Bad Request** | `{"error": "l'identifiant doit être un nombre"}` |

Le jeu 999 n'existe pas (404). « abc » n'est pas un nombre : la requête est mal formée (400).

## 2. Tournois du jeu 4 (Street Fighter 6)

Commande : `curl.exe -i http://localhost:8081/api/games/4/tournaments`
Code HTTP : **200 OK**

Il y a **4 tournois** :

| Tournoi | Ville |
| --- | --- |
| Fight Night Marseille | Marseille |
| Combo Breaker EU | Paris |
| Rising Dragon | Tokyo |
| Hadoken Series | Bruxelles |

## 3. Filtres sur `/api/games`

Les paramètres sont `genre` et `year` (d'après le README).

- Jeux de course : `curl.exe -i "http://localhost:8081/api/games?genre=course"`
  Code **200 OK**, 3 jeux : Rocket League, Trackmania, Mario Kart 8 Deluxe.
- Jeux sortis en 2018 : `curl.exe -i "http://localhost:8081/api/games?year=2018"`
  Code **200 OK**, 2 jeux : Celeste, Dead Cells.

## 4. `/api/stats` et la clé d'API

La clé se passe dans l'en-tête **`X-API-Key`** (d'après le README et le code de l'API).

| Cas | Commande | Code HTTP | Réponse |
| --- | --- | --- | --- |
| Sans clé | `curl.exe -i http://localhost:8081/api/stats` | **401 Unauthorized** | `{"error": "clé absente ou inconnue : ajoutez l'en-tête X-API-Key"}` |
| Clé `eleve` | `curl.exe -i -H "X-API-Key: eleve" http://localhost:8081/api/stats` | **403 Forbidden** | `{"error": "clé reconnue, mais pas le droit de voir les statistiques"}` |
| Clé `prof` | `curl.exe -i -H "X-API-Key: prof" http://localhost:8081/api/stats` | **200 OK** | `{"games": 12, "tournaments": 20}` |

401 : le serveur ne sait pas qui on est. 403 : il sait qui on est, mais on n'a pas le droit.

## 5. Bonus

- `curl.exe -i http://localhost:8081/api/broken` : **500 Internal Server Error**, `{"error": "panne simulée"}`.
- `curl.exe -s -o /dev/null -w "code=%{http_code} temps=%{time_total}s\n" http://localhost:8081/api/slow` : **200 OK** après **5,001794 s** (`code=200 temps=5.001794s`).
