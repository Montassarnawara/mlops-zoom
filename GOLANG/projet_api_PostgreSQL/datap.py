import psycopg2
import pandas as pd

# Lire CSV
prod = pd.read_csv('/workspaces/mlops-zoom/GOLANG/projet_api_sqllite/product.csv')

# Connexion PostgreSQL
conn = psycopg2.connect(
    dbname="products_db",
    user="myuser",
    password="mypassword",
    host="localhost",
    port="5433"
)

cur = conn.cursor()

# Création table (adapter types selon données)
cur.execute("""
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name TEXT,
    price NUMERIC,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
""")

# Insérer ligne par ligne
for _, row in prod.iterrows():
    cur.execute(
        "INSERT INTO products (name, price, created_at, updated_at) VALUES (%s, %s, %s, %s)",
        (row['name'], row['price'], row['created_at'], row['updated_at'])
    )

conn.commit()
cur.close()
conn.close()
