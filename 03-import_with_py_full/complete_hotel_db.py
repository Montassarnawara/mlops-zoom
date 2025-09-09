"""
🏨 HOTEL DATABASE - CRÉATION COMPLÈTE DES 24 TABLES
==================================================

Ce script crée et peuple une base de données complète pour un système de gestion hôtelière
avec 24 tables interconnectées selon l'architecture définie.

Tables principales:
- Tables de base: COUNTRY, CITY, HOTEL, USER, ROLE, CONTRACT
- Tables de tarification: BASE_PRICE, BOARD_PRICES, ROOM_PRICES, ROOM_TYPE_PRICES
- Tables de promotion: SPECIAL_OFFER, DAY_PROMOTION, ROOM_PROMOTION, PAX_PROMOTION
- Tables de structure: ROOM, ROOM_TYPE, BOARD, PERIOD
- Tables de liaison: CONTRACT_ROOM, PERIOD_BOARD
- Tables de gestion: SUPPLEMENT, CHILD_PRICE, CANCELLATION_CONDITION, ALLOTMENT

Auteur: Système automatisé
Date: 2025-08-05
Version: 1.0
"""

import psycopg2
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT
from faker import Faker
import random
from datetime import timedelta, date
import uuid

fake = Faker()

# ===============================================
# 🔌 CONFIGURATION ET CONNEXION
# ===============================================

def create_database():
    """Créer la base de données hotel_db si elle n'existe pas"""
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
        print(f"❌ Erreur création base: {e}")

def connect_db():
    """Connexion à la base hotel_db"""
    try:
        conn = psycopg2.connect(
            dbname="hotel_db",
            user="myuser",
            password="mypassword", 
            host="localhost",
            port="5433"
        )
        return conn
    except Exception as e:
        print(f"❌ Erreur connexion: {e}")
        return None

# ===============================================
# 🏗️ CRÉATION DES TABLES (ordre des dépendances)
# ===============================================

def create_all_tables(cur):
    """Créer toutes les 24 tables dans l'ordre correct des dépendances"""
    
    print("🗑️ Suppression des tables existantes...")
    # Suppression dans l'ordre inverse des dépendances
    tables_to_drop = [
        'ALLOTMENT', 'ROOM_TYPE_PRICES', 'CANCELLATION_CONDITION', 'CHILD_PRICE',
        'ROOM_PRICES', 'SUPPLEMENT', 'BOARD_PRICES', 'PERIOD_BOARD', 'BASE_PRICE',
        'PAX_PROMOTION', 'ROOM_PROMOTION', 'DAY_PROMOTION', 'SPECIAL_OFFER',
        'CONTRACT_ROOM', 'PERIOD', 'CONTRACT', 'BOARD', 'ROOM', 'ROOM_TYPE',
        'USER_APP', 'ROLE', 'HOTEL', 'CITY', 'COUNTRY'
    ]
    
    for table in tables_to_drop:
        cur.execute(f"DROP TABLE IF EXISTS {table} CASCADE;")
    
    print("🏗️ Création des tables de base...")
    
    # 1. COUNTRY - Table de base (pas de dépendances)
    cur.execute("""
        CREATE TABLE COUNTRY (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            code VARCHAR(10) UNIQUE NOT NULL
        );
    """)
    
    # 2. CITY - Dépend de COUNTRY
    cur.execute("""
        CREATE TABLE CITY (
            id SERIAL PRIMARY KEY,
            citycode VARCHAR(20) NOT NULL,
            name VARCHAR(100) NOT NULL,
            countryname VARCHAR(100),
            countryid INT REFERENCES COUNTRY(id) ON DELETE CASCADE
        );
    """)
    
    # 3. HOTEL - Dépend de CITY et COUNTRY (références string)
    cur.execute("""
        CREATE TABLE HOTEL (
            hotel_key VARCHAR(20) PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            city VARCHAR(100),
            country VARCHAR(100),
            stars INT CHECK (stars BETWEEN 1 AND 5),
            address TEXT,
            mail VARCHAR(100),
            latitude FLOAT,
            longitude FLOAT,
            phone VARCHAR(50),
            description TEXT,
            short_description TEXT
        );
    """)
    
    # 4. ROLE - Table de base pour les utilisateurs
    cur.execute("""
        CREATE TABLE ROLE (
            id SERIAL PRIMARY KEY,
            name VARCHAR(50) UNIQUE NOT NULL
        );
    """)
    
    # 5. USER_APP - Dépend de ROLE
    cur.execute("""
        CREATE TABLE USER_APP (
            id SERIAL PRIMARY KEY,
            status VARCHAR(50) DEFAULT 'actif',
            marge FLOAT DEFAULT 0.0,
            marge_operation FLOAT DEFAULT 0.0,
            solde FLOAT DEFAULT 0.0,
            solde_rouge FLOAT DEFAULT 0.0,
            currency VARCHAR(10) DEFAULT 'EUR',
            maxrequest INT DEFAULT 100,
            "group" VARCHAR(50),
            marge_b2b FLOAT DEFAULT 0.0,
            marge_xml FLOAT DEFAULT 0.0,
            role_id INT REFERENCES ROLE(id) ON DELETE SET NULL
        );
    """)
    
    # 6. CONTRACT - Dépend de HOTEL et USER_APP
    cur.execute("""
        CREATE TABLE CONTRACT (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            hotel_id VARCHAR(20) REFERENCES HOTEL(hotel_key) ON DELETE CASCADE,
            start_at DATE NOT NULL,
            end_at DATE NOT NULL,
            access VARCHAR(50),
            active BOOLEAN DEFAULT TRUE,
            currency VARCHAR(10) DEFAULT 'EUR',
            market INT DEFAULT 1,
            client_id INT REFERENCES USER_APP(id) ON DELETE SET NULL,
            CHECK (end_at > start_at)
        );
    """)
    
    print("🏗️ Création des tables de structure...")
    
    # 7. ROOM_TYPE - Types de chambres
    cur.execute("""
        CREATE TABLE ROOM_TYPE (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            code VARCHAR(20) UNIQUE NOT NULL,
            client_id INT REFERENCES USER_APP(id) ON DELETE SET NULL
        );
    """)
    
    # 8. ROOM - Chambres individuelles
    cur.execute("""
        CREATE TABLE ROOM (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            name_code VARCHAR(20) UNIQUE NOT NULL,
            client_id INT REFERENCES USER_APP(id) ON DELETE SET NULL,
            max_pax INT DEFAULT 4,
            min_pax INT DEFAULT 1,
            child INT DEFAULT 2,
            min_adult INT DEFAULT 1,
            max_adult INT DEFAULT 4,
            CHECK (max_pax >= min_pax),
            CHECK (max_adult >= min_adult)
        );
    """)
    
    # 9. BOARD - Types de pension
    cur.execute("""
        CREATE TABLE BOARD (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            definition TEXT,
            client_id INT REFERENCES USER_APP(id) ON DELETE SET NULL
        );
    """)
    
    # 10. PERIOD - Périodes de contrat
    cur.execute("""
        CREATE TABLE PERIOD (
            id SERIAL PRIMARY KEY,
            contract_id INT REFERENCES CONTRACT(id) ON DELETE CASCADE,
            start_at DATE NOT NULL,
            end_at DATE NOT NULL,
            minimum_stay INT DEFAULT 1,
            delai INT DEFAULT 0,
            min_stay INT DEFAULT 1,
            name VARCHAR(100),
            code VARCHAR(20),
            CHECK (end_at > start_at),
            CHECK (minimum_stay > 0)
        );
    """)
    
    print("🏗️ Création des tables de liaison...")
    
    # 11. CONTRACT_ROOM - Liaison contrat/chambre
    cur.execute("""
        CREATE TABLE CONTRACT_ROOM (
            id SERIAL PRIMARY KEY,
            contract_id INT REFERENCES CONTRACT(id) ON DELETE CASCADE,
            room_id INT REFERENCES ROOM(id) ON DELETE CASCADE,
            room_type_id INT REFERENCES ROOM_TYPE(id) ON DELETE CASCADE,
            maximum_number_people INT DEFAULT 4,
            minimum_number_people INT DEFAULT 1,
            max_adult INT DEFAULT 4,
            min_adult INT DEFAULT 1,
            child INT DEFAULT 2,
            code_room VARCHAR(20),
            code_type VARCHAR(20),
            CHECK (maximum_number_people >= minimum_number_people),
            CHECK (max_adult >= min_adult)
        );
    """)
    
    # 12. PERIOD_BOARD - Liaison période/pension
    cur.execute("""
        CREATE TABLE PERIOD_BOARD (
            id SERIAL PRIMARY KEY,
            period_id INT REFERENCES PERIOD(id) ON DELETE CASCADE,
            board_id INT REFERENCES BOARD(id) ON DELETE CASCADE,
            extra_board_id INT REFERENCES BOARD(id) ON DELETE SET NULL
        );
    """)
    
    print("🏗️ Création des tables de tarification...")
    
    # 13. BASE_PRICE - Prix de base
    cur.execute("""
        CREATE TABLE BASE_PRICE (
            id SERIAL PRIMARY KEY,
            room_id INT REFERENCES ROOM(id) ON DELETE CASCADE,
            period_id INT REFERENCES PERIOD(id) ON DELETE CASCADE,
            type_id INT REFERENCES ROOM_TYPE(id) ON DELETE CASCADE,
            board_id INT REFERENCES BOARD(id) ON DELETE CASCADE,
            price FLOAT NOT NULL DEFAULT 0.0,
            par_pax BOOLEAN DEFAULT FALSE,
            extb_price FLOAT DEFAULT 0.0,
            operation VARCHAR(10) DEFAULT '+',
            CHECK (price >= 0),
            CHECK (extb_price >= 0)
        );
    """)
    
    # 14. BOARD_PRICES - Prix des pensions
    cur.execute("""
        CREATE TABLE BOARD_PRICES (
            id SERIAL PRIMARY KEY,
            board_id INT REFERENCES BOARD(id) ON DELETE CASCADE,
            period_id INT REFERENCES PERIOD(id) ON DELETE CASCADE,
            price FLOAT NOT NULL DEFAULT 0.0,
            extb_price FLOAT DEFAULT 0.0,
            par_pax BOOLEAN DEFAULT FALSE,
            CHECK (price >= 0),
            CHECK (extb_price >= 0)
        );
    """)
    
    # 15. ROOM_PRICES - Prix par chambre
    cur.execute("""
        CREATE TABLE ROOM_PRICES (
            id SERIAL PRIMARY KEY,
            room_id INT REFERENCES ROOM(id) ON DELETE CASCADE,
            period_id INT REFERENCES PERIOD(id) ON DELETE CASCADE,
            price FLOAT NOT NULL DEFAULT 0.0,
            extb_price FLOAT DEFAULT 0.0,
            par_pax BOOLEAN DEFAULT FALSE,
            CHECK (price >= 0),
            CHECK (extb_price >= 0)
        );
    """)
    
    # 16. ROOM_TYPE_PRICES - Prix par type de chambre
    cur.execute("""
        CREATE TABLE ROOM_TYPE_PRICES (
            id SERIAL PRIMARY KEY,
            type_id INT REFERENCES ROOM_TYPE(id) ON DELETE CASCADE,
            period_id INT REFERENCES PERIOD(id) ON DELETE CASCADE,
            price FLOAT NOT NULL DEFAULT 0.0,
            extb_price FLOAT DEFAULT 0.0,
            par_pax BOOLEAN DEFAULT FALSE,
            CHECK (price >= 0),
            CHECK (extb_price >= 0)
        );
    """)
    
    print("🏗️ Création des tables de promotion...")
    
    # 17. SPECIAL_OFFER - Offres spéciales
    cur.execute("""
        CREATE TABLE SPECIAL_OFFER (
            id SERIAL PRIMARY KEY,
            receive_date FLOAT,
            requestdate_from FLOAT,
            requestdate_to FLOAT,
            checkin_date FLOAT,
            checkout_date FLOAT,
            min_stay INT DEFAULT 1,
            max_stay INT DEFAULT 30,
            "order" INT DEFAULT 0,
            offer_type VARCHAR(50),
            origine VARCHAR(100),
            definition TEXT,
            note TEXT,
            bookingdate_range_exclusif BOOLEAN DEFAULT FALSE,
            in_contrat BOOLEAN DEFAULT FALSE,
            priority BOOLEAN DEFAULT FALSE,
            contrat_id INT REFERENCES CONTRACT(id) ON DELETE CASCADE,
            client_id INT REFERENCES USER_APP(id) ON DELETE SET NULL,
            CHECK (max_stay >= min_stay)
        );
    """)
    
    # 18. DAY_PROMOTION - Promotions par jour
    cur.execute("""
        CREATE TABLE DAY_PROMOTION (
            id SERIAL PRIMARY KEY,
            special_offer_id INT REFERENCES SPECIAL_OFFER(id) ON DELETE CASCADE,
            rooms TEXT,
            room_types TEXT,
            boards TEXT,
            stay_day VARCHAR(50),
            pay_day VARCHAR(50)
        );
    """)
    
    # 19. ROOM_PROMOTION - Promotions par chambre
    cur.execute("""
        CREATE TABLE ROOM_PROMOTION (
            id SERIAL PRIMARY KEY,
            special_offer_id INT REFERENCES SPECIAL_OFFER(id) ON DELETE CASCADE,
            rooms TEXT,
            room_types TEXT,
            boards TEXT,
            value FLOAT DEFAULT 0.0,
            value_operation VARCHAR(10) DEFAULT '%'
        );
    """)
    
    # 20. PAX_PROMOTION - Promotions par nombre de personnes
    cur.execute("""
        CREATE TABLE PAX_PROMOTION (
            id SERIAL PRIMARY KEY,
            special_offer_id INT REFERENCES SPECIAL_OFFER(id) ON DELETE CASCADE,
            rooms TEXT,
            room_types TEXT,
            boards TEXT,
            pax_stay INT DEFAULT 1,
            pax_pay INT DEFAULT 1,
            CHECK (pax_stay > 0),
            CHECK (pax_pay > 0)
        );
    """)
    
    print("🏗️ Création des tables de gestion...")
    
    # 21. SUPPLEMENT - Suppléments tarifaires
    cur.execute("""
        CREATE TABLE SUPPLEMENT (
            id SERIAL PRIMARY KEY,
            period_id INT REFERENCES PERIOD(id) ON DELETE CASCADE,
            room_id INT REFERENCES ROOM(id) ON DELETE CASCADE,
            type_id INT REFERENCES ROOM_TYPE(id) ON DELETE CASCADE,
            board_id INT REFERENCES BOARD(id) ON DELETE CASCADE,
            adults INT DEFAULT 1,
            contract_id INT REFERENCES CONTRACT(id) ON DELETE CASCADE,
            price FLOAT DEFAULT 0.0,
            operation VARCHAR(10) DEFAULT '+',
            CHECK (adults > 0),
            CHECK (price >= 0)
        );
    """)
    
    # 22. CHILD_PRICE - Prix enfants
    cur.execute("""
        CREATE TABLE CHILD_PRICE (
            id SERIAL PRIMARY KEY,
            period_id INT REFERENCES PERIOD(id) ON DELETE CASCADE,
            room_id INT REFERENCES ROOM(id) ON DELETE CASCADE,
            type_id INT REFERENCES ROOM_TYPE(id) ON DELETE CASCADE,
            board_id INT REFERENCES BOARD(id) ON DELETE CASCADE,
            adults INT DEFAULT 1,
            contract_id INT REFERENCES CONTRACT(id) ON DELETE CASCADE,
            childs TEXT,
            CHECK (adults > 0)
        );
    """)
    
    # 23. CANCELLATION_CONDITION - Conditions d'annulation
    cur.execute("""
        CREATE TABLE CANCELLATION_CONDITION (
            id SERIAL PRIMARY KEY,
            period_id INT REFERENCES PERIOD(id) ON DELETE CASCADE,
            board_id INT REFERENCES BOARD(id) ON DELETE CASCADE,
            room_type_id_cancel INT REFERENCES ROOM_TYPE(id) ON DELETE CASCADE,
            room_id INT REFERENCES ROOM(id) ON DELETE CASCADE,
            no_show_operation VARCHAR(20),
            max_days_before_arrival INT DEFAULT 0,
            min_days_before_arrival INT DEFAULT 0,
            nights_to_bill INT DEFAULT 1,
            no_show BOOLEAN DEFAULT FALSE,
            no_show_nights_to_bill INT DEFAULT 1,
            free_cancel BOOLEAN DEFAULT TRUE,
            free_cancel_before VARCHAR(20),
            refundable BOOLEAN DEFAULT TRUE,
            operation VARCHAR(10) DEFAULT '%',
            CHECK (max_days_before_arrival >= min_days_before_arrival),
            CHECK (nights_to_bill >= 0),
            CHECK (no_show_nights_to_bill >= 0)
        );
    """)
    
    # 24. ALLOTMENT - Quotas de chambres
    cur.execute("""
        CREATE TABLE ALLOTMENT (
            id SERIAL PRIMARY KEY,
            period_id INT REFERENCES PERIOD(id) ON DELETE CASCADE,
            room_id INT REFERENCES ROOM(id) ON DELETE CASCADE,
            number INT DEFAULT 1,
            CHECK (number > 0)
        );
    """)
    
    print("✅ Toutes les 24 tables créées avec succès !")

# ===============================================
# 📊 PEUPLEMENT DES DONNÉES
# ===============================================

def seed_countries(cur, n=10):
    """Peupler la table COUNTRY avec des pays réalistes"""
    countries = [
        ('France', 'FR'), ('Espagne', 'ES'), ('Italie', 'IT'), ('Allemagne', 'DE'),
        ('Royaume-Uni', 'GB'), ('Portugal', 'PT'), ('Grèce', 'GR'), ('Turquie', 'TR'),
        ('Maroc', 'MA'), ('Tunisie', 'TN'), ('Égypte', 'EG'), ('États-Unis', 'US'),
        ('Canada', 'CA'), ('Japon', 'JP'), ('Australie', 'AU')
    ]
    
    for i, (name, code) in enumerate(countries[:n]):
        cur.execute(
            "INSERT INTO COUNTRY (name, code) VALUES (%s, %s)",
            (name, code)
        )
    print(f"✅ {n} pays insérés")

def seed_cities(cur, n=20):
    """Peupler la table CITY"""
    cur.execute("SELECT id, name FROM COUNTRY")
    countries = cur.fetchall()
    
    cities = [
        'Paris', 'Madrid', 'Rome', 'Berlin', 'Londres', 'Lisbonne', 'Athènes', 'Istanbul',
        'Casablanca', 'Tunis', 'Le Caire', 'New York', 'Toronto', 'Tokyo', 'Sydney',
        'Barcelone', 'Milan', 'Munich', 'Manchester', 'Porto', 'Thessalonique', 'Antalya',
        'Marrakech', 'Sousse', 'Alexandrie', 'Los Angeles', 'Vancouver', 'Osaka', 'Melbourne'
    ]
    
    for i, city_name in enumerate(cities[:n]):
        country = random.choice(countries)
        city_code = f"CITY{i+1:03d}"
        cur.execute("""
            INSERT INTO CITY (citycode, name, countryname, countryid)
            VALUES (%s, %s, %s, %s)
        """, (city_code, city_name, country[1], country[0]))
    print(f"✅ {n} villes insérées")

def seed_roles(cur):
    """Peupler la table ROLE avec des rôles standards"""
    roles = ['Admin', 'Agent', 'Client', 'Manager', 'Operator']
    for role in roles:
        cur.execute("INSERT INTO ROLE (name) VALUES (%s)", (role,))
    print(f"✅ {len(roles)} rôles insérés")

def seed_users(cur, n=15):
    """Peupler la table USER_APP"""
    cur.execute("SELECT id FROM ROLE")
    roles = cur.fetchall()
    
    statuses = ['actif', 'suspendu', 'en_attente', 'bloqué']
    currencies = ['EUR', 'USD', 'TND', 'GBP', 'CAD']
    groups = ['A', 'B', 'C', 'VIP', 'STANDARD']
    
    for i in range(n):
        role_id = random.choice(roles)[0]
        cur.execute("""
            INSERT INTO USER_APP (
                status, marge, marge_operation, solde, solde_rouge, currency,
                maxrequest, "group", marge_b2b, marge_xml, role_id
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            random.choice(statuses),
            round(random.uniform(5, 25), 2),
            round(random.uniform(1, 10), 2),
            round(random.uniform(100, 5000), 2),
            -100.0,
            random.choice(currencies),
            random.randint(50, 1000),
            random.choice(groups),
            round(random.uniform(1, 8), 2),
            round(random.uniform(0.5, 5), 2),
            role_id
        ))
    print(f"✅ {n} utilisateurs insérés")

def seed_hotels(cur, n=12):
    """Peupler la table HOTEL"""
    cur.execute("SELECT name FROM CITY")
    cities = cur.fetchall()
    cur.execute("SELECT name FROM COUNTRY")
    countries = cur.fetchall()
    
    for i in range(n):
        city = random.choice(cities)[0]
        country = random.choice(countries)[0]
        hotel_key = f"HOT{i+1:03d}"
        
        cur.execute("""
            INSERT INTO HOTEL (
                hotel_key, name, city, country, stars, address, mail,
                latitude, longitude, phone, description, short_description
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            hotel_key,
            fake.company() + " Hotel",
            city,
            country,
            random.randint(1, 5),
            fake.address().replace('\n', ', '),
            fake.company_email(),
            round(fake.latitude(), 6),
            round(fake.longitude(), 6),
            fake.phone_number(),
            fake.text(max_nb_chars=300),
            fake.sentence()
        ))
    print(f"✅ {n} hôtels insérés")

def seed_contracts(cur, n=15):
    """Peupler la table CONTRACT"""
    cur.execute("SELECT hotel_key FROM HOTEL")
    hotels = cur.fetchall()
    cur.execute("SELECT id FROM USER_APP")
    users = cur.fetchall()
    
    accesses = ['Basic', 'Premium', 'Exclusive', 'VIP', 'Standard']
    currencies = ['EUR', 'USD', 'TND', 'GBP']
    
    for i in range(n):
        hotel_id = random.choice(hotels)[0]
        client_id = random.choice(users)[0]
        start_date = fake.date_between(start_date='-6M', end_date='today')
        end_date = start_date + timedelta(days=random.randint(90, 730))
        
        cur.execute("""
            INSERT INTO CONTRACT (
                name, hotel_id, start_at, end_at, access, active,
                currency, market, client_id
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            f"Contrat {fake.company()} {i+1}",
            hotel_id,
            start_date,
            end_date,
            random.choice(accesses),
            random.choice([True, True, True, False]),  # 75% actifs
            random.choice(currencies),
            random.randint(1, 5),
            client_id
        ))
    print(f"✅ {n} contrats insérés")

def seed_room_types(cur, n=8):
    """Peupler la table ROOM_TYPE"""
    room_types = [
        ('Single Standard', 'SS'), ('Double Standard', 'DS'), ('Double Superior', 'DSP'),
        ('Junior Suite', 'JS'), ('Suite', 'ST'), ('Family Room', 'FR'),
        ('Triple', 'TR'), ('Quadruple', 'QR')
    ]
    
    cur.execute("SELECT id FROM USER_APP LIMIT 3")
    users = cur.fetchall()
    
    for i, (name, code) in enumerate(room_types[:n]):
        client_id = random.choice(users)[0] if random.random() < 0.3 else None
        cur.execute("""
            INSERT INTO ROOM_TYPE (name, code, client_id)
            VALUES (%s, %s, %s)
        """, (name, code, client_id))
    print(f"✅ {n} types de chambres insérés")

def seed_rooms(cur, n=20):
    """Peupler la table ROOM"""
    cur.execute("SELECT id FROM USER_APP LIMIT 5")
    users = cur.fetchall()
    
    for i in range(n):
        room_name = f"Chambre {i+1:03d}"
        room_code = f"R{i+1:03d}"
        client_id = random.choice(users)[0] if random.random() < 0.2 else None
        
        min_pax = random.randint(1, 2)
        max_pax = random.randint(min_pax + 1, 6)
        min_adult = random.randint(1, 2)
        max_adult = random.randint(min_adult, max_pax)
        
        cur.execute("""
            INSERT INTO ROOM (
                name, name_code, client_id, max_pax, min_pax,
                child, min_adult, max_adult
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            room_name, room_code, client_id, max_pax, min_pax,
            random.randint(0, 2), min_adult, max_adult
        ))
    print(f"✅ {n} chambres insérées")

def seed_boards(cur, n=6):
    """Peupler la table BOARD"""
    boards = [
        ('Room Only', 'Chambre seule sans repas'),
        ('Bed & Breakfast', 'Chambre avec petit-déjeuner'),
        ('Half Board', 'Demi-pension (petit-déjeuner + dîner)'),
        ('Full Board', 'Pension complète (tous les repas)'),
        ('All Inclusive', 'Tout inclus (repas + boissons)'),
        ('All Inclusive Premium', 'Tout inclus premium avec extras')
    ]
    
    cur.execute("SELECT id FROM USER_APP LIMIT 3")
    users = cur.fetchall()
    
    for name, definition in boards[:n]:
        client_id = random.choice(users)[0] if random.random() < 0.1 else None
        cur.execute("""
            INSERT INTO BOARD (name, definition, client_id)
            VALUES (%s, %s, %s)
        """, (name, definition, client_id))
    print(f"✅ {n} types de pensions insérés")

def seed_periods(cur, n=25):
    """Peupler la table PERIOD"""
    cur.execute("SELECT id FROM CONTRACT")
    contracts = cur.fetchall()
    
    for i in range(n):
        contract_id = random.choice(contracts)[0]
        start_date = fake.date_between(start_date='-3M', end_date='+6M')
        end_date = start_date + timedelta(days=random.randint(30, 180))
        
        cur.execute("""
            INSERT INTO PERIOD (
                contract_id, start_at, end_at, minimum_stay, delai,
                min_stay, name, code
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            contract_id, start_date, end_date,
            random.randint(1, 7), random.randint(0, 5),
            random.randint(1, 3), f"Période {i+1}",
            f"P{i+1:03d}"
        ))
    print(f"✅ {n} périodes insérées")

def seed_contract_rooms(cur, n=30):
    """Peupler la table CONTRACT_ROOM"""
    cur.execute("SELECT id FROM CONTRACT")
    contracts = cur.fetchall()
    cur.execute("SELECT id FROM ROOM")
    rooms = cur.fetchall()
    cur.execute("SELECT id FROM ROOM_TYPE")
    room_types = cur.fetchall()
    
    for i in range(n):
        contract_id = random.choice(contracts)[0]
        room_id = random.choice(rooms)[0]
        room_type_id = random.choice(room_types)[0]
        
        min_people = random.randint(1, 2)
        max_people = random.randint(min_people + 1, 6)
        min_adult = random.randint(1, 2)
        max_adult = random.randint(min_adult, max_people)
        
        cur.execute("""
            INSERT INTO CONTRACT_ROOM (
                contract_id, room_id, room_type_id, maximum_number_people,
                minimum_number_people, max_adult, min_adult, child,
                code_room, code_type
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            contract_id, room_id, room_type_id, max_people,
            min_people, max_adult, min_adult, random.randint(0, 2),
            f"CR{i+1:03d}", f"CT{i+1:03d}"
        ))
    print(f"✅ {n} liaisons contrat-chambre insérées")

def seed_period_boards(cur, n=40):
    """Peupler la table PERIOD_BOARD"""
    cur.execute("SELECT id FROM PERIOD")
    periods = cur.fetchall()
    cur.execute("SELECT id FROM BOARD")
    boards = cur.fetchall()
    
    for i in range(n):
        period_id = random.choice(periods)[0]
        board_id = random.choice(boards)[0]
        extra_board_id = random.choice(boards)[0] if random.random() < 0.3 else None
        
        cur.execute("""
            INSERT INTO PERIOD_BOARD (period_id, board_id, extra_board_id)
            VALUES (%s, %s, %s)
        """, (period_id, board_id, extra_board_id))
    print(f"✅ {n} liaisons période-pension insérées")

def seed_pricing_tables(cur):
    """Peupler toutes les tables de tarification"""
    cur.execute("SELECT id FROM ROOM")
    rooms = cur.fetchall()
    cur.execute("SELECT id FROM PERIOD") 
    periods = cur.fetchall()
    cur.execute("SELECT id FROM ROOM_TYPE")
    room_types = cur.fetchall()
    cur.execute("SELECT id FROM BOARD")
    boards = cur.fetchall()
    
    # BASE_PRICE
    for i in range(50):
        room_id = random.choice(rooms)[0]
        period_id = random.choice(periods)[0]
        type_id = random.choice(room_types)[0]
        board_id = random.choice(boards)[0]
        
        cur.execute("""
            INSERT INTO BASE_PRICE (
                room_id, period_id, type_id, board_id, price,
                par_pax, extb_price, operation
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            room_id, period_id, type_id, board_id,
            round(random.uniform(50, 300), 2),
            random.choice([True, False]),
            round(random.uniform(10, 50), 2),
            random.choice(['+', '%'])
        ))
    
    # BOARD_PRICES
    for i in range(30):
        board_id = random.choice(boards)[0]
        period_id = random.choice(periods)[0]
        
        cur.execute("""
            INSERT INTO BOARD_PRICES (
                board_id, period_id, price, extb_price, par_pax
            ) VALUES (%s, %s, %s, %s, %s)
        """, (
            board_id, period_id,
            round(random.uniform(20, 100), 2),
            round(random.uniform(5, 25), 2),
            random.choice([True, False])
        ))
    
    # ROOM_PRICES
    for i in range(40):
        room_id = random.choice(rooms)[0]
        period_id = random.choice(periods)[0]
        
        cur.execute("""
            INSERT INTO ROOM_PRICES (
                room_id, period_id, price, extb_price, par_pax
            ) VALUES (%s, %s, %s, %s, %s)
        """, (
            room_id, period_id,
            round(random.uniform(80, 400), 2),
            round(random.uniform(15, 60), 2),
            random.choice([True, False])
        ))
    
    # ROOM_TYPE_PRICES
    for i in range(35):
        type_id = random.choice(room_types)[0]
        period_id = random.choice(periods)[0]
        
        cur.execute("""
            INSERT INTO ROOM_TYPE_PRICES (
                type_id, period_id, price, extb_price, par_pax
            ) VALUES (%s, %s, %s, %s, %s)
        """, (
            type_id, period_id,
            round(random.uniform(60, 350), 2),
            round(random.uniform(12, 55), 2),
            random.choice([True, False])
        ))
    
    print("✅ Tables de tarification peuplées (BASE_PRICE, BOARD_PRICES, ROOM_PRICES, ROOM_TYPE_PRICES)")

def seed_promotion_tables(cur):
    """Peupler les tables de promotion"""
    cur.execute("SELECT id FROM CONTRACT")
    contracts = cur.fetchall()
    cur.execute("SELECT id FROM USER_APP")
    users = cur.fetchall()
    
    # SPECIAL_OFFER
    offer_types = ['Early Booking', 'Last Minute', 'Stay 7 Pay 6', 'Honeymoon', 'Family']
    origines = ['Website', 'Agent', 'Direct', 'Email', 'Phone']
    
    special_offers = []
    for i in range(20):
        contract_id = random.choice(contracts)[0]
        client_id = random.choice(users)[0]
        
        cur.execute("""
            INSERT INTO SPECIAL_OFFER (
                receive_date, requestdate_from, requestdate_to, checkin_date,
                checkout_date, min_stay, max_stay, "order", offer_type,
                origine, definition, note, bookingdate_range_exclusif,
                in_contrat, priority, contrat_id, client_id
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            random.uniform(1640995200, 1672531200),  # 2022-2023 timestamps
            random.uniform(1672531200, 1704067200),  # 2023-2024
            random.uniform(1704067200, 1735689600),  # 2024-2025
            random.uniform(1704067200, 1735689600),
            random.uniform(1735689600, 1767225600),  # 2025-2026
            random.randint(1, 7), random.randint(7, 30), i,
            random.choice(offer_types), random.choice(origines),
            fake.text(max_nb_chars=200), fake.sentence(),
            random.choice([True, False]), random.choice([True, False]),
            random.choice([True, False]), contract_id, client_id
        ))
        special_offers.append(cur.lastrowid if hasattr(cur, 'lastrowid') else i+1)
    
    cur.execute("SELECT id FROM SPECIAL_OFFER")
    special_offers = [row[0] for row in cur.fetchall()]
    
    # DAY_PROMOTION
    days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
    for i in range(15):
        offer_id = random.choice(special_offers)
        stay_days = ','.join(random.sample(days, random.randint(1, 4)))
        pay_days = ','.join(random.sample(days, random.randint(1, 3)))
        
        cur.execute("""
            INSERT INTO DAY_PROMOTION (
                special_offer_id, rooms, room_types, boards, stay_day, pay_day
            ) VALUES (%s, %s, %s, %s, %s, %s)
        """, (
            offer_id, 'Room101,Room102', 'DS,JS', 'BB,HB',
            stay_days, pay_days
        ))
    
    # ROOM_PROMOTION
    operations = ['%', '-', '+']
    for i in range(18):
        offer_id = random.choice(special_offers)
        
        cur.execute("""
            INSERT INTO ROOM_PROMOTION (
                special_offer_id, rooms, room_types, boards, value, value_operation
            ) VALUES (%s, %s, %s, %s, %s, %s)
        """, (
            offer_id, 'Room201,Room202', 'SS,DS',
            'RO,BB', round(random.uniform(5, 30), 2),
            random.choice(operations)
        ))
    
    # PAX_PROMOTION
    for i in range(12):
        offer_id = random.choice(special_offers)
        pax_stay = random.randint(1, 4)
        pax_pay = random.randint(1, pax_stay)
        
        cur.execute("""
            INSERT INTO PAX_PROMOTION (
                special_offer_id, rooms, room_types, boards, pax_stay, pax_pay
            ) VALUES (%s, %s, %s, %s, %s, %s)
        """, (
            offer_id, 'Room301,Room302', 'FR,QR',
            'FB,AI', pax_stay, pax_pay
        ))
    
    print("✅ Tables de promotion peuplées (SPECIAL_OFFER, DAY_PROMOTION, ROOM_PROMOTION, PAX_PROMOTION)")

def seed_management_tables(cur):
    """Peupler les tables de gestion"""
    cur.execute("SELECT id FROM PERIOD")
    periods = cur.fetchall()
    cur.execute("SELECT id FROM ROOM")
    rooms = cur.fetchall()
    cur.execute("SELECT id FROM ROOM_TYPE")
    room_types = cur.fetchall()
    cur.execute("SELECT id FROM BOARD")
    boards = cur.fetchall()
    cur.execute("SELECT id FROM CONTRACT")
    contracts = cur.fetchall()
    
    # SUPPLEMENT
    operations = ['+', '%', '-']
    for i in range(25):
        cur.execute("""
            INSERT INTO SUPPLEMENT (
                period_id, room_id, type_id, board_id, adults,
                contract_id, price, operation
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            random.choice(periods)[0], random.choice(rooms)[0],
            random.choice(room_types)[0], random.choice(boards)[0],
            random.randint(1, 4), random.choice(contracts)[0],
            round(random.uniform(10, 80), 2), random.choice(operations)
        ))
    
    # CHILD_PRICE
    child_descriptions = ['0-2 ans gratuit', '3-12 ans 50%', '13-17 ans 80%', 'Enfant supplémentaire']
    for i in range(20):
        cur.execute("""
            INSERT INTO CHILD_PRICE (
                period_id, room_id, type_id, board_id, adults,
                contract_id, childs
            ) VALUES (%s, %s, %s, %s, %s, %s, %s)
        """, (
            random.choice(periods)[0], random.choice(rooms)[0],
            random.choice(room_types)[0], random.choice(boards)[0],
            random.randint(1, 2), random.choice(contracts)[0],
            random.choice(child_descriptions)
        ))
    
    # CANCELLATION_CONDITION
    no_show_ops = ['charge_full', 'charge_partial', 'no_charge']
    operations = ['%', 'fixed', 'nights']
    for i in range(30):
        max_days = random.randint(5, 30)
        min_days = random.randint(0, max_days-1)
        
        cur.execute("""
            INSERT INTO CANCELLATION_CONDITION (
                period_id, board_id, room_type_id_cancel, room_id,
                no_show_operation, max_days_before_arrival, min_days_before_arrival,
                nights_to_bill, no_show, no_show_nights_to_bill, free_cancel,
                free_cancel_before, refundable, operation
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            random.choice(periods)[0], random.choice(boards)[0],
            random.choice(room_types)[0], random.choice(rooms)[0],
            random.choice(no_show_ops), max_days, min_days,
            random.randint(1, 3), random.choice([True, False]),
            random.randint(1, 2), random.choice([True, False]),
            f"{random.randint(1, 15)} days", random.choice([True, False]),
            random.choice(operations)
        ))
    
    # ALLOTMENT
    for i in range(35):
        cur.execute("""
            INSERT INTO ALLOTMENT (period_id, room_id, number)
            VALUES (%s, %s, %s)
        """, (
            random.choice(periods)[0], random.choice(rooms)[0],
            random.randint(1, 20)
        ))
    
    print("✅ Tables de gestion peuplées (SUPPLEMENT, CHILD_PRICE, CANCELLATION_CONDITION, ALLOTMENT)")

def get_database_stats(cur):
    """Afficher les statistiques de la base de données"""
    print("\n📊 STATISTIQUES DE LA BASE DE DONNÉES")
    print("=" * 50)
    
    tables = [
        'COUNTRY', 'CITY', 'HOTEL', 'ROLE', 'USER_APP', 'CONTRACT',
        'ROOM_TYPE', 'ROOM', 'BOARD', 'PERIOD', 'CONTRACT_ROOM', 'PERIOD_BOARD',
        'BASE_PRICE', 'BOARD_PRICES', 'ROOM_PRICES', 'ROOM_TYPE_PRICES',
        'SPECIAL_OFFER', 'DAY_PROMOTION', 'ROOM_PROMOTION', 'PAX_PROMOTION',
        'SUPPLEMENT', 'CHILD_PRICE', 'CANCELLATION_CONDITION', 'ALLOTMENT'
    ]
    
    total_records = 0
    for table in tables:
        cur.execute(f"SELECT COUNT(*) FROM {table}")
        count = cur.fetchone()[0]
        total_records += count
        print(f"📋 {table:<25} : {count:>6} enregistrements")
    
    print("=" * 50)
    print(f"🎯 TOTAL                     : {total_records:>6} enregistrements")
    print(f"🏗️ TABLES CRÉÉES             : {len(tables):>6} tables")
    print("=" * 50)

def main():
    """Fonction principale d'exécution"""
    print("🚀 DÉMARRAGE - Création complète de la base hôtelière")
    print("=" * 60)
    
    # Étape 1: Créer la base de données
    create_database()
    
    # Étape 2: Se connecter à la base
    conn = connect_db()
    if not conn:
        print("❌ Impossible de se connecter à la base de données")
        return
    
    cur = conn.cursor()
    
    try:
        # Étape 3: Créer toutes les tables
        print("\n🏗️ CRÉATION DES TABLES")
        print("-" * 30)
        create_all_tables(cur)
        conn.commit()
        
        # Étape 4: Peupler les tables de base
        print("\n📊 PEUPLEMENT DES DONNÉES DE BASE")
        print("-" * 40)
        seed_countries(cur, 10)
        seed_cities(cur, 20)
        seed_roles(cur)
        seed_users(cur, 15)
        seed_hotels(cur, 12)
        seed_contracts(cur, 15)
        conn.commit()
        
        # Étape 5: Peupler les tables de structure
        print("\n🏗️ PEUPLEMENT DES STRUCTURES")
        print("-" * 35)
        seed_room_types(cur, 8)
        seed_rooms(cur, 20)
        seed_boards(cur, 6)
        seed_periods(cur, 25)
        conn.commit()
        
        # Étape 6: Peupler les tables de liaison
        print("\n🔗 PEUPLEMENT DES LIAISONS")
        print("-" * 30)
        seed_contract_rooms(cur, 30)
        seed_period_boards(cur, 40)
        conn.commit()
        
        # Étape 7: Peupler les tables de tarification
        print("\n💰 PEUPLEMENT DES TARIFICATIONS")
        print("-" * 35)
        seed_pricing_tables(cur)
        conn.commit()
        
        # Étape 8: Peupler les tables de promotion
        print("\n🎁 PEUPLEMENT DES PROMOTIONS")
        print("-" * 32)
        seed_promotion_tables(cur)
        conn.commit()
        
        # Étape 9: Peupler les tables de gestion
        print("\n⚙️ PEUPLEMENT DE LA GESTION")
        print("-" * 30)
        seed_management_tables(cur)
        conn.commit()
        
        # Étape 10: Afficher les statistiques
        get_database_stats(cur)
        
        print("\n🎉 SUCCÈS COMPLET !")
        print("=" * 60)
        print("✅ Base de données hôtelière créée avec succès")
        print("✅ 24 tables créées et peuplées") 
        print("✅ Données cohérentes avec contraintes respectées")
        print("✅ Prêt pour l'intégration avec Hazelcast")
        print("=" * 60)
        
    except Exception as e:
        print(f"❌ Erreur lors de l'exécution: {e}")
        conn.rollback()
    finally:
        cur.close()
        conn.close()

if __name__ == "__main__":
    main()
