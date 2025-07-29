ssh -p 50175 root@192.99.232.126
TGBYHN2025++++
TGBYHN2025++++
curl http://localhost:8888/college/echo/CLG0015
curl http://localhost:8888/college/fiber/CLG0015
curl http://localhost:8888/college/gin/CLG0015
curl http://localhost:8888/college/net/CLG0015
curl http://localhost:8888/user


# POST
curl -X POST http://localhost:8080/college/echo \
  -H "Content-Type: application/json" \
  -d '{
    "communication_skills": 8,
    "projects_completed": 3,
    "prev_sem_result": 15.5,
    "cgpa": 3.6,
    "academic_performance": 9,
    "iq": 120,
    "extra_curricular_score": 6,
    "placement": "no",
    "internship_experience": "2 months at Google",
    "college_id": "CLG0015"
}'
# PUT
curl -X PUT http://localhost:8080/college/echo/CLG0015/increment
curl -X PUT http://localhost:8080/college/echo/CLG0015/toggle-placement
# DEL
curl -X DELETE http://localhost:8080/college/echo/CLG0015

 


lsof -i :8080
kill -9 7192
go run main.go
cd GO_test

//get 
curl -i http://localhost:8888/users
curl -i http://localhost:8888/users/1
//post
curl -i -X POST http://localhost:8888/users \
  -H "Content-Type: application/json" \
  -d '{"id":"5", "name":"youssef", "validation":true}'

//pacht
curl -i -X PATCH http://localhost:8888/users/2 \
  -H "Content-Type: application/json" \
  -d '{"name":"amine", "validation":false}'

// delete 
curl -i -X DELETE http://localhost:8888/users


// head
curl -i -X HEAD http://localhost:8888/users

 go mod init github.com/postg_api
 go get github.com/gin-gonic/gin

go mod init github.com/Montassarnawara/postg_api
go get github.com/gin-gonic/gin


# Créer l'env
python3 -m venv envpy

# Activer l'env
source envpy/bin/activate

# Vérifier l'interpréteur Python dans l'env
which python3
python3 --version

# Sortir de l'env
deactivate   # Si ça ne marche pas, ferme terminal ou lance un nouveau shell bash

python3 -c "import numpy; print(numpy.__version__)"
python3 -c "import pandas as pd; print(pd.__version__)"

///////////////////coloare la commende 
nano ~/.bashrc
PS1='\[\e[1;31m\](envpy)\[\e[0m\]\[\e[1;34m\]\u@\h:\w\[\e[0m\]# '
source ~/.bashrc

 root@go-175:~/postg_api# su - postgres
postgres@go-175:~$ psql         


scp -P 50175 college_student_placement_dataset.csv root@192.99.232.126:/root/postg_api/
root@192.99.232.126's password: 
college_student_placement_dataset.csv                                                     

scp -P 50175 imdb_tvshows.csv root@192.99.232.126:/root/postg_api/

scp -P 50175 import_with_py root@192.99.232.126:/root

scp -P 50175 -r hotel_hazelcast root@192.99.232.126:/root
scp -P 50175 -r sys_api_hotel root@192.99.232.126:/root
scp -P 50175 -r api_zeta_max root@192.99.232.126:/root
TGBYHN2025++++


scp -P 50175 -r root@192.99.232.126:/root/analyser_api .


mv duration_by_path.png analyser_api/
