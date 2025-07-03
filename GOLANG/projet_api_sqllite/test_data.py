import sqlite3

conn = sqlite3.connect("product.db")
cursor = conn.cursor()

# Voir les noms de tables
cursor.execute("SELECT name FROM sqlite_master WHERE type='table';")
print("Tables:", cursor.fetchall())

# Voir un échantillon de la table
cursor.execute("SELECT * FROM product LIMIT 5;")
for row in cursor.fetchall():
    print(row)

conn.close()
