
ssh -p 50175 root@192.99.232.126
TGBYHN2025++++

go get github.com/gin-gonic/gin

# post
POST http://localhost:8080/user
Content-Type: application/json

{
  "Id": "10",
  "Name": "Yassine",
  "Validation": true
}

#GET :

GET http://localhost:8080/user/10
curl http://localhost:8080/college/CLG0057


 Tu importes :

    context : pour gérer les opérations concurrentes (asynchrone / timeout).

    json : pour encoder/décoder les objets entre Go ↔ JSON.

    log : pour les messages d’erreur/diagnostic.

    net/http : pour les constantes HTTP (200, 404, etc.).

    gin : framework HTTP pour créer une API rapide.

    hazelcast-go-client : client Go pour Hazelcast (base en RAM distribuée).


+++++++++++++
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








go get github.com/hazelcast/hazelcast-go-client

docker run -d --name hazelcast \
  -p 5701:5701 \
  -e HZ_CLUSTERNAME=dev-go \
  hazelcast/hazelcast:5.3


# Lancer le fichier Go
go run main.go

# Voir les logs de Hazelcast (optionnel pour debug)
docker logs -f hazelcast

# Arrêter Hazelcast
docker stop hazelcast

# Relancer
docker start hazelcast

go mod init hizelcon

go get github.com/hazelcast/hazelcast-go-client


docker rm -f hazelcast

docker run -d --name hazelcast \
  -p 5701:5701 \
  -e HZ_CLUSTERNAME=dev-go \
  hazelcast/hazelcast:5.3

docker stop hazelcast
docker start hazelcast


vérifie si ton port est bien ouvert localement.
telnet 127.0.0.1 5701
nc -zv 127.0.0.1 5701



ngrok http 8000

 Étapes d’installation rapide de ngrok (mode root, sans sudo)
📌 1. Télécharge le binaire ngrok :

wget https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-stable-linux-amd64.zip

📌 2. Dézippe-le :

unzip ngrok-stable-linux-amd64.zip

📌 3. Rends-le exécutable et déplace-le dans /usr/local/bin :

chmod +x ngrok
mv ngrok /usr/local/bin/

apt update
apt install unzip


# Après inscription sur ngrok.com, récupérer le token puis dans la VM :
ngrok config add-authtoken <token>

# Lancer l’API dans un terminal :
go run main.go

# Dans un autre terminal :
ngrok http 8888


TRUNCATE TABLE statistiques_api RESTART IDENTITY;


# Créer l'env
python3 -m venv envpy

# Activer l'env
source envpy/bin/activate

# Vérifier l'interpréteur Python dans l'env
which python3
python3 --version

# Sortir de l'env
deactivate   # Si ça ne marche pas, ferme terminal ou lance un nouveau shell bash

go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/gin-gonic/gin
go get github.com/lib/pq

Étapes générales
1. 📁 Créer le dossier du projet

mkdir echo_api && cd echo_api

2. 🧰 Initialiser un module Go

go mod init echo_api

3. ⬇️ Installer Echo et PostgreSQL

go get github.com/labstack/echo/v4
go get github.com/lib/pq

Téléchargement Go (version stable la plus récente)

wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz

2. Supprimer toute ancienne version (par sécurité)

rm -rf /usr/local/go

3. Décompresser et installer

tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz

4. Ajouter Go au PATH :

echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

5. Vérifier l’installation :

go versionapt install curl -y

apt install postgresql postgresql-contrib -y
systemctl start postgresql

Ouvre ou édite ton fichier .bashrc :

nano ~/.bashrc

🔍 Trouve ou ajoute une ligne comme :

PS1='\[\e[1;31m\](envpy)\[\e[0m\] \[\e[1;34m\]\u@\h:\w\[\e[0m\]# '



Installation manuelle (recommandée par Go)

    Télécharge la dernière version stable (par exemple Go 1.22.3) :

wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz

    Extrais l'archive dans /usr/local :

sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz

    Ajoute Go au PATH (ajoute dans ton ~/.bashrc, ~/.zshrc, ou ~/.profile) :

echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

scp -P 50175 -r mon_dossier root@192.99.232.126:/root/

scp -P 50175 -r root@192.99.232.126:/root/analyser_api .


scp -P 50175 -r go-hazelcast-api root@192.99.232.126:/root/



 root@go-175:~# sudo nano /etc/postgresql/*/main/pg_hba.conf
 root@go-175:~# sudo nano /etc/postgresql/*/main/postgresql.conf
 root@go-175:~# sudo systemctl restart postgresql
 root@go-175:~# sudo ss -tunlp | grep 5432
tcp   LISTEN 0      244        127.0.0.1:5432      0.0.0.0:*    users:(("postgres",pid=18366,fd=6))                  
tcp   LISTEN 0      244            [::1]:5432         [::]:*    users:(("postgres",pid=18366,fd=5))                  
 root@go-175:~# psql -h 10.10.50.175 -U aurora -d hazelcastdb -p 5432
Password for user aurora: 
psql: error: connection to server at "10.10.50.175", port 5432 failed: FATAL:  password authentication failed for user "aurora"
connection to server at "10.10.50.175", port 5432 failed: FATAL:  password authentication failed for user "aurora"