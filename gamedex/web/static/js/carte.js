// La carte des tournois, avec Leaflet.
// La variable "tournois" est créée par le template jeu.html, avant ce fichier.
// Les clés sont celles des tags JSON de Go : name, city, lat, lng.

// 1. On crée la carte dans la balise <div id="carte">
const carte = L.map("carte");

// 2. On ajoute le fond de carte OpenStreetMap
L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
}).addTo(carte);

// 3. Un marqueur et une bulle par tournoi
const positions = [];
for (const tournoi of tournois) {
    const position = [tournoi.lat, tournoi.lng];
    positions.push(position);

    // textContent écrit du texte simple : rien n'est interprété comme du HTML
    const bulle = document.createElement("div");
    bulle.textContent = tournoi.name + " · " + tournoi.city;

    L.marker(position).addTo(carte).bindPopup(bulle);
}

// 4. On règle la vue pour voir tous les marqueurs, sans zoomer trop près
carte.fitBounds(positions, { padding: [30, 30], maxZoom: 6 });
