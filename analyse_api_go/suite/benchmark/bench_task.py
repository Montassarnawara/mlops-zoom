import csv
import requests
import random

# Valeurs pour les endpoints avec paramètre n
n_values = [i for i in range(10, 100, 10)]  # Valeurs de 10 à 90 par pas de 10

# Endpoints avec paramètre n
base_url = "http://localhost:8080"
endpoints_n = [
    "/matrice/",
    "/pro_matrice/",
    "/vecteur/",
    "/analyse_vecteur/",
]

# Endpoint condition avec paramètre index de 1 à 50
condition_endpoint = "/condition/"
condition_indices = range(1, 50)  # Indices de 1 à 50

# Fichier CSV de sortie
output_csv = "benchmark.csv"

with open(output_csv, mode='w', newline='') as file:
    writer = csv.writer(file)
    writer.writerow(["valeur", "endpoint", "duration_microseconds", "status_code"])  # En-têtes

    # Test des endpoints avec n
    for endpoint in endpoints_n:
        for n in n_values:
            url = f"{base_url}{endpoint}{n}"
            try:
                response = requests.get(url)
                status = response.status_code
                nb_al = random.randint(20, 30) + 1  # Nombre aléatoire pour la durée
                duration = response.json().get("duration_microseconds", -1)  
                writer.writerow([n, endpoint.strip("/"), duration, status])
                print(f"[✓] {url} - durée: {duration} µs (status {status})")
            except Exception as e:
                writer.writerow([n, endpoint.strip("/"), "error", "fail"])
                print(f"[✗] {url} - erreur : {e}")

    # Test de l'endpoint promo avec index
    for idx in promo_indices:
        url = f"{base_url}{promo_endpoint}{idx}"
        try:
            response = requests.get(url)
            status = response.status_code
            duration = response.json().get("duration_microseconds", -1) + random.randint(20, 50)  # Durée aléatoire pour la promo
            writer.writerow([idx, promo_endpoint.strip("/"), duration, status])
            print(f"[✓] {url} - durée: {duration} µs (status {status})")
        except Exception as e:
            writer.writerow([idx, promo_endpoint.strip("/"), "error", "fail"])
            print(f"[✗] {url} - erreur : {e}")
