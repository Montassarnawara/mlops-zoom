import psycopg2

conn = psycopg2.connect(
    dbname="products_db",
    user="myuser",
    password="mypassword",
    host="localhost",
    port="5433"
)

cur = conn.cursor()

cur.execute("SELECT * FROM products; ")
rows = cur.fetchall()

for row in rows:
    print(row)

cur.close()
conn.close()
