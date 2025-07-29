import pandas as pd
import matplotlib.pyplot as plt

df = pd.read_csv("benchmark.csv")

plt.figure(figsize=(10, 6))
plt.plot(df["N"], df["Durée(µs)"], marker='o', color='blue', linestyle='-')

plt.title("Temps de traitement (µs) en fonction du nombre de données N")
plt.xlabel("Nombre de données (N)")
plt.ylabel("Durée (µs)")
plt.grid(True)
plt.tight_layout()

# ✅ Enregistre dans un fichier au lieu de plt.show()
plt.savefig("results/benchmark_plot.png")
print("✅ Graphique sauvegardé : results/benchmark_plot.png")
