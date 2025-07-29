#!/bin/bash

# Vérifie que l'utilisateur donne un nombre n
if [ -z "$1" ]; then
  echo "❌ Utilisation : $0 <nombre_de_requetes>"
  exit 1
fi

n=$1

for i in $(seq 1 $n); do
  # Formater l'index en 4 chiffres avec padding (ex: 1 => 0001)
  id=$(printf "CLG%04d" $i)

  echo "🔍 Requête pour $id"
  curl -s http://localhost:8080/college/$id | jq  # jq pour formater JSON (si installé)
  echo -e "\n----------------------"
done
