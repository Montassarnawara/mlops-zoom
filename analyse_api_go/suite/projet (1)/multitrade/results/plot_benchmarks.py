
import os
import pandas as pd
import matplotlib.pyplot as plt


# S'assurer que le dossier 'results/' existe
os.makedirs("results", exist_ok=True)

# Lecture du CSV
csv_path = "benchmark.csv"
df = pd.read_csv(csv_path)

# Liste des colonnes à tracer
fields = df['Colonne'].unique()

# 1. Graphe : Temps d'exécution (Durée) en fonction de N pour chaque colonne
plt.figure(figsize=(10,6))
for field in fields:
    sub = df[df['Colonne'] == field]
    plt.plot(sub['N'], sub['Durée(µs)'], marker='o', label=field)
plt.xlabel('N')
plt.ylabel('Durée (µs)')
plt.title('Temps d\'exécution en fonction de N')
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.savefig('results/benchmark_time.png')
plt.close()

# 2. Graphe : Moyenne en fonction de N pour chaque colonne
plt.figure(figsize=(10,6))
for field in fields:
    sub = df[df['Colonne'] == field]
    plt.plot(sub['N'], sub['Moyenne'], marker='o', label=field)
plt.xlabel('N')
plt.ylabel('Moyenne')
plt.title('Moyenne par colonne en fonction de N')
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.savefig('results/benchmark_avg.png')
plt.close()

# 3. Graphe : Somme en fonction de N pour chaque colonne
plt.figure(figsize=(10,6))
for field in fields:
    sub = df[df['Colonne'] == field]
    plt.plot(sub['N'], sub['Somme'], marker='o', label=field)
plt.xlabel('N')
plt.ylabel('Somme')
plt.title('Somme par colonne en fonction de N')
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.savefig('results/benchmark_sum.png')
plt.close()

print("Graphiques enregistrés dans le dossier results/")
