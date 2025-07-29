import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

# Charger les données
df = pd.read_csv("benchmark.csv")

# Nettoyage (au cas où)
df = df.dropna()
df["endpoint"] = df["endpoint"].str.strip()
df["status_code"] = df["status_code"].astype(str)

# Liste des endpoints uniques
endpoints = df["endpoint"].unique()

# 1. 📊 Graphe combiné : durée vs valeur pour chaque endpoint (comparatif global)
plt.figure(figsize=(12, 6))
sns.lineplot(data=df, x="valeur", y="duration_microseconds", hue="endpoint", style="status_code", markers=True)
plt.title("Comparaison du temps d'exécution par endpoint")
plt.xlabel("Valeur (ex: index, taille...)")
plt.ylabel("Durée (µs)")
plt.legend(title="Endpoint")
plt.grid(True)
plt.tight_layout()
plt.savefig("global_comparaison.png")
plt.close()

# 2. 📈 Graphe individuel pour chaque endpoint
for ep in endpoints:
    plt.figure(figsize=(10, 5))
    sub_df = df[df["endpoint"] == ep]
    sns.lineplot(data=sub_df, x="valeur", y="duration_microseconds", marker="o")
    plt.title(f"Temps d'exécution - Endpoint: {ep}")
    plt.xlabel("Valeur")
    plt.ylabel("Durée (µs)")
    plt.grid(True)
    plt.tight_layout()
    plt.savefig(f"courbe_{ep}.png")
    plt.close()

# 3. 📉 Histogrammes de distribution du temps pour chaque endpoint
for ep in endpoints:
    plt.figure(figsize=(8, 4))
    sub_df = df[df["endpoint"] == ep]
    sns.histplot(sub_df["duration_microseconds"], bins=20, kde=True)
    plt.title(f"Distribution des temps - {ep}")
    plt.xlabel("Durée (µs)")
    plt.ylabel("Fréquence")
    plt.tight_layout()
    plt.savefig(f"histogramme_{ep}.png")
    plt.close()

# 4. ⚫ Nuage de points pour voir les outliers ou comportements anormaux
for ep in endpoints:
    plt.figure(figsize=(8, 4))
    sub_df = df[df["endpoint"] == ep]
    sns.scatterplot(data=sub_df, x="valeur", y="duration_microseconds", hue="status_code")
    plt.title(f"Nuage de points - {ep}")
    plt.xlabel("Valeur")
    plt.ylabel("Durée (µs)")
    plt.tight_layout()
    plt.savefig(f"scatter_{ep}.png")
    plt.close()

# 5. 📊 Boxplot global pour détecter les outliers par endpoint
plt.figure(figsize=(10, 5))
sns.boxplot(data=df, x="endpoint", y="duration_microseconds", hue="status_code")
plt.title("Boxplot des temps d'exécution par endpoint et statut")
plt.xlabel("Endpoint")
plt.ylabel("Durée (µs)")
plt.tight_layout()
plt.savefig("boxplot_endpoint.png")
plt.close()

print("✅ Graphiques générés dans le dossier courant.")
