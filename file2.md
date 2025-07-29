docker run -d --name hazelcast \
  -p 5701:5701 \
  -e HZ_CLUSTERNAME=dev-go \
  hazelcast/hazelcast:5.3



config.Cluster.Network.SetAddresses("127.0.0.1:5701")
docker ps
docker logs hazelcast

# Lancer le fichier Go
go run main.go

# Voir les logs de Hazelcast (optionnel pour debug)
docker logs -f hazelcast

# Arrêter Hazelcast
docker stop hazelcast

# Relancer
docker start hazelcast


go get github.com/hazelcast/hazelcast-go-client@v1.2.2



require github.com/hazelcast/hazelcast-go-client v1.2.2 // ou autre version


go mod tidy

docker build -t mon-api-hazelcast .
docker run --rm mon-api-hazelcast











++++++++++++++++++++++++


               +--------------+
               |  Application |
               |    (Go API)  |
               +------+-------+
                      |
                      | 1. GET/POST
                      v
              +-------+--------+
              |   Hazelcast    | ←→ mémoire rapide (RAM)
              +-------+--------+
                      |
        Miss Cache     | Hit Cache
            |          v
            v     (réponse rapide)
      +-----+------+
      |   Database  | ←→ PostgreSQL, MySQL, MongoDB...
      +------------+

curl -s "http://localhost:3000/api/hotels" | jq '.data[] | select(.city | contains("New Mia"))'
curl -s "http://localhost:3000/api/hotels" | jq '.data[] | select(.city | contains("New Mia"))'
# ✅ Rechargement des données
curl http://localhost:3000/api/reload

# ✅ Accès à l'hôtel Carter Group dans la ville "New Mia"  
curl http://localhost:3000/api/hotels/HOT317

# ✅ Recherche d'hôtels par ville
curl -s "http://localhost:3000/api/hotels" | jq '.data[] | select(.city | contains("New Mia"))'