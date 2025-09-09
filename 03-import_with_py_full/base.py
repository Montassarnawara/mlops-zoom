import psycopg2
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT
from faker import Faker
import random
from datetime import timedelta

fake = Faker()

# 1. Création base de données hotel_db (si pas existante)
def create_database():
    try:
        conn = psycopg2.connect(
            dbname="postgres",
            user="myuser",
            password="mypassword",
            host="localhost",
            port="5433"
        )
        conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT)
        cur = conn.cursor()
        cur.execute("CREATE DATABASE hotel_db")
        print("✅ Base 'hotel_db' créée.")
        cur.close()
        conn.close()
    except psycopg2.errors.DuplicateDatabase:
        print("⚠️ La base 'hotel_db' existe déjà.")
    except Exception as e:
        print("Erreur création base:", e)


# 2. Connexion à hotel_db
def connect_db():
    conn = psycopg2.connect(
        dbname="hotel_db",
        user="myuser",
        password="mypassword",
        host="localhost",
        port="5433"
    )
    return conn


# 3. Création des tables
def create_tables(cur):
    cur.execute("""
    DROP TABLE IF EXISTS CONTRACT;
    DROP TABLE IF EXISTS USER_APP;
    DROP TABLE IF EXISTS ROLE;
    DROP TABLE IF EXISTS HOTEL;
    DROP TABLE IF EXISTS CITY;
    DROP TABLE IF EXISTS COUNTRY;

    CREATE TABLE COUNTRY (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        code VARCHAR(10)
    );
    
    CREATE TABLE CITY (
        id SERIAL PRIMARY KEY,
        citycode VARCHAR(20),
        name VARCHAR(100),
        countryname VARCHAR(100),
        countryid INT REFERENCES COUNTRY(id)
    );

    CREATE TABLE HOTEL (
        hotel_key VARCHAR(20) PRIMARY KEY,
        name VARCHAR(100),
        city VARCHAR(100),
        country VARCHAR(100),
        stars INT,
        address TEXT,
        mail VARCHAR(100),
        latitude FLOAT,
        longitude FLOAT,
        phone VARCHAR(50),
        description TEXT,
        short_description TEXT
    );

    CREATE TABLE ROLE (
        id SERIAL PRIMARY KEY,
        name VARCHAR(50)
    );

    CREATE TABLE USER_APP (
        id SERIAL PRIMARY KEY,
        status VARCHAR(50),
        marge FLOAT,
        marge_operation FLOAT,
        solde FLOAT,
        solde_rouge FLOAT,
        currency VARCHAR(10),
        maxrequest INT,
        "group" VARCHAR(50),
        marge_b2b FLOAT,
        marge_xml FLOAT,
        role_id INT REFERENCES ROLE(id)
    );

    CREATE TABLE CONTRACT (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        hotel_id VARCHAR(20) REFERENCES HOTEL(hotel_key),
        start_at DATE,
        end_at DATE,
        access VARCHAR(50),
        active BOOLEAN,
        currency VARCHAR(10),
        market INT,
        client_id INT REFERENCES USER_APP(id)
    );
    """)

# 4. Insertion données factices
def seed_country(cur, n=5):
    for _ in range(n):
        cur.execute("INSERT INTO COUNTRY (name, code) VALUES (%s, %s)",
                    (fake.country(), fake.country_code()))
        
def seed_city(cur, n=10):
    cur.execute("SELECT id, name FROM COUNTRY")
    countries = cur.fetchall()
    for _ in range(n):
        country = random.choice(countries)
        cur.execute("""
            INSERT INTO CITY (citycode, name, countryname, countryid)
            VALUES (%s, %s, %s, %s)
        """, (fake.city_suffix(), fake.city(), country[1], country[0]))

def seed_hotel(cur, n=10):
    cur.execute("SELECT name FROM CITY")
    cities = cur.fetchall()
    for _ in range(n):
        city = random.choice(cities)[0]
        country = fake.country()
        cur.execute("""
            INSERT INTO HOTEL (hotel_key, name, city, country, stars, address, mail,
            latitude, longitude, phone, description, short_description)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            fake.unique.bothify(text='HOT###'),
            fake.company(),
            city,
            country,
            random.randint(1, 5),
            fake.address(),
            fake.company_email(),
            fake.latitude(),
            fake.longitude(),
            fake.phone_number(),
            fake.text(max_nb_chars=200),
            fake.sentence()
        ))

def seed_role(cur):
    roles = ['Admin', 'Agent', 'Client']
    for r in roles:
        cur.execute("INSERT INTO ROLE (name) VALUES (%s)", (r,))

def seed_users(cur, n=10):
    cur.execute("SELECT id FROM ROLE")
    roles = cur.fetchall()
    for _ in range(n):
        role_id = random.choice(roles)[0]
        cur.execute("""
            INSERT INTO USER_APP (status, marge, marge_operation, solde, solde_rouge,
            currency, maxrequest, "group", marge_b2b, marge_xml, role_id)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            random.choice(['actif', 'suspendu']),
            round(random.uniform(5, 20), 2),
            round(random.uniform(1, 10), 2),
            round(random.uniform(100, 1000), 2),
            -100.0,
            random.choice(['EUR', 'USD', 'TND']),
            random.randint(10, 1000),
            random.choice(['A', 'B', 'C']),
            round(random.uniform(1, 5), 2),
            round(random.uniform(0.5, 3), 2),
            role_id
        ))

def seed_contract(cur, n=10):
    cur.execute("SELECT hotel_key FROM HOTEL")
    hotels = cur.fetchall()
    cur.execute("SELECT id FROM USER_APP")
    users = cur.fetchall()

    for _ in range(n):
        hotel_id = random.choice(hotels)[0]
        client_id = random.choice(users)[0]
        start = fake.date_between(start_date='-1y', end_date='today')
        end = start + timedelta(days=random.randint(30, 365))
        cur.execute("""
            INSERT INTO CONTRACT (name, hotel_id, start_at, end_at, access, active,
            currency, market, client_id)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            fake.bs().title(),
            hotel_id,
            start,
            end,
            random.choice(['Basic', 'Premium', 'Exclusive']),
            random.choice([True, False]),
            random.choice(['EUR', 'USD', 'TND']),
            random.randint(1, 5),
            client_id
        ))

def main():
    create_database()

    conn = connect_db()
    cur = conn.cursor()

    create_tables(cur)
    conn.commit()
    print("✅ Tables créées.")

    seed_country(cur)
    seed_city(cur)
    seed_role(cur)
    seed_hotel(cur)
    seed_users(cur)
    seed_contract(cur)
    conn.commit()

    print("✅ Données factices insérées.")

    # Fermer connexion
    cur.close()
    conn.close()

if __name__ == "__main__":
    main()
