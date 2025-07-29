hotel-sync/
│
├── main.go                        # Point d’entrée
├── config/
│   └── postgres.go               # Connexion PostgreSQL
├── cache/
│   └── hazelcast.go             # Connexion Hazelcast
├── models/
│   ├── contract.go
│   ├── hotel.go
│   ├── city.go
│   ├── country.go
│   ├── user.go
│   └── role.go
├── services/
│   └── importer.go              # Logique d’importation des données
└── go.mod / go.sum
