import psycopg2

def connect_db():
    conn = psycopg2.connect(
        dbname="hotel_db",
        user="myuser",
        password="mypassword",
        host="localhost",
        port="5433"
    )
    return conn

def show_table_data(cur, table_name, limit=5):
    print(f"\n===== Contenu de la table {table_name} (limit {limit}) =====")
    cur.execute(f"SELECT * FROM {table_name} LIMIT {limit}")
    rows = cur.fetchall()
    for row in rows:
        print(row)

def main():
    conn = connect_db()
    cur = conn.cursor()

    tables = ["COUNTRY", "CITY", "HOTEL", "ROLE", "USER_APP", "CONTRACT"]

    for table in tables:
        show_table_data(cur, table)

    cur.close()
    conn.close()

if __name__ == "__main__":
    main()
