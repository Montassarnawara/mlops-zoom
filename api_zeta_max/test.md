# ✅ Rechargement des données
curl http://localhost:3000/api/reload

# ✅ Accès à l'hôtel Carter Group dans la ville "New Mia"  
curl http://localhost:3000/api/hotels/HOT317

# ✅ Recherche d'hôtels par ville
curl -s "http://localhost:3000/api/hotels" | jq '.data[] | select(.city | contains("New Mia"))'

